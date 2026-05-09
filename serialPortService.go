package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"github.com/wailsapp/wails/v3/pkg/application"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type SerialPortService struct {
	platformAdapter PlatformAdapter // 平台适配器
	ctx             context.Context
	conn            io.ReadWriteCloser
	app             *application.App
	cancel          context.CancelFunc // 用于通知 readLoop 退出
	mu              sync.Mutex         // 互斥锁，保护 conn
}

func (s *SerialPortService) GetSerialPorts() ([]string, error) {
	var ports []string
	switch runtime.GOOS {
	case "windows":
		for i := 1; i <= 255; i++ {
			portName := fmt.Sprintf("COM%d", i)
			if s.isPortOpenable(portName) {
				ports = append(ports, portName)
			}
		}
	case "darwin":
		files, err := s.listDevFiles()
		if err == nil {
			for _, f := range files {
				if len(f) > 4 && (f[:4] == "tty." || f[:3] == "cu.") {
					if strings.HasPrefix(f, "cu.") {
						ports = append(ports, "/dev/"+f)
					}
				}
			}
		}
	case "linux":
		files, err := s.listDevFiles()
		if err == nil {
			for _, f := range files {
				if len(f) > 3 && (f[:6] == "ttyUSB" || f[:6] == "ttyACM" || f[:5] == "ttyS") {
					ports = append(ports, "/dev/"+f)
				}
			}
		}
	}

	return ports, nil
}

func (s *SerialPortService) listDevFiles() ([]string, error) {
	var files []string
	entries, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}

func (s *SerialPortService) isPortOpenable(portName string) bool {
	options := serial.OpenOptions{
		PortName: portName,
		BaudRate: 9600, // 随便给个波特率
	}
	_, err := serial.Open(options)
	if err == nil {
		return true
	}
	return false
}

func (s *SerialPortService) OpenSerial(portName string, baudRate uint, dataBits uint, stopBits uint, parityMode int) error {
	//var ParityMode = serial.PARITY_NONE
	//switch {
	//case parityMode == 1:
	//	ParityMode = serial.PARITY_ODD
	//case parityMode == 2:
	//	ParityMode = serial.PARITY_EVEN
	//default:
	//	ParityMode = serial.PARITY_NONE
	//}
	//options := serial.OpenOptions{
	//	PortName:              portName,
	//	BaudRate:              baudRate,
	//	DataBits:              dataBits,
	//	StopBits:              stopBits,
	//	ParityMode:            ParityMode,
	//	MinimumReadSize:       1,
	//	InterCharacterTimeout: 1000,
	//}
	//return s.platformAdapter.Open(options)
	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		s.conn.Close()
	}
	var ParityMode = serial.PARITY_NONE
	switch {
	case parityMode == 1:
		ParityMode = serial.PARITY_ODD
	case parityMode == 2:
		ParityMode = serial.PARITY_EVEN
	default:
		ParityMode = serial.PARITY_NONE
	}
	options := serial.OpenOptions{
		PortName:              portName,
		BaudRate:              baudRate,
		DataBits:              dataBits,
		StopBits:              stopBits,
		ParityMode:            ParityMode,
		MinimumReadSize:       1,
		InterCharacterTimeout: 1000,
	}
	var err error
	s.conn, err = serial.Open(options)
	if err != nil {
		return fmt.Errorf("无法打开端口 %s: %v", portName, err)
	}
	log.Printf("成功打开串口: %s", portName)
	// 启动读取协程
	s.cancel = cancel // 保存 cancel 函数
	s.ctx = ctx
	go s.readLoop()
	return nil
}

func (s *SerialPortService) readLoop() {
	defer func() {
		log.Println("读取协程已退出")
		s.mu.Lock()
		s.cancel = nil
		s.mu.Unlock()
	}()

	reader := bufio.NewReader(s.conn)

	for {
		// 创建一个 1 秒的超时定时器
		readTimeout := time.After(1 * time.Second)

		// 用于接收读取结果的通道
		readResult := make(chan struct {
			line string
			err  error
		}, 1)

		// 在一个新的 goroutine 中执行阻塞读取
		go func() {
			line, err := reader.ReadString('\n')
			readResult <- struct {
				line string
				err  error
			}{line, err}
		}()

		// 使用 select 等待：要么是读取完成，要么是超时，要么是收到取消信号
		select {
		case <-s.ctx.Done():
			// 1. 收到取消信号，立即退出
			log.Println("收到取消信号，退出读取循环")
			return

		case <-readTimeout:
			continue

		case result := <-readResult:
			if result.err != nil {
				if result.err == io.EOF || s.ctx.Err() != nil {
					return
				}
				log.Printf("读取错误: %v", result.err)
				return
			}
			s.app.Event.Emit("serial_data", result.line)
		}
	}
}

func (s *SerialPortService) CloseSerial() error {
	s.mu.Lock()
	conn := s.conn
	s.conn = nil
	cancel := s.cancel
	s.mu.Unlock()
	if conn == nil {
		return nil
	}
	if cancel != nil {
		cancel()
	}
	done := make(chan error, 1)
	go func() {
		done <- conn.Close()
	}()
	select {
	case err := <-done:
		if err != nil && !strings.Contains(err.Error(), "already closed") {
			return err
		}
		return nil
	case <-time.After(5 * time.Second):
		log.Println("警告：串口关闭超时，强制返回")
		return fmt.Errorf("串口关闭超时（可能是设备已断开或驱动卡死）")
	}
}

func (s *SerialPortService) SendData(data string) error {
	if s.conn == nil {
		return fmt.Errorf("串口未打开，请先连接")
	}
	// 写入数据
	// 注意：如果协议需要换行符，请在这里拼接，例如 data + "\n"
	// fmt.Println("接收数据", data)
	_, err := s.conn.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("发送失败: %v", err)
	}
	return nil
}
