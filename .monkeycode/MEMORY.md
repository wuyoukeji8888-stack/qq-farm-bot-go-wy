# User Instruction Memory

This file records user instructions, preferences, and teachings for reference in future interactions.

## Format

### User Instruction Entry
[User Instruction Summary]
- Date: 2026-09-25
- Context: 推送到仓库前更新助手版本
- Instructions:
  - 每次执行 git push 前，先重新生成 version.go（运行 bash install.sh 或手动执行：VERSION=$(git rev-parse --short HEAD) && printf 'package main\n\nvar buildVersion = "%s"\nvar buildTime = "%s"\n' "$VERSION" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > version.go）
  - 然后构建前端（cd web && npm run build）
  - 再编译后端（go build -o go-farm-bot .）
  - 最后才执行 git add -A && git commit && git push

### Project Knowledge Entry
[Project Knowledge Summary]
- Date: 2026-09-25
- Context: Agent 发现的项目构建规范
- Category: Build Methods
- Instructions:
  - Go 模块名仍是 github.com/Aoluis1005/go-farm-bot（与 GitHub 仓库名 qq-farm-bot-go-wy 不同）
  - 可执行文件名为 go-farm-bot，安装到 /opt/go-farm-bot
  - 前端构建产物嵌入后端二进制（//go:embed all:web/dist），必须先构建前端再编译后端
  - 助手版本号由 /api/health 接口返回，格式为 buildTime（UTC ISO 8601），前端转换为北京时间显示
