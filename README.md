# My-OpenClaw

A minimal implementation of OpenClaw's Gateway architecture, built with **Go backend** and **React frontend**.

## Features

- WebSocket-based real-time chat
- Session management
- Message history
- Clean UI with React
- Extensible architecture

## Quick Start

### 1. Install Dependencies

```bash
# Install Go dependencies
cd backend
go mod tidy

# Install React dependencies
cd ../frontend
npm install
```

### 2. Start Backend

```bash
cd backend
go run ./cmd/gateway
```

Gateway will start at `ws://127.0.0.1:18789/ws`

### 3. Start Frontend

```bash
cd frontend
npm run dev
```

Frontend will start at `http://localhost:3000`

### 4. Or Run Both

```bash
# Terminal 1 - Backend
cd backend && go run ./cmd/gateway

# Terminal 2 - Frontend
cd frontend && npm run dev
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Gateway info |
| `GET /health` | Health check |
| `GET /status` | Gateway status |
| `WS /ws` | WebSocket connection |

## WebSocket Protocol

### Connect
```json
{
  "type": "req",
  "id": "1",
  "method": "connect",
  "params": {
    "client": { "id": "web", "version": "1.0.0", "platform": "web" }
  }
}
```

### Send Message
```json
{
  "type": "req",
  "id": "2",
  "method": "chat.send",
  "params": {
    "session": "main",
    "message": "Hello!"
  }
}
```

### Response
```json
{
  "type": "res",
  "id": "2",
  "ok": true,
  "payload": {
    "messageId": "msg_123",
    "content": "Hello! How can I help you?",
    "session": "main"
  }
}
```

## Project Structure

```
My-OpenClaw/
├── backend/
│   ├── cmd/gateway/main.go     # Entry point
│   ├── internal/gateway/        # Gateway server
│   └── pkg/types/              # Shared types
├── frontend/
│   ├── src/
│   │   ├── App.tsx            # Main React app
│   │   └── main.tsx           # React entry
│   ├── index.html
│   └── package.json
├── live2d/                     # Live2D 桌面宠物
└── README.md
```

## Next Steps

To extend this minimal version:

1. **Add real LLM integration** - Replace mock agent with actual API calls
2. **Add more methods** - `sessions.list`, `status.get`, etc.
3. **Add authentication** - Token-based auth for remote connections
4. **Add channels** - Telegram, Discord, Slack adapters
5. **Add tools** - Shell execution, file operations, browser automation
6. **Add skills** - SKILL.md-based capability extension

## References

- [OpenClaw Official](https://openclaw.ai)
- [OpenClaw Documentation](https://docs.openclaw.ai)
- [Gateway Protocol](https://docs.openclaw.ai/gateway/protocol)
