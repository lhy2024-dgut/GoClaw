package channels

import (
	"context"
	"fmt"

	"github.com/goclaw/goclaw/config"
	"github.com/goclaw/goclaw/logs"
)

type DiscordChannel struct {
	logger logs.Logger
}

func NewDiscordChannel(logger logs.Logger) *DiscordChannel {
	return &DiscordChannel{
		logger: logger,
	}
}

func (d *DiscordChannel) Name() string {
	return "discord"
}

func (d *DiscordChannel) Start(ctx context.Context, cfg *config.Config) error {
	if cfg.Discord.BotToken == "" {
		return fmt.Errorf("discord bot token is not configured")
	}

	d.logger.Info("Discord channel started (placeholder)")
	// TODO: Implement actual Discord bot logic
	// This would use disgo or similar library

	// Simulate running
	go func() {
		<-ctx.Done()
		d.logger.Info("Discord channel stopping")
	}()

	return nil
}

func (d *DiscordChannel) Stop() error {
	d.logger.Info("Discord channel stopped")
	return nil
}
