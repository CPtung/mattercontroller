package lighting

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"syscall"
	"time"
)

type LightingApp interface {
	Start(ctx context.Context) error
	Stop() error
	PairingCode() string
	Status() string
}

type LightingAppImpl struct {
	BinPath     string
	InterfaceID string
	StopTimeout time.Duration
	pairingCode string
	mu          sync.Mutex
	cmd         *exec.Cmd
	startedAt   time.Time
}

var (
	rePairingCode = regexp.MustCompile(`Manual pairing code: \[(\d+)\]`)
	rePaired      = regexp.MustCompile(`Released\s-\sType:1`)
	reUnpaired    = regexp.MustCompile(`Released\s-\sType:2`)
)

// NewLightingAppController 建構器
func NewApp(binPath, iface string) LightingApp {
	if binPath == "" {
		binPath = "/usr/local/bin/chip-lighting-app"
	}
	return &LightingAppImpl{
		BinPath:     binPath,
		InterfaceID: iface,
		StopTimeout: 5 * time.Second,
	}
}

// Start
func (c *LightingAppImpl) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.InterfaceID == "" {
		return errors.New("InterfaceID is empty")
	}

	if c.cmd != nil && c.cmd.Process != nil {
		// terminate if previous process still alive
		if err := c.cmd.Process.Signal(syscall.Signal(0)); err == nil {
			return errors.New("lighting app already running")
		}
		c.cmd = nil
	}

	if _, err := os.Stat(c.BinPath); err != nil {
		return fmt.Errorf("binary not found: %s (%w)", c.BinPath, err)
	}

	args := []string{"--interface-id", c.InterfaceID}
	cmd := exec.CommandContext(ctx, c.BinPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	// 處理 stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			//fmt.Printf("[STDOUT] %s\n", line)

			// 擷取 pairing code
			if c.pairingCode == "" {
				if matches := rePairingCode.FindStringSubmatch(line); len(matches) > 1 {
					c.pairingCode = matches[1]
					fmt.Printf("\n✅ 找到 Manual pairing code: %s\n", c.pairingCode)
				}
			} else if rePaired.MatchString(line) {
				fmt.Printf("\n✅ 配對成功\n")
			} else if reUnpaired.MatchString(line) {
				c.pairingCode = ""
				fmt.Printf("\n✅ 解除成功\n")
			}
			if err := ctx.Err(); err != nil {
				return
			}
		}
	}()

	c.cmd = cmd
	return nil
}

// Stop
func (c *LightingAppImpl) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd == nil || c.cmd.Process == nil {
		return errors.New("lighting app not running")
	}

	pid := c.cmd.Process.Pid
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM: %w", err)
	}

	done := make(chan error, 1)
	go func(cmd *exec.Cmd) {
		done <- cmd.Wait()
	}(c.cmd)

	select {
	case <-time.After(3 * time.Second):
	case err := <-done:
		if err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Printf("lighting app exited with error: %v", err)
		}
		c.cmd = nil
		return nil
	}

	// Force kill
	_ = syscall.Kill(pid, syscall.SIGKILL)
	c.cmd = nil
	return nil
}

func (c *LightingAppImpl) PairingCode() string {
	return c.pairingCode
}

func (c *LightingAppImpl) Status() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd == nil || c.cmd.Process == nil {
		return "stopped"
	}
	if err := c.cmd.Process.Signal(syscall.Signal(0)); err == nil {
		return "running"
	}
	return "stopped"
}
