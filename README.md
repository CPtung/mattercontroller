# Matter Controller Demo

## 目錄

- [Demo 操作流程](#demo-操作流程)
- [CLI 使用方式](#cli-使用方式)
- [Matter & Light API 使用說明](#matter--light-api-使用說明)
    - [伺服器資訊](#伺服器資訊)
    - [範例操作](#範例操作)
        - [1. 配對裝置 (Pairing)](#1-配對裝置-pairing)
        - [2. 解除配對 (Unpairing)](#2-解除配對-unpairing)
        - [3. 查詢燈狀態 (Get Light State)](#3-查詢燈狀態-get-light-state)
        - [4. 控制燈開關 (Set Light State)](#4-控制燈開關-set-light-state)
    - [錯誤回應](#錯誤回應)
- [進一步閱讀](#進一步閱讀)

## 文件簡介

本文件說明如何啟動 **Matter Controller Demo**，並介紹相關 CLI 使用方式。

如果已了解操作，可直接轉跳至底下 [進一步閱讀](#進一步閱讀)。

## Demo 操作流程

1. 先將編譯好的 `chip-tool` `chip-lighting-app` 放置到系統 PATH 中。
    
    ```bash
    # 如何編譯，可以參考 進一步閱讀中 的建置方法
    cp chip-tool chip-lighting-app /usr/local/bin
    ```
    
2. 同時也將編譯好的 matter controller 的 binary 檔一樣放置到 系統 PATH 中，並將matter 的 systemd config 安置到系統中。
    
    ```bash
    # 如何編譯，可以參考 Makefile
    cp matter /usr/local/bin/
    
    # matter.service 檔案可參考 data 資料夾
    cp matter.service /lib/systemd/system/
    systemctl daemon-reload
    ```
    
3. 啟動 OpenThread 模擬網路，當看到 `🎉 狀態檢查完成！設備已成為 leader` 表示網路已經建立。
    
    ```bash
    # 複製OTBR ENV檔到指定位置
    cp otbr-env.list /var/run/matter/
    systemctl start matter
    
    : ✅ start_socat.sh 執行成功
    : 💾 PTS 編號變數: 3
    : 🎯 PTY 完整路徑: /dev/pts/3
    : 執行: docker run --name otbr -d --rm \
    : 	--cap-add=net_admin \
    : 	--env-file=/run/matter/otbr-env.list \
    : 	--network=host \
    : 	-v /dev/pts:/dev/pts \
    : 	--device=/dev/net/tun \
    : 	--volume=/var/lib/otbr:/data \
    : 	openthread/border-router:latest
    : ✅ 啟動 Border Router 成功
    : 開始 OpenThread 初始化流程...
    : 執行: 初始化新數據集
    : 輸出:
    : Done
    : 等待 1 秒...
    : 執行: 提交活動數據集
    : 輸出:
    : Done
    : 等待 1 秒...
    : 執行: 啟用網絡接口
    : 輸出:
    : Done
    : 等待 1 秒...
    : 執行: 啟動 Thread 網絡
    : 輸出:
    : Done
    : OpenThread 初始化流程命令完成
    : 開始檢查 leader 狀態...
    : 第 1 次狀態檢查:
    : 檢查 State 狀態...
    : 狀態輸出:
    : detached
    : Done
    : ❌ 尚未成為 leader
    : 等待 5 秒後重新檢查...
    : 第 2 次狀態檢查:
    : 檢查 State 狀態...
    : 狀態輸出:
    : leader
    : Done
    : ✅ 已成為 leader
    : 🎉 狀態檢查完成！設備已成為 leader
    ```
    
4. 找到網路名稱，並啟動 Lighting App，取得配對碼
    
    ```bash
    # 此 Demo 預設 enp0s31f6，可到 otbr-env.list中修改
    ip -brief a | grep UP
    enp0s31f6        UP             10.123.13.105/23 fe80::290:e8ff:fea6:61a2/64
    
    matter light enp0s31f6
    
    ✅ 找到 Manual pairing code: 34970112332
    ```
    
5. 開啟另一個視窗，執行配對 (Pairing)
    
    ```bash
    # 此 Demo nodeID 為 1
    matter chip pairing 1 34970112332
    === Matter 設備控制器範例 ===
    
    1. 配對設備範例:
    開始配對設備 (Node ID: 1, Pairing Code: 34970112332)
    配對結果: true
    Device ID: 8
    
    --------------------------------------------
    
    # light app 的視窗 會顯示成功訊息
    ✅ 配對成功
    ```
    
6. 驗證裝置狀態
    
    ```bash
    matter light state 8
    
    === Matter 設備控制器範例 ===
    
    3. 讀取狀態範例:
    當前狀態: off
    ```
    
    `<deviceID>` 為配對成功後由 Matter Controller 回傳的裝置 ID，此處為 `8`。
    
7. **控制燈具**
    - 開燈：
        
        ```bash
        matter chip on <deviceID>
        
        === Matter 設備控制器範例 ===
        
        2. 控制開關範例:
        開燈 (Node ID: 1, Endpoint: 1)
        
        3. 讀取狀態範例:
        當前狀態: on
        ```
        
    - 關燈：
        
        ```bash
        matter chip off <deviceID>
        
        === Matter 設備控制器範例 ===
        關燈 (Node ID: 1, Endpoint: 1)
        
        3. 讀取狀態範例:
        當前狀態: off
        
        ```
        
8. 執行解除配對
    
    ```bash
    # 此 Demo deviceID 為 8
    matter chip pairing 8
    === Matter 設備控制器範例 ===
    
    1. 解除配對設備範例:
    解除配對設備 (Device ID: 8)
    解除配對結果: true
    
    --------------------------------------------
    
    # light app 的視窗 會顯示成功訊息
    ✅ 解除成功
    ```
    
9. 結束 OpenThread 模擬網路
    
    ```bash
    systemctl stop matter
    
    : Stopping Matter Controller Service...
    : 開始 OpenThread 關閉流程...
    : 執行: 停止 Thread 網絡
    : 輸出:
    : Done
    : 等待 1 秒...
    : 執行: 重置初始化
    : 輸出:
    : OpenThread 關閉流程命令完成
    : 開始檢查 disabled 狀態...
    : 第 1 次狀態檢查:
    : 檢查 State 狀態...
    : 狀態輸出:
    : disabled
    : Done
    : ✅ 網絡已停用
    : 🎉 狀態檢查完成！設備已成為 disabled
    : 執行: docker stop otbr
    : [1]00:05:59.027 [C] Platform------: platformUartProcess() at uart.c:238: Invalid argument
    : ✅ stop_socat.sh 執行成功
    : 💾 執行結果:
    : matter.service: Succeeded.
    : Stopped Matter Controller Service.
    
    ```
    

---

## CLI 使用方式

透過 `matter` 指令，可以操作多種功能，以下為主要子指令：

### 基本用法

```bash
matter [command]

Matter Controller

Usage:
  matter [command]

Available Commands:
  chip        Run Chip tool
  completion  Generate the autocompletion script for the specified shell
  daemon      Run daemon
  help        Help about any command
  info        Show version/build info
  light       Lighting App

Flags:
  -h, --help   help for matter

Use "matter [command] --help" for more information about a command.

```

### 可用指令

- **chip**
    
    執行 CHIP 工具，支援配對相關操作。
    
    - 範例：`matter chip pairing <nodeID> <pairCode>`
- **light**
    
    控制與查詢燈具狀態。
    
    - 範例：
        - `matter light <deviceID> on`
        - `matter light <deviceID> off`
        - `matter light <deviceID> state`
- **daemon**
    
    以背景服務模式啟動 Matter Controller。
    
- **info**
    
    顯示版本與編譯資訊。
    
- **help**
    
    顯示指令的幫助說明。
    

> 可使用 matter [command] --help 取得特定指令的詳細用法。
> 


---

## Matter & Light API 使用說明

在完成上述步驟 4. 時，後續 chip 指令可以透過 REST API 達成 (步驟 5.~ 8.)

Web Server 透過 **Unix Domain Socket** 提供服務，路徑為：

```
/var/run/matter/matter.sock
```

在 `curl` 中呼叫時，需要使用 `--unix-socket` 參數指定該 socket。
OpenAPI 文件可參考 [openapi/openapi.yml](openapi/openapi.yml)。

---

### 伺服器資訊

- 協定：`http+unix`
- Socket 路徑：`/var/run/matter/matter.sock`
- Base URL：`http://localhost`

---

### 範例操作

### 1. 配對裝置 (Pairing)

**Endpoint**

`POST /matter/pairing`

**Request 範例**

```bash
curl --unix-socket /var/run/matter/matter.sock \
  -X POST http://localhost/matter/pairing \
  -H "Content-Type: application/json" \
  -d '{
    "nodeID": "1",
    "pairCode": "34970112332"
  }'

```

**成功回應範例**

```json
{
  "id": 1,
  "nodeId": "1",
  "endpointId": "1"
}
```

---

### 2. 解除配對 (Unpairing)

**Endpoint**

`POST /matter/unpairing/{deviceID}`

**Request 範例**

```bash
curl --unix-socket /var/run/matter/matter.sock \
  -X POST http://localhost/matter/unpairing/1
```

**成功回應範例**

```json
{ "success": true }
```

---

### 3. 查詢燈狀態 (Get Light State)

**Endpoint**

`GET /light/{deviceID}`

**Request 範例**

```bash
curl --unix-socket /var/run/matter/matter.sock \
  http://localhost/light/1
```

**成功回應範例**

```json
{ "state": "on" }
```

---

### 4. 控制燈開關 (Set Light State)

**Endpoint**

`PUT /light/{deviceID}`

**Request 範例：開燈**

```bash
curl --unix-socket /var/run/matter/matter.sock \
  -X PUT http://localhost/light/1 \
  -H "Content-Type: application/json" \
  -d '{"state": "on"}'
```

**Request 範例：關燈**

```bash
curl --unix-socket /var/run/matter/matter.sock \
  -X PUT http://localhost/light/1 \
  -H "Content-Type: application/json" \
  -d '{"state": "off"}'
```

**成功回應範例**

```json
{ "state": "on" }
```

---

## 錯誤回應

常見錯誤格式如下：

```json
{ "error": "error message" }
```

- **400**：輸入參數驗證失敗
- **404**：裝置不存在
- **500**：內部伺服器錯誤

---

## 進一步閱讀

更完整的架構設計與安裝設定文件，請參考 docs/ 目錄：
- [Architecture Overview](docs/Architecture%20Overview.md)
- [Build Chip-tool and Lighting app](docs/Build%20Chip-tool%20and%20Lighting%20app.md)
- [Setup OpenThread Border Router](docs/Setup%20OpenThread%20Border%20Router.md)
- [Build Matter Controller](docs/Build%20Matter%20Controller.md)
