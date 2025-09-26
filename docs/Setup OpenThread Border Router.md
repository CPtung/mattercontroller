# Setup OpenThread Border Router

# 環境

- x64_86 Debian 11

# 前置作業

由於沒有相關支援Thread的晶片系統 (SoC) 和網路共同處理器 (NCP)，以下採用官方提供的Simulation方式進行模擬。

### 1. 在 Host 上啟動模擬 RCP

直接從官方的Docker container中將 `ot-rcp` 複製到 DUT上。

```bash
docker cp openthread-environment:/openthread/build/examples/apps/ncp/ot-rcp /usr/local/bin/
```

沒有NCP的條件下，RCP必須有可以對應的 tty 串口，因此這邊採用 socat 的方式讓RCP可以正常啟動。

```bash
$ socat -d -d EXEC:"ot-rcp 1",pty,raw,echo=0 pty,raw,echo=0

socat[1880053] N forking off child, using pty for reading and writing
socat[1880053] N forked off child process 1880055
socat[1880053] N forked off child process 1880055
socat[1880055] N execvp'ing "ot-rcp"
socat[1880053] N PTY is **/dev/pts/1**
socat[1880053] N starting data transfer loop with FDs [5,5] and [7,7]
```

### 2. 在 Host 上啟動模擬 Border Router

首先從官方 Docker Repository上下載 `openthread/border-router:latest` image。

```bash
docker pull openthread/border-router:latest:latest
```

建立 Border Router 用的設定檔，讓它連 RCP的 tty port `/dev/pts/1`：

```bash
OT_RCP_DEVICE=spinel+hdlc+uart:///dev/pts/1?uart-baudrate=115200
OT_INFRA_IF=eth0      # 或 wlan0，看你的環境
OT_THREAD_IF=wpan0
OT_LOG_LEVEL=7
```

啟動容器

```bash
docker run --name otbr -d --rm \
		--cap-add=net_admin \
		--env-file=/var/lib/matter/otbr-env.list \
		--network=host \
		-v /dev/pts:/dev/pts \
		-v /var/lib/matter:/var/lib/matter \
		--device=/dev/net/tun \
		--volume=/var/lib/otbr:/data \
		openthread/border-router:latest
```

# 驗證

進容器

```bash
docker exec -it otbr bash
```

在 CLI 裡面建置網路

```bash
$ ot-ctl
> dataset init new
> dataset commit active
> ifconfig up
> thread start
> state

detached
Done

----------------------------------  狀態需要等幾秒才會切換成 leader

> state

leader
Done

> router table
| ID | RLOC16 | Next Hop | Path Cost | LQ In | LQ Out | Age | Extended MAC     | Link |
+----+--------+----------+-----------+-------+--------+-----+------------------+------+
| 22 | 0x5800 |       63 |         0 |     0 |      0 |   0 | 56d12fd822a56371 |    0 |

Done
```