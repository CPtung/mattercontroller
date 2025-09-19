package model

import "time"

// MatterDevice 代表一個 Matter 設備
type MatterDevice struct {
	ID        uint      `structs:"-" gorm:"primaryKey;not null;" json:"id"`
	CreatedAt time.Time `structs:"-" json:"-"`
	UpdatedAt time.Time `structs:"-" json:"-"`

	NodeID     string `json:"nodeId"`
	EndpointID string `json:"endpointId"`
}

// CommandResponse 代表命令回應
type CommandResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	RawOutput string      `json:"rawOutput"`
}

// OnOffState 開關狀態
type OnOffState struct {
	IsOn bool `json:"isOn"`
}

type MatterPairReqeust struct {
	NodeID   string `json:"nodeID" binding:"required"`
	PairCode string `json:"pairCode" binding:"required"`
}

type MatterUnpairingReqeust struct {
	DeviceID string `json:"deviceID" binding:"required"`
}

type MatterLightCmdRequest struct {
	NodeID     string `json:"nodeID" binding:"required"`
	EndpointID string `json:"endpointID" binding:"required"`
}

type MatterLightConfig struct {
	State string `json:"state" binding:"required,oneof=on off"`
}
