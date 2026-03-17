package tools

import (
	"context"
	"fmt"

	"github.com/goclaw/goclaw/logs"
)

type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]interface{} // JSON schema for parameters
	Execute(ctx context.Context, params map[string]interface{}) (string, error)
}

type Registry struct {
	tools  map[string]Tool
	logger logs.Logger
}

func NewRegistry(logger logs.Logger) *Registry {
	return &Registry{
		tools:  make(map[string]Tool),
		logger: logger,
	}
}

func (r *Registry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
	r.logger.Infof("Registered tool: %s", tool.Name())
}

func (r *Registry) Execute(name string, params map[string]interface{}) (string, error) {
	tool, exists := r.tools[name]
	if !exists {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	r.logger.Infof("Executing tool: %s with params: %v", name, params)
	return tool.Execute(context.Background(), params)
}

func (r *Registry) ListTools() []string {
	var toolNames []string
	for name := range r.tools {
		toolNames = append(toolNames, name)
	}
	return toolNames
}

func (r *Registry) GetTool(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}
