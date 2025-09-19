package main

import (
	"fmt"

	"github.com/CPtung/mattercontroller/internal/matter/chiptool"
	"github.com/CPtung/mattercontroller/pkg/model"
	"github.com/spf13/cobra"
)

var (
	chipCmd = &cobra.Command{
		Use:   "chip",
		Short: "Run Chip tool",
	}

	pairCmd = &cobra.Command{
		Use:   "pairing <nodeID> <pairCode>",
		Short: "Pairing a node",
		Args:  cobra.ExactArgs(2),
		RunE:  pairing,
	}

	unpairCmd = &cobra.Command{
		Use:   "unpairing <deviceID>",
		Short: "Unpairing a node",
		Args:  cobra.ExactArgs(1),
		RunE:  unpairing,
	}

	onCmd = &cobra.Command{
		Use:   "on <deviceID>",
		Short: "Light on",
		Args:  cobra.ExactArgs(1),
		RunE:  lightOn,
	}

	offCmd = &cobra.Command{
		Use:   "off <deviceID>",
		Short: "Light off",
		Args:  cobra.ExactArgs(1),
		RunE:  lightOff,
	}

	stateCmd = &cobra.Command{
		Use:   "state <deviceID>",
		Short: "Get state",
		Args:  cobra.ExactArgs(1),
		RunE:  getState,
	}
)

func init() {
	chipCmd.AddCommand(pairCmd)
	chipCmd.AddCommand(unpairCmd)
	chipCmd.AddCommand(onCmd)
	chipCmd.AddCommand(offCmd)
	chipCmd.AddCommand(stateCmd)
	rootCmd.AddCommand(chipCmd)
}

func pairing(cmd *cobra.Command, args []string) error {
	nodeID := args[0]
	pairCode := args[1]

	// 建立 Matter 控制器
	controller := chiptool.New(cmd.Context(), "")

	fmt.Println("=== Matter 設備控制器範例 ===")

	// 範例 1: 配對設備
	fmt.Println("\n1. 配對設備範例:")
	resp, err := controller.PairDevice(nodeID, pairCode)
	if err == nil {
		fmt.Printf("配對結果: %+v\n", resp.Success)
		fmt.Printf("Device ID: %d\n", resp.Data.(*model.MatterDevice).ID)
	}
	return err
}

func unpairing(cmd *cobra.Command, args []string) error {
	nodeID := args[0]

	// 建立 Matter 控制器
	controller := chiptool.New(cmd.Context(), "")

	fmt.Println("=== Matter 設備控制器範例 ===")

	// 範例 1: 配對設備
	fmt.Println("\n1. 解除配對設備範例:")
	resp, err := controller.UnpairDevice(nodeID)
	if err == nil {
		fmt.Printf("解除配對結果: %+v\n", resp.Success)
	}
	return err
}

func lightOn(cmd *cobra.Command, args []string) error {
	deviceID := args[0]

	// 建立 Matter 控制器
	controller := chiptool.New(cmd.Context(), "")

	fmt.Println("=== Matter 設備控制器範例 ===")

	// 範例 2: 控制開關
	fmt.Println("\n2. 控制開關範例:")
	controller.TurnOn(deviceID)

	// 範例 3: 讀取狀態
	fmt.Println("\n3. 讀取狀態範例:")
	if resp, err := controller.GetOnOffState(deviceID); err == nil && resp.Success {
		data := resp.Data.(*model.MatterLightConfig)
		fmt.Printf("當前狀態: %v\n", data.State)
	}

	return nil
}

func lightOff(cmd *cobra.Command, args []string) error {
	deviceID := args[0]

	// 建立 Matter 控制器
	controller := chiptool.New(cmd.Context(), "")

	fmt.Println("=== Matter 設備控制器範例 ===")

	// 範例 2: 控制開關
	controller.TurnOff(deviceID)

	// 範例 3: 讀取狀態
	fmt.Println("\n3. 讀取狀態範例:")
	if resp, err := controller.GetOnOffState(deviceID); err == nil && resp.Success {
		data := resp.Data.(*model.MatterLightConfig)
		fmt.Printf("當前狀態: %v\n", data.State)
	}
	return nil
}

func getState(cmd *cobra.Command, args []string) error {
	deviceID := args[0]

	// 建立 Matter 控制器
	controller := chiptool.New(cmd.Context(), "")

	fmt.Println("=== Matter 設備控制器範例 ===")

	// 範例 3: 讀取狀態
	fmt.Println("\n3. 讀取狀態範例:")
	if resp, err := controller.GetOnOffState(deviceID); err == nil && resp.Success {
		data := resp.Data.(*model.MatterLightConfig)
		fmt.Printf("當前狀態: %v\n", data.State)
	}
	return nil
}
