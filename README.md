# GoClaw - Personal AI Assistant

GoClaw is a personal AI assistant built with Go, inspired by OpenClaw. It runs 24/7, connects to messaging apps, and can execute tasks.

## Features

- **24/7 Operation**: Runs as a daemon with graceful shutdown
- **Multi-Channel Support**: Connect to Telegram, Discord, and WebChat
- **AI Integration**: LLM-powered conversations with memory
- **Tool Execution**: Email querying and calendar management
- **Persistent Storage**: SQLite database for conversations and state

## Project Structure

```
.
├── cmd/              # Main application entry point
├── channels/         # Messaging channel implementations
│   ├── telegram.go   # Telegram bot integration
│   ├── discord.go    # Discord bot integration
│   └── webchat.go    # WebChat WebSocket server
├── ai/               # AI engine and LLM integration
├── tools/            # Tool execution system
│   ├── email_tool.go # Email querying tool
│   └── calendar_tool.go # Calendar management tool
├── storage/          # Database and persistence
├── config/           # Configuration management
├── logs/             # Logging system
└── tests/            # Test files
```

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the application:
   ```bash
   go build -o goclaw.exe ./cmd
   ```

## Configuration

Create a configuration file at `~/.goclaw/config.yaml`:

```yaml
app:
  name: "GoClaw"
  version: "0.1.0"
  port: 8080

log:
  level: "info"
  format: "text"

database:
  type: "sqlite"
  path: "goclaw.db"

ai:
  provider: "openai"
  model: "gpt-4o-mini"
  api_key: ""
  max_tokens: 1024
  temperature: 0.7

telegram:
  bot_token: ""
  webhook: ""

discord:
  bot_token: ""
  guild_id: ""

tools:
  email:
    imap_server: ""
    username: ""
    password: ""
  calendar:
    provider: "google"
    api_key: ""
```

## Usage

### Start the server

```bash
./goclaw.exe serve
```

### Show version

```bash
./goclaw.exe version
```

## Development

### Running tests

```bash
go test ./...
```

### Adding a new channel

1. Create a new file in `channels/`
2. Implement the `Channel` interface
3. Register it in `channels/manager.go`

### Adding a new tool

1. Create a new file in `tools/`
2. Implement the `Tool` interface
3. Register it in the tool registry

## Architecture

### 24/7 Operation

The application runs as a daemon with:
- Graceful shutdown on SIGINT/SIGTERM
- Health check endpoint at `/health`
- Metrics endpoint at `/metrics`

### Multi-Channel Architecture

Each channel runs independently and communicates with the central AI engine:
```
Telegram → Channel Manager → AI Engine → Response
Discord  ↗                  ↗
WebChat  ↗                  ↗
```

### Tool Execution System

Tools are registered dynamically and can be called by the AI:
```go
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params map[string]interface{}) (string, error)
}
```

## License

MIT License
