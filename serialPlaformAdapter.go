package main

import (
	"github.com/jacobsa/go-serial/serial"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type PlatformAdapter interface {
	Open(serial.OpenOptions) error
	Close() error
}

func NewPlatformAdapter(service *SerialPortService) PlatformAdapter {
	switch runtime.GOOS {
	case "windows":
		return &WindowsAdapter{service: service}
	case "darwin": // macOS
		return &MacOSAdapter{service: service}
	default:
		return &MacOSAdapter{service: service}
		//default: // Linux 等其他平台
		//	return &GenericAdapter{service: service}
	}
}

type WindowsAdapter struct {
	service *SerialPortService
}

func (a *WindowsAdapter) Open(openOptions serial.OpenOptions) error {
	a.service.mu.Lock()
	if a.service.conn != nil {
		a.service.cancel()
		a.service.conn = nil
	}
	a.service.mu.Unlock()
	time.Sleep(200 * time.Millisecond)
	options := serial.OpenOptions{
		PortName:              openOptions.PortName,
		BaudRate:              openOptions.BaudRate,
		DataBits:              openOptions.DataBits,
		StopBits:              openOptions.StopBits,
		ParityMode:            openOptions.ParityMode,
		MinimumReadSize:       1,
		InterCharacterTimeout: 1000,
	}
	_, err := serial.Open(options)
	return err
}

func (a *WindowsAdapter) Close() error {
	a.service.mu.Lock()
	_ = a.service.conn
	a.service.conn = nil
	cancel := a.service.cancel
	a.service.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	return nil
}

type MacOSAdapter struct {
	service *SerialPortService
}

func (a *MacOSAdapter) Open(openOptions serial.OpenOptions) error {
	_ = normalizeMacOSPath(openOptions.PortName)
	killSystemModemManager()
	options := serial.OpenOptions{
		PortName:              openOptions.PortName,
		BaudRate:              openOptions.BaudRate,
		DataBits:              openOptions.DataBits,
		StopBits:              openOptions.StopBits,
		ParityMode:            openOptions.ParityMode,
		MinimumReadSize:       1,
		InterCharacterTimeout: 1000,
	}
	_, err := serial.Open(options)
	return err
}

func (a *MacOSAdapter) Close() error {
	a.service.mu.Lock()
	conn := a.service.conn
	a.service.conn = nil
	cancel := a.service.cancel
	a.service.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	// 1. 带超时关闭（macOS 必须）
	done := make(chan error, 1)
	go func() { done <- conn.Close() }()

	select {
	case err := <-done:
		return err
	case <-time.After(500 * time.Millisecond):
		log.Println("macOS: 强制跳过串口关闭（系统占用）")
		return nil
	}
}

// ====== macOS 专属辅助函数 ======
func normalizeMacOSPath(path string) string {
	if runtime.GOOS != "darwin" {
		return path
	}
	if strings.Contains(path, "/dev/tty.") {
		return strings.Replace(path, "/dev/tty.", "/dev/cu.", 1)
	}
	return path
}

func killSystemModemManager() {
	exec.Command("killall", "-HUP", "ModemManager").Run()
	exec.Command("networksetup", "-setmodemconnectiondisabled", "on").Run()
}
