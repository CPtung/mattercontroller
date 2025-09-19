package chiptool

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/CPtung/mattercontroller/internal/database"
	"github.com/CPtung/mattercontroller/pkg/model"
)

type MatterController interface {
	// PairDevice
	PairDevice(nodeID, pairingCode string) (*model.CommandResponse, error)
	// UnpairDevice
	UnpairDevice(deviceID string) (*model.CommandResponse, error)
	// TurnOn
	TurnOn(deviceID string) (*model.CommandResponse, error)
	// TurnOff
	TurnOff(deviceID string) (*model.CommandResponse, error)
	// GetOnOffState
	GetOnOffState(deviceID string) (*model.CommandResponse, error)
}

// MatterController Matter
type MatterControllerImpl struct {
	chipToolPath string
	timeout      time.Duration
}

var (
	boolRegex   = regexp.MustCompile(`(?i)(Data|BOOL)\s*=\s*(true|false|1|0)`)
	stringRegex = regexp.MustCompile(`Data\s*=\s*"([^"]*)"`)
)

// NewMatterController 建立新的 Matter 控制器
func New(ctx context.Context, chipToolPath string) MatterController {
	if chipToolPath == "" {
		chipToolPath = "/usr/local/bin/chip-tool"
	}

	return &MatterControllerImpl{
		chipToolPath: chipToolPath,
		timeout:      30 * time.Second,
	}
}

// executeCommand 執行 chip-tool 命令
func (mc *MatterControllerImpl) executeCommand(args ...string) (*model.CommandResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, mc.chipToolPath, args...)
	output, err := cmd.CombinedOutput()

	response := &model.CommandResponse{
		RawOutput: string(output),
	}

	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("命令執行失敗: %v", err)
		return response, err
	}

	response.Success = true
	return response, nil
}

// PairDevice 配對設備
func (mc *MatterControllerImpl) PairDevice(nodeID, pairingCode string) (*model.CommandResponse, error) {
	fmt.Printf("開始配對設備 (Node ID: %s, Pairing Code: %s)\n", nodeID, pairingCode)

	response, err := mc.executeCommand("pairing", "code", nodeID, pairingCode)
	if !response.Success {
		return response, fmt.Errorf("配對失敗: err: %s, resp: %s", err.Error(), response.Error)
	}
	// Assign a mock device ID
	data := model.MatterDevice{
		NodeID:     nodeID,
		EndpointID: "1",
	}
	response.Data = &data

	// mc.devices.Store(data.DeviceID, data)
	database.Store(&data)
	return response, nil
}

// UnpairDevice 解除配對設備
func (mc *MatterControllerImpl) UnpairDevice(deviceID string) (*model.CommandResponse, error) {
	fmt.Printf("解除配對設備 (Device ID: %s)\n", deviceID)
	device, err := database.Load(deviceID)
	if err != nil {
		return nil, fmt.Errorf("找不到設備 %s", err.Error())
	}

	resp, err := mc.executeCommand("pairing", "unpair", device.NodeID)
	if err == nil && resp.Success {
		database.Delete(deviceID)
	}
	return resp, err
}

// TurnOn 開燈
func (mc *MatterControllerImpl) TurnOn(deviceID string) (*model.CommandResponse, error) {

	device, err := database.Load(deviceID)
	if err != nil {
		return nil, fmt.Errorf("找不到設備")
	}

	fmt.Printf("開燈 (Node ID: %s, Endpoint: %s)\n", device.NodeID, device.EndpointID)
	resp, err := mc.executeCommand("onoff", "on", device.NodeID, device.EndpointID)
	if err != nil {
		return resp, err
	}
	resp.Data = &model.MatterLightConfig{State: "on"}
	return resp, nil
}

// TurnOff 關燈
func (mc *MatterControllerImpl) TurnOff(deviceID string) (*model.CommandResponse, error) {
	device, err := database.Load(deviceID)
	if err != nil {
		return nil, fmt.Errorf("找不到設備")
	}

	fmt.Printf("關燈 (Node ID: %s, Endpoint: %s)\n", device.NodeID, device.EndpointID)
	resp, err := mc.executeCommand("onoff", "off", device.NodeID, device.EndpointID)
	if err != nil {
		return resp, err
	}
	resp.Data = &model.MatterLightConfig{State: "off"}
	return resp, nil
}

// GetOnOffState 讀取開關狀態
func (mc *MatterControllerImpl) GetOnOffState(deviceID string) (*model.CommandResponse, error) {
	device, err := database.Load(deviceID)
	if err != nil {
		return nil, fmt.Errorf("找不到設備")
	}

	response, err := mc.executeCommand("onoff", "read", "on-off", device.NodeID, device.EndpointID)
	if err != nil {
		return response, err
	}

	if response.Success {
		isOn := mc.parseOnOffState(response.RawOutput)
		response.Data = &model.MatterLightConfig{State: func() string {
			if isOn {
				return "on"
			}
			return "off"
		}()}
	}

	return response, nil
}

// parseOnOffState 解析開關狀態
func (mc *MatterControllerImpl) parseOnOffState(output string) bool {
	// 尋找 "Data = true" 或 "Data = false" 或 "BOOL = true/false"
	matches := boolRegex.FindStringSubmatch(output)

	if len(matches) >= 3 {
		value := strings.ToLower(matches[2])
		return value == "true" || value == "1"
	}

	return false
}
