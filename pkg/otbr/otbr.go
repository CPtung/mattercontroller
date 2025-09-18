package otbr

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	command     string
	description string
}

// executeCommand 執行docker命令並返回輸出
func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// checkOutput 檢查輸出是否符合預期
func checkOutput(output, expected string) bool {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == expected {
			return true
		}
	}
	return false
}

// executeOTCommand 執行OpenThread命令並驗證結果
func executeOTCommand(otCommand, description string) bool {
	dockerCmd := fmt.Sprintf(`docker exec -i otbr ot-ctl %s`, otCommand)

	fmt.Printf("執行: %s\n", description)
	output, err := executeCommand(dockerCmd)

	if err != nil {
		fmt.Printf("❌ 命令執行失敗: %v\n", err)
		return false
	}

	fmt.Printf("輸出:\n%s\n", output)

	/*
		if checkOutput(output, "Done") {
			fmt.Printf("✅ %s 成功\n", description)
			return true
		} else {
			fmt.Printf("❌ %s 失敗 - 未找到 'Done'\n", description)
			return false
		}
	*/
	return true
}

// checkLeaderStatus 檢查是否為leader狀態
func checkStatus(state string) bool {
	dockerCmd := `docker exec -i otbr ot-ctl state`

	fmt.Println("檢查 State 狀態...")
	output, err := executeCommand(dockerCmd)

	if err != nil {
		fmt.Printf("❌ 狀態檢查失敗: %v\n", err)
		return false
	}

	fmt.Printf("狀態輸出:\n%s\n", output)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "leader" {
			fmt.Println("✅ 已成為 leader")
			return true
		} else if line == "disabled" {
			fmt.Println("✅ 網絡已停用")
			return true
		}
	}

	fmt.Printf("❌ 尚未成為 %s\n", state)
	return false
}

func writeTtyNum(ttyNum string) error {
	// write ttyNum to /etc/mtctrl/otbr-env.list
	envPath := "/run/mtctrl/otbr-env.list"
	content := fmt.Sprintf(`OT_RCP_DEVICE=spinel+hdlc+uart:///dev/pts/%s?uart-baudrate=115200
OT_INFRA_IF=enp0s31f6
OT_THREAD_IF=wpan0
OT_LOG_LEVEL=6`, ttyNum)

	return os.WriteFile(envPath, []byte(content), 0644)
}

func CreateBorderRouter(ttyNum string) error {
	// write ttyNum to /run/mtctrl/otbr-env.list
	err := writeTtyNum(ttyNum)
	if err != nil {
		return err
	}

	// docker run border router container
	dockerCmd := `docker run --name otbr -d --rm \
	--cap-add=net_admin \
	--env-file=/run/mtctrl/otbr-env.list \
	--network=host \
	-v /dev/pts:/dev/pts \
	--device=/dev/net/tun \
	--volume=/var/lib/otbr:/data \
	openthread/border-router:latest`
	fmt.Printf("執行: %s\n", dockerCmd)
	output, err := executeCommand(dockerCmd)
	if err != nil {
		fmt.Printf("❌ 命令執行失敗: %s\n", output)
		return err
	}
	// retry 10 times and interval 1 second
	for i := 0; i < 10; i++ {
		statsCmd := `docker stats otbr --no-stream --format "{{.CPUPerc}}"`
		output, err = executeCommand(statsCmd)
		if err == nil {
			fmt.Println("✅ 啟動 Border Router 成功")
			return nil
		}
		fmt.Printf("❌ 命令執行失敗: %s\n", output)
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("啟動 Border Router 失敗")
}

func CloseBorderRouter() error {
	// docker stop border router container
	dockerCmd := `docker stop otbr`
	fmt.Printf("執行: %s\n", dockerCmd)
	output, err := executeCommand(dockerCmd)
	if err != nil {
		fmt.Printf("❌ 命令執行失敗: %s\n", output)
		return err
	}
	return nil
}

// executeCommandSequence 執行一系列 OpenThread 命令
func executeCommandSequence(commands []Command, sequenceName string) error {
	fmt.Printf("開始 %s...\n", sequenceName)

	for i, cmd := range commands {
		success := executeOTCommand(cmd.command, cmd.description)

		if !success {
			return fmt.Errorf("命令執行失敗: %s", cmd.description)
		}

		// 每個命令之間間隔1秒（最後一個命令後不需要等待）
		if i < len(commands)-1 {
			fmt.Println("等待 1 秒...")
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("\n%s命令完成\n", sequenceName)
	return nil
}

// waitForStatus 等待指定狀態，可配置檢查間隔和超時
func waitForStatus(targetStatus, statusName string, maxAttempts int, checkInterval time.Duration) error {
	fmt.Printf("開始檢查 %s 狀態...\n", statusName)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("\n第 %d 次狀態檢查:\n", attempt)

		if checkStatus(targetStatus) {
			fmt.Printf("\n🎉 狀態檢查完成！設備已成為 %s\n", targetStatus)
			return nil
		}

		if attempt < maxAttempts {
			fmt.Printf("等待 %.0f 秒後重新檢查...\n", checkInterval.Seconds())
			time.Sleep(checkInterval)
		}
	}

	totalTime := time.Duration(maxAttempts) * checkInterval
	fmt.Printf("\n⚠️  經過 %d 次檢查（%.0f 秒），設備仍未成為 %s\n",
		maxAttempts, totalTime.Seconds(), targetStatus)
	return fmt.Errorf("狀態檢查失敗: 未能達到 %s 狀態", targetStatus)
}

// SetupThreadNetwork 設置 Thread 網絡
func SetupThreadNetwork() error {
	commands := []Command{
		{"dataset init new", "初始化新數據集"},
		{"dataset commit active", "提交活動數據集"},
		{"ifconfig up", "啟用網絡接口"},
		{"thread start", "啟動 Thread 網絡"},
	}

	// 執行命令序列
	if err := executeCommandSequence(commands, "OpenThread 初始化流程"); err != nil {
		log.Printf("初始化失敗: %v", err)
		return err
	}

	// 等待 leader 狀態
	return waitForStatus("leader", "leader", 10, 5*time.Second)
}

// TearDownThreadNetwork 關閉 Thread 網絡
func TearDownThreadNetwork() error {
	commands := []Command{
		{"thread stop", "停止 Thread 網絡"},
		{"factoryreset", "重置初始化"},
	}

	// 執行命令序列
	if err := executeCommandSequence(commands, "OpenThread 關閉流程"); err != nil {
		log.Printf("關閉失敗: %v", err)
		return err
	}

	// 等待 disabled 狀態
	return waitForStatus("disabled", "disabled", 10, 5*time.Second)
}
