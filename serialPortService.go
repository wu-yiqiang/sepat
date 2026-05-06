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
)

type SerialPortService struct {
	ctx    context.Context
	conn   io.ReadWriteCloser
	app    *application.App
	cancel context.CancelFunc // 用于通知 readLoop 退出
	mu     sync.Mutex         // 互斥锁，保护 conn
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

func (s *SerialPortService) OpenSerial(portName string, baudRate uint, dataBits uint, stopBits uint) error {
	ctx, cancel := context.WithCancel(context.Background())

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		s.conn.Close()
	}
	options := serial.OpenOptions{
		PortName:              portName,
		BaudRate:              baudRate,
		DataBits:              dataBits,
		StopBits:              stopBits,
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
		// 清理 cancel，防止内存泄漏
		s.mu.Lock()
		s.cancel = nil
		s.mu.Unlock()
	}()

	reader := bufio.NewReader(s.conn)

	for {
		fmt.Println("开始获取数据")
		// 检查上下文是否已取消
		select {
		case <-s.ctx.Done():
			fmt.Println("关闭连接了")
			return
		default:
			fmt.Println("读取数据前")
			line, err := reader.ReadString('\n')
			fmt.Println("读取数据后")
			if err != nil {
				// 如果是 EOF 或者上下文取消，直接退出
				if err == io.EOF || s.ctx.Err() != nil {
					return
				}
				log.Printf("读取错误: %v", err)
				continue
			}
			fmt.Println("获取数据了", line)
			s.app.Event.Emit("serial_data", line)
		}
	}
}

func (s *SerialPortService) CloseSerial() error {
	s.mu.Lock()
	conn := s.conn // 保存引用
	s.conn = nil   // 立即置空
	s.mu.Unlock()

	if conn == nil {
		return nil
	}

	// 1. 先关闭底层连接
	// 这会强制 readLoop 中的 ReadString 立即返回错误 (通常是 "file already closed")
	err := conn.Close()

	// 2. 再取消 Context，确保 readLoop 能收到信号
	if s.cancel != nil {
		s.cancel()
	}

	// 忽略 "file already closed" 错误，因为可能是我们自己关闭的
	if err != nil && !strings.Contains(err.Error(), "already closed") {
		return err
	}
	return nil
}

func (s *SerialPortService) SendData(data string) error {
	if s.conn == nil {
		return fmt.Errorf("串口未打开，请先连接")
	}

	// 写入数据
	// 注意：如果协议需要换行符，请在这里拼接，例如 data + "\n"
	fmt.Println("接受数据", data)
	_, err := s.conn.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("发送失败: %v", err)
	}
	return nil
}
