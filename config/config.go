package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"app"`

	Log struct {
		Level  string `mapstructure:"level"`
		File   string `mapstructure:"file"`
		Format string `mapstructure:"format"`
	} `mapstructure:"log"`

	Database struct {
		Type string `mapstructure:"type"`
		Path string `mapstructure:"path"`
	} `mapstructure:"database"`

	AI struct {
		Provider    string  `mapstructure:"provider"`
		Model       string  `mapstructure:"model"`
		APIKey      string  `mapstructure:"api_key"`
		APIBase     string  `mapstructure:"api_base"`
		MaxTokens   int     `mapstructure:"max_tokens"`
		Temperature float64 `mapstructure:"temperature"`
	} `mapstructure:"ai"`

	Telegram struct {
		BotToken string `mapstructure:"bot_token"`
		Webhook  string `mapstructure:"webhook"`
	} `mapstructure:"telegram"`

	Discord struct {
		BotToken string `mapstructure:"bot_token"`
		GuildID  string `mapstructure:"guild_id"`
	} `mapstructure:"discord"`

	Tools struct {
		Email struct {
			IMAPServer string `mapstructure:"imap_server"`
			Username   string `mapstructure:"username"`
			Password   string `mapstructure:"password"`
		} `mapstructure:"email"`

		Calendar struct {
			Provider string `mapstructure:"provider"`
			APIKey   string `mapstructure:"api_key"`
		} `mapstructure:"calendar"`
	} `mapstructure:"tools"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add config paths in order of priority
	viper.AddConfigPath(".")
	viper.AddConfigPath("./$HOME/.goclaw")
	viper.AddConfigPath("$HOME/.goclaw")
	viper.AddConfigPath("/etc/goclaw")

	// Add support for Windows user profile path
	homeDir, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(homeDir, ".goclaw"))
	}

	// Also try the current working directory's $HOME/.goclaw
	cwd, err := os.Getwd()
	if err == nil {
		viper.AddConfigPath(filepath.Join(cwd, "$HOME", ".goclaw"))
	}

	// Set defaults
	viper.SetDefault("app.name", "GoClaw")
	viper.SetDefault("app.version", "0.1.0")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.path", "goclaw.db")
	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "gpt-4o-mini")
	viper.SetDefault("ai.max_tokens", 1024)
	viper.SetDefault("ai.temperature", 0.7)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			if err := createDefaultConfig(); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func createDefaultConfig() error {
	configDir := filepath.Join("$HOME", ".goclaw")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	defaultConfig := `
app:
  name: "GoClaw"
  version: "0.1.0"
  port: 8080

log:
  level: "info"
  file: ""
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
`

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}
