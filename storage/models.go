package storage

import (
	"time"
)

type Conversation struct {
	ID        uint `gorm:"primaryKey"`
	UserID    string
	Messages  []byte `gorm:"type:json"` // JSON serialized messages
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Config struct {
	UserID      string `gorm:"primaryKey"`
	LLMProvider string
	Model       string
	MaxTokens   int
	UpdatedAt   time.Time
}

type ToolExecution struct {
	ID        uint `gorm:"primaryKey"`
	UserID    string
	ToolName  string
	Params    []byte `gorm:"type:json"`
	Result    string
	Error     string
	CreatedAt time.Time
}

type UserState struct {
	UserID    string `gorm:"primaryKey"`
	State     string
	Context   []byte `gorm:"type:json"`
	UpdatedAt time.Time
}
