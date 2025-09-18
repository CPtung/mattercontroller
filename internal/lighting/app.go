package lighting

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type LightingAppController struct {
	BinPath     string
	InterfaceID string
	StopTimeout time.Duration
	mu          sync.Mutex
	cmd         *exec.Cmd
	startedAt   time.Time
}

// NewLightingAppController 建構器
func NewLightingAppController(binPath, iface string) *LightingAppController {
	if binPath == "" {
		binPath = "/usr/local/bin/chip-lighting-app"
	}
	return &LightingAppController{
		BinPath:     binPath,
		InterfaceID: iface,
		StopTimeout: 5 * time.Second,
	}
}

// Start 如果尚未啟動，啟動外部程式
func (c *LightingAppController) Start(ctx context.Context) error {
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	c.cmd = cmd
	return nil
}

// Stop 優雅停止（TERM），逾時後強制（KILL）
func (c *LightingAppController) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd == nil || c.cmd.Process == nil {
		return errors.New("lighting app not running")
	}

	pid := c.cmd.Process.Pid
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM: %w", err)
	}

	// 等候優雅結束
	done := make(chan error, 1)
	go func(cmd *exec.Cmd) {
		done <- cmd.Wait()
	}(c.cmd)

	select {
	case <-ctx.Done():
	case <-time.After(c.StopTimeout):
	case err := <-done:
		if err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Printf("lighting app exited with error: %v", err)
		}
		c.cmd = nil
		return nil
	}

	// 強制 KILL
	_ = syscall.Kill(pid, syscall.SIGKILL)
	c.cmd = nil
	return nil
}

func (c *LightingAppController) Status() string {
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

func (c *LightingAppController) PostLightingStartStop(g *gin.Context) {
	action := g.Param("action")
	ctx, cancel := context.WithTimeout(g.Request.Context(), c.StopTimeout+3*time.Second)
	defer cancel()

	var err error
	switch action {
	case "start":
		err = c.Start(ctx)
		if err == nil {
			g.JSON(http.StatusOK, gin.H{"status": "started", "interfaceId": c.InterfaceID})
			return
		}
	case "stop":
		err = c.Stop(ctx)
		if err == nil {
			g.JSON(http.StatusOK, gin.H{"status": "stopped"})
			return
		}
	default:
		g.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid action"})
		return
	}

	// 錯誤情況
	g.JSON(http.StatusConflict, gin.H{"status": "error", "error": err.Error()})
}

func (c *LightingAppController) GetLightingStatus(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{
		"status":      c.Status(),
		"interfaceId": c.InterfaceID,
	})
}
