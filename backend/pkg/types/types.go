package types

import "time"

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type Session struct {
	Key        string    `json:"key"`
	AgentID    string    `json:"agentId"`
	Transcript []Message `json:"transcript"`
	Memory     string    `json:"memory"`
	CreatedAt  time.Time `json:"createdAt"`
	LastActive time.Time `json:"lastActive"`
}

type RequestFrame struct {
	Type   string         `json:"type"`
	ID     string         `json:"id"`
	Method string         `json:"method"`
	Params map[string]any `json:"params,omitempty"`
}

type ResponseFrame struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	OK      bool   `json:"ok"`
	Payload any    `json:"payload,omitempty"`
	Error   string `json:"error,omitempty"`
}

type EventFrame struct {
	Type    string `json:"type"`
	Event   string `json:"event"`
	Payload any    `json:"payload,omitempty"`
	Seq     int64  `json:"seq,omitempty"`
}

type AuthContext struct {
	OK       bool
	Role     string
	Scopes   []string
	DeviceID string
}

type Config struct {
	Gateway GatewayConfig `yaml:"gateway"`
	Secrets Secrets       `yaml:"secrets"`
	Agent   AgentConfig   `yaml:"agent"`
}

type GatewayConfig struct {
	Port int    `yaml:"port"`
	Bind string `yaml:"bind"`
	Auth Auth   `yaml:"auth"`
}

type Auth struct {
	Mode  string `yaml:"mode"`
	Token string `yaml:"token"`
}

type Secrets map[string]string

type AgentConfig struct {
	Model string `yaml:"model"`
}
