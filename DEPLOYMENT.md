# GoClaw 部署指南

## 开发环境

### 启动服务
```bash
# Windows
./run.bat

# 或直接
./goclaw.exe serve
```

### 访问界面
- Web 聊天界面: http://localhost:8080/chat
- 健康检查: http://localhost:8080/health
- 根路径: http://localhost:8080/

## 生产部署

### 1. 使用 systemd (Linux)

创建服务文件 `/etc/systemd/system/goclaw.service`:

```ini
[Unit]
Description=GoClaw AI Assistant
After=network.target

[Service]
Type=simple
User=goclaw
WorkingDirectory=/opt/goclaw
ExecStart=/opt/goclaw/goclaw.exe serve
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启用并启动服务:
```bash
sudo systemctl enable goclaw
sudo systemctl start goclaw
sudo systemctl status goclaw
```

### 2. 使用 Docker

创建 Dockerfile:

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o goclaw ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/goclaw .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./goclaw", "serve"]
```

构建和运行:
```bash
docker build -t goclaw .
docker run -d -p 8080:8080 --name goclaw goclaw
```

### 3. 使用 Supervisor

创建配置文件 `/etc/supervisor/conf.d/goclaw.conf`:

```ini
[program:goclaw]
command=/path/to/goclaw.exe serve
directory=/path/to/goclaw
autostart=true
autorestart=true
stderr_logfile=/var/log/goclaw.err.log
stdout_logfile=/var/log/goclaw.out.log
```

## 配置生产环境

### 环境变量

```bash
export GOCRAW_APP_PORT=8080
export GOCRAW_LOG_LEVEL=info
export GOCRAW_AI_API_KEY=your-api-key
export GOCRAW_TELEGRAM_BOT_TOKEN=your-bot-token
export GOCRAW_DISCORD_BOT_TOKEN=your-bot-token
```

### 安全配置

1. **HTTPS**: 使用反向代理 (Nginx/Apache)
2. **防火墙**: 只开放必要端口
3. **API 密钥**: 使用环境变量，不要硬编码
4. **日志**: 配置日志轮转

### Nginx 反向代理示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## 监控

### 健康检查
```bash
curl http://localhost:8080/health
```

### 指标
```bash
curl http://localhost:8080/metrics
```

## 故障排除

### 端口被占用
```bash
# Windows
netstat -ano | findstr :8080

# Linux
lsof -i :8080
```

### 日志查看
```bash
# Windows (查看事件查看器)
# Linux
tail -f /var/log/goclaw.log
```

### 常见问题

1. **无法启动**: 检查配置文件路径和权限
2. **404 错误**: 确保模板文件存在
3. **WebSocket 连接失败**: 检查防火墙设置

## 备份和恢复

### 备份数据库
```bash
cp goclaw.db goclaw.db.backup
```

### 恢复
```bash
cp goclaw.db.backup goclaw.db
```
