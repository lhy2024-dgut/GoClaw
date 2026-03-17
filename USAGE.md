# GoClaw 使用指南

## 启动服务

### Windows
```bash
./goclaw.exe serve
```

### Linux/macOS
```bash
./goclaw serve
```

## Web 界面

启动后，访问 `http://localhost:8080/chat` 打开聊天界面。

## API 端点

### 1. 根路径
```
GET /
```
返回服务信息和可用端点。

### 2. 健康检查
```
GET /health
```
检查服务状态。

### 3. 指标
```
GET /metrics
```
获取服务指标。

### 4. 聊天 API
```
POST /api/chat
```
发送消息（待实现）。

### 5. 聊天历史
```
GET /api/chat/history
```
获取聊天历史（待实现）。

## 配置

配置文件位于 `~/.goclaw/config.yaml`：

```yaml
app:
  port: 8080  # HTTP 服务器端口

log:
  level: info  # 日志级别：debug, info, warn, error

ai:
  provider: openai  # LLM 提供商
  model: gpt-4o-mini  # 模型名称
  api_key: ""  # API 密钥（必须配置）

telegram:
  bot_token: ""  # Telegram 机器人令牌（可选）

discord:
  bot_token: ""  # Discord 机器人令牌（可选）
```

## 功能状态

### ✅ 已实现
- HTTP 服务器和路由
- 健康检查和指标端点
- Web 聊天界面（基础）
- 多渠道框架（Telegram、Discord、WebChat）
- 配置管理
- 日志系统

### 🚧 开发中
- LLM 集成
- 真实聊天功能
- 工具执行系统
- 对话记忆

## 示例

### 测试服务
```bash
# 检查服务状态
curl http://localhost:8080/health

# 查看服务信息
curl http://localhost:8080/
```

### 使用 Web 界面
1. 启动服务
2. 打开浏览器访问 `http://localhost:8080/chat`
3. 输入消息并发送

## 下一步

1. 配置 AI API 密钥
2. 配置聊天渠道（Telegram/Discord）
3. 测试完整聊天功能
4. 部署到服务器
