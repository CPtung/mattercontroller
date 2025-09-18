package simulation

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/CPtung/mtctrl/internal/config"
)

// isNumeric
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func StartRCP() (string, error) {
	ttyNumber := ""
	cmd := exec.Command(filepath.Join(config.LibPath, "start_socat.sh"))
	output, err := cmd.Output()
	if err != nil {
		return ttyNumber, fmt.Errorf("執行 start_socat.sh 失敗: %v", err)
	}
	ttyNumber = strings.TrimSpace(string(output))

	fmt.Printf("✅ start_socat.sh 執行成功\n")
	fmt.Printf("💾 PTS 編號變數: %s\n", ttyNumber)

	// 驗證結果是否為數字
	if ttyNumber != "" && isNumeric(ttyNumber) {
		fmt.Printf("🎯 PTY 完整路徑: /dev/pts/%s\n", ttyNumber)
	} else {
		fmt.Printf("⚠️  輸出格式可能異常: '%s'\n", ttyNumber)
	}
	return ttyNumber, nil
}

func StopRCP() error {
	cmd := exec.Command(filepath.Join(config.LibPath, "stop_socat.sh"))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("執行 stop_socat.sh 失敗: %v", err)
	}
	fmt.Printf("✅ stop_socat.sh 執行成功\n")
	fmt.Printf("💾 執行結果: %s\n", string(output))
	return nil
}
