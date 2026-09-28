# QQ农场自动化助手（Go 版）

本项目为二改，原项目：https://github.com/Aoluis1005/QQ-farm-BOT-GO

## 功能亮点

- **多账号管理**：支持添加多个账号，统一界面管理
- **多用户系统**：支持多用户登录、注册、卡密激活，完善的权限管理（admin/user角色）
- **卡密系统**：支持卡密激活、用户续期、卡密发放与领取，防刷机制
- **自动巡查**：自动偷菜、帮忙、收菜，省时省力
- **活动中心**：查看活动状态，一键领取秋祈良愿、快乐不独享等奖励
- **Web 管理面板**：响应式设计，支持移动端访问
- **一键部署**：install.sh 脚本自动安装依赖、编译前后端
- **Docker 支持**：docker-compose 一键启动，数据持久化

### 多用户与卡密系统

本项目集成了完善的多用户管理系统和卡密激活机制：

- **用户注册/登录**：支持多用户注册账号，密码经 PBKDF2 加密存储
- **角色权限**：admin 角色拥有全部权限，user 角色受限于卡密时长
- **卡密激活**：用户可通过卡密激活账号，获得指定天数的使用权限
- **续期续费**：支持使用卡密为已有用户续期
- **卡密防刷**：每IP每24小时限领1张卡密， UA 绑定防作弊

### API 接口概览

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/users/register` | POST | 用户注册 |
| `/api/users/login` | POST | 用户登录 |
| `/api/users/me` | GET | 获取用户信息 |
| `/api/users/renew` | POST | 用户续期 |
| `/api/users/delete` | POST | 删除用户 |
| `/api/users/claim-card` | POST | 领取卡密（防刷） |
| `/api/cards` | GET | 获取用户卡列表 |

### 一键部署（推荐，Rocky Linux 9.6 / Debian / Ubuntu）
root 执行即可。Rocky 9 走 dnf，脚本会安装 Go 1.25 与 Node 22、2G 机自动加 1G swap、编译前后端、装到 `/opt/go-farm-bot`、注册 systemd、放行 3009。

```bash
git clone https://github.com/wuyoukeji8888-stack/qq-farm-bot-go-wy.git
cd qq-farm-bot-go-wy
bash install.sh
```

部署完成后打开 **`http://<服务器 IP>:3009`**，注意云厂商安全组也要放行 TCP 3009。

**后续更新**同样是这两步，脚本会自己先停服务再覆盖、装完自动拉起，不用手动 `systemctl stop`：

```bash
cd qq-farm-bot-go-wy
git pull
sudo bash install.sh
```

更新完想确认线上跑的是哪一版，登录后在个人资料栏以年月日时分秒的方式标记版本更新时间

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
# 注意：需先构建前端（web/dist），否则 embed 会失败
cd web && npm ci && npm run build && cd ..
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
2. **登录**：使用默认账号 `admin` / `admin` 登录
3. **注册新用户**：在登录页面选择「注册」或访问 `/api/users/register` 注册新账号
4. **获取卡密**：管理员可在 `users.json` 中添加卡密，用户可用卡密注册/续期
5. **开启自动化**：登录后进入设置，添加微信账号后打开自动化开关

> **⚠️ 安全提示**：首次登录后请尽快修改 `admin` 默认密码！
> - 访问 `POST /api/users/login` 登录后记下 token
> - 使用 `POST /api/admin/change-password` 修改密码
> - 或直接修改 `users.json` 中的 password 字段

### 多用户系统

| 角色 | 说明 |
|------|------|
| `admin` | 拥有全部权限，默认账号 `admin/admin` |
| `user` | 受卡密时长限制，无卡密则无法使用 |

登录方式：
```bash
# 1. 登录获取 token
curl -X POST http://localhost:3009/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# 2. 使用 token 访问功能面板
curl http://localhost:3009/api/accounts \
  -H "Authorization: Bearer <token>"
```

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
- 当前仓库：https://github.com/wuyoukeji8888-stack/qq-farm-bot-go-wy
