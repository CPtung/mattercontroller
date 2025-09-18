package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

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
	dockerCmd := fmt.Sprintf(`docker exec -i obtr ot-ctl %s`, otCommand)

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
func checkLeaderStatus() bool {
	dockerCmd := `docker exec -i obtr ot-ctl state`

	fmt.Println("檢查 leader 狀態...")
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
		}
	}

	fmt.Println("❌ 尚未成為 leader")
	return false
}

func main() {
	fmt.Println("開始 OpenThread 初始化流程...")

	// 定義要執行的命令
	commands := []struct {
		command     string
		description string
	}{
		{"dataset init new", "初始化新數據集"},
		{"dataset commit active", "提交活動數據集"},
		{"ifconfig up", "啟用網絡接口"},
		{"thread start", "啟動 Thread 網絡"},
	}

	// 執行初始化命令
	for i, cmd := range commands {
		success := executeOTCommand(cmd.command, cmd.description)

		if !success {
			log.Fatalf("命令執行失敗，停止執行: %s", cmd.description)
		}

		// 每個命令之間間隔1秒（最後一個命令後不需要等待）
		if i < len(commands)-1 {
			fmt.Println("等待 1 秒...")
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("\n初始化命令完成，開始檢查 leader 狀態...")

	// 檢查leader狀態，每3秒檢查一次，最多檢查10次（30秒）
	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("\n第 %d 次狀態檢查:\n", attempt)

		if checkLeaderStatus() {
			fmt.Println("\n🎉 OpenThread 網絡初始化完成！設備已成為 leader")
			return
		}

		if attempt < maxAttempts {
			fmt.Println("等待 3 秒後重新檢查...")
			time.Sleep(3 * time.Second)
		}
	}

	fmt.Printf("\n⚠️  經過 %d 次檢查（%d 秒），設備仍未成為 leader\n", maxAttempts, maxAttempts*3)
}
