package channels

import (
	"context"

	"github.com/goclaw/goclaw/ai"
	"github.com/goclaw/goclaw/config"
	"github.com/goclaw/goclaw/logs"
)

type Channel interface {
	Start(ctx context.Context, cfg *config.Config) error
	Stop() error
	Name() string
}

type Manager struct {
	channels map[string]Channel
	logger   logs.Logger
}

func NewManager(logger logs.Logger) *Manager {
	return &Manager{
		channels: make(map[string]Channel),
		logger:   logger,
	}
}

func (m *Manager) Register(name string, channel Channel) {
	m.channels[name] = channel
}

func (m *Manager) StartAll(ctx context.Context, cfg *config.Config) error {
	// Register available channels
	m.registerChannels(cfg)

	// Start all registered channels
	for name, channel := range m.channels {
		m.logger.Infof("Starting channel: %s", name)
		if err := channel.Start(ctx, cfg); err != nil {
			m.logger.Errorf("Failed to start channel %s: %v", name, err)
			// Continue starting other channels even if one fails
		}
	}

	return nil
}

// SetAIEngine sets the AI engine for channels that need it
func (m *Manager) SetAIEngine(aiEngine *ai.AIEngine) {
	for _, channel := range m.channels {
		if webchat, ok := channel.(*WebChatChannel); ok {
			webchat.aiEngine = aiEngine
		}
	}
}

func (m *Manager) StopAll() error {
	for name, channel := range m.channels {
		m.logger.Infof("Stopping channel: %s", name)
		if err := channel.Stop(); err != nil {
			m.logger.Errorf("Failed to stop channel %s: %v", name, err)
		}
	}
	return nil
}

func (m *Manager) GetChannel(name string) Channel {
	return m.channels[name]
}

func (m *Manager) registerChannels(cfg *config.Config) {
	// Register Telegram channel if configured
	if cfg.Telegram.BotToken != "" {
		m.Register("telegram", NewTelegramChannel(m.logger))
	}

	// Register Discord channel if configured
	if cfg.Discord.BotToken != "" {
		m.Register("discord", NewDiscordChannel(m.logger))
	}

	// Always register WebChat channel (AI engine will be set later)
	m.Register("webchat", NewWebChatChannel(m.logger, nil))
}
