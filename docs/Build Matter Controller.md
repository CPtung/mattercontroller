# Build Matter Controller

# 環境

- x64_86 Debian 11

# 前置作業

Matter Controller是基於 Go 撰寫開發的Service

## 1. 建置編譯環境

從 Go 官方 Docker repository 下載 `linux/amd64` 版本 `debian:bullseye` 

```bash
docker pull --platform=linux/amd64 golang:bullseye
```

啟動容器

```bash
docker create --platform=linux/amd64 \
		--name matter-builder \
		-v $PWD:/go/src \
		-w /go/src \
		-it golang:bullseye
		
docker start -ai matter-builder
```

## 2. Project Structure

```
.
├── cmd
│   ├── chip.go
│   ├── daemon.go
│   ├── light.go
│   └── main.go
├── internal
│   ├── config
│   ├── database
│   ├── matter
│   ├── otbr
│   └── server
└── pkg
    ├── model
    └── restapi
```

- cmd: 包含各種應用程式元件的進入點：
    - main.go：主要應用程式進入點
    - daemon.go：背景服務控制操作
    - chip.go：Matter/chip-tool相關命令
    - light.go：lighting app相關指令
- internal:
    - config: service相關的filepath 以及 config/scripts template
    - database: matter device資料庫實作
    - matter: 包含 chip-tool 以及 lighting app 相關指令實作
    - otbr: 包含rcp 以及 border router 相關指令實作
    - server: 基於 gin webframework 的 API server 相關實作
- pkg:
    - model: 各 Object 定義
    - restapi: Service所支援REST API 相關實作

# 編譯 Binary

細節可參考 `Makefile`

```
go build \
	-buildmode=pie \
	-ldflags "-s -w -X main.Version=4dea4a9-dirty -X main.BuildTime=2025-09-18T13:55:06Z" \
	-o build/matter ./cmd
```