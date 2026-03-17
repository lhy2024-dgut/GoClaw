# GoClaw 快速启动指南

## 前提条件

- Go 1.22 或更高版本
- Git

## 快速开始

### 1. 克隆项目（如果还没有）
```bash
git clone <repository-url>
cd GoClaw
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 构建项目
```bash
go build -o goclaw.exe ./cmd
```

### 4. 创建配置文件
```bash
# Windows
mkdir %USERPROFILE%\.goclaw
copy config.example.yaml %USERPROFILE%\.goclaw\config.yaml

# Linux/macOS
mkdir -p ~/.goclaw
cp config.example.yaml ~/.goclaw/config.yaml
```

### 5. 编辑配置文件
打开 `~/.goclaw/config.yaml` 并填入你的 API 密钥：

```yaml
ai:
  api_key: "your-openai-api-key"  # 替换为你的 OpenAI API 密钥

telegram:
  bot_token: "your-telegram-bot-token"  # 可选，从 @BotFather 获取

discord:
  bot_token: "your-discord-bot-token"  # 可选，从 Discord 开发者门户获取
```

### 6. 启动服务器
```bash
# Windows
./goclaw.exe serve

# Linux/macOS
./goclaw serve
```

### 7. 测试服务
访问健康检查端点：
```bash
curl http://localhost:8080/health
```

应该返回：
```json
{
  "status": "healthy",
  "time": "2026-03-16T10:00:00Z"
}
```

## 常用命令

```bash
# 查看版本
./goclaw.exe version

# 查看帮助
./goclaw.exe --help
./goclaw.exe serve --help

# 运行测试
go test ./tests/...
```

## 配置说明

### 基本配置
- `app.port`: HTTP 服务器端口（默认 8080）
- `log.level`: 日志级别（debug、info、warn、error）

### AI 配置
- `ai.provider`: LLM 提供商（openai、anthropic 等）
- `ai.model`: 模型名称
- `ai.api_key`: API 密钥

### 聊天渠道配置
- `telegram.bot_token`: Telegram 机器人令牌
- `discord.bot_token`: Discord 机器人令牌

### 工具配置
- `tools.email`: 邮件查询配置
- `tools.calendar`: 日历管理配置

## 故障排除

### 问题：端口被占用
解决方案：修改 `config.yaml` 中的 `app.port`

### 问题：API 密钥错误
解决方案：检查 `ai.api_key` 是否正确配置

### 问题：数据库连接失败
解决方案：检查 `database.path` 是否有写入权限

## 下一步

1. 配置聊天渠道（Telegram/Discord）
2. 配置工具（邮件、日历）
3. 测试 AI 对话
4. 部署到服务器（使用 systemd 或 Docker）

## 获取帮助

- 查看 `README.md` 了解项目详情
- 查看 `PROJECT_SUMMARY.md` 了解开发计划
- 运行 `./goclaw.exe --help` 查看命令帮助
