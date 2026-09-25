本项目为二改，原项目：https://github.com/Aoluis1005/QQ-farm-BOT-GO 

### 一键部署（推荐，Rocky Linux 9.6 / Debian / Ubuntu）
root 执行即可。Rocky 9 走 dnf，脚本会安装 Go 1.25 与 Node 22、2G 机自动加1G swap、编译前后端、装到 `/opt/go-farm-bot`、注册 systemd、放行 3009。

```bash
git clone https://github.com/Aoluis1005/QQ-farm-BOT-GO.git
cd QQ-farm-BOT-GO
bash install.sh
```

部署完成后打开 **`http://<服务器IP>:3009`云厂商安全组也要放行 TCP 3009。

**后续更新**同样是这两步，脚本会自己先停服务再覆盖、装完自动拉起，不用手动 `systemctl stop`：

```bash
cd QQ-farm-BOT-GO
git pull
sudo bash install.sh
```

更新完想确认线上跑的是哪一版，返回里的 `version` 就是 git 短哈希：

```bash
curl -s http://127.0.0.1:3009/api/health
```

> 💡 `install.sh` 每次都会**自动重新构建前端**（检测并自动安装 Node → `npm ci` + `vite build`）后再编译后端，保证页面主题/样式始终完整。请勿删除或替换 `web/dist`，也不要手动放置旧版可执行文件——否则可能导致页面白屏、无任何 UI 样式。

### Docker 部署（可选）
仓库自带 `Dockerfile` + `docker-compose.yml`，适合没有 systemd 的环境（LXC 容器 / 群晖 / 其他主机），与 install.sh 互不影响：

```bash
# 方式一：docker compose（推荐）
docker compose up -d --build

# 方式二：纯 docker（可选指定版本号，默认 dev）
docker build --build-arg VERSION=$(git rev-parse --short HEAD) -t go-farm-bot .
docker run -d --name go-farm-bot -p 3009:3009 -e ADMIN_PORT=3009 \
  -v qq-farm-bot-data:/root/.qq-farm-bot go-farm-bot
```

- 管理页面：`http://<服务器IP>:3009`
- **数据持久化**：账号/配置/日志全部在容器卷 `qq-farm-bot-data`（挂载到 `/root/.qq-farm-bot`），重建容器不丢数据
- **前端内嵌**：`web/dist` 已随镜像打进二进制，无需额外静态服务器
- **资源目录**：`game-config`（素材配置）与 `yyb-resource`（YYB 扫码）已随镜像复制，无需手动放置
- 升级：`docker compose up -d --build` 重新构建即可；想换版本号先 `FARM_VERSION=<hash> docker compose build`
- 时区无关：日志固定输出北京时间（UTC+8）

## 📦 快速开始

### 编译
```bash
# Go 1.20+
CGO_ENABLED=0 go build -o go-farm-bot .
```

### 运行
```bash
./go-farm-bot
# 默认监听 :3009
```

### 首次配置
1. 浏览器打开 `http://<服务器IP>:3009`
2. 进入管理面板 → **设置** → 设置管理密码（`/api/admin/setup`）
3. 添加账号 → 微信扫码 / 授权登录（YYB 渠道 = wx 平台）
4. 打开自动化开关，Bot 即开始挂机

### systemd 部署（推荐，裸机）
```ini
# /etc/systemd/system/go-farm-bot.service
[Unit]
Description=QQ Farm Bot (Go)
After=network.target

[Service]
WorkingDirectory=/opt/go-farm-bot
ExecStart=/opt/go-farm-bot/go-farm-bot
Environment=ADMIN_PORT=3009
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
```bash
sudo systemctl enable --now go-farm-bot
```


## ⚠️ 免责声明
- 仅用于学习与个人自动化研究，请遵守游戏用户协议
- 本项目完全免费

## 📝 相关项目
- 本项目为二改，原项目：https://github.com/Aoluis1005/QQ-farm-BOT-GO  —— 本项目的协议参考与功能对照
