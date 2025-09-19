package config

import (
	"fmt"
	"os"
	"sync"
)

const (
	LibPath     = "/var/lib/matter"
	EtcPath     = "/etc/matter"
	RuntimePath = "/run/matter"
)

var oncePaths sync.Once

func createServicePaths() {
	oncePaths.Do(func() {
		if _, err := os.Stat(LibPath); os.IsNotExist(err) {
			if err := os.MkdirAll(LibPath, 0755); err != nil {
				fmt.Printf("")
			}
			// Create default rcp scripts
			if err := createScripts(); err != nil {
				fmt.Println(err.Error())
			}
		}
		if _, err := os.Stat(EtcPath); os.IsNotExist(err) {
			if err := os.MkdirAll(EtcPath, 0755); err != nil {
				fmt.Printf("")
			}
		}
		if _, err := os.Stat(RuntimePath); os.IsNotExist(err) {
			if err := os.MkdirAll(RuntimePath, 0755); err != nil {
				fmt.Printf("")
			}
		}
	})
}
