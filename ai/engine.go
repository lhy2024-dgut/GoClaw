package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goclaw/goclaw/config"
	"github.com/goclaw/goclaw/logs"
	"github.com/goclaw/goclaw/tools"
	"github.com/sashabaranov/go-openai"
)

type AIEngine struct {
	logger       logs.Logger
	provider     string
	model        string
	client       *openai.Client
	toolRegistry *tools.Registry
}

func NewAIEngine(logger logs.Logger, cfg *config.Config, toolRegistry *tools.Registry) *AIEngine {
	var client *openai.Client
	if cfg.AI.APIKey != "" {
		config := openai.DefaultConfig(cfg.AI.APIKey)

		// Set custom API base for different providers
		if cfg.AI.APIBase != "" {
			config.BaseURL = cfg.AI.APIBase
		} else if cfg.AI.Provider == "qwen" || cfg.AI.Provider == "alibaba" {
			// Default Alibaba Cloud Bailian endpoint
			config.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
		}

		client = openai.NewClientWithConfig(config)
	}

	return &AIEngine{
		logger:       logger,
		provider:     cfg.AI.Provider,
		model:        cfg.AI.Model,
		client:       client,
		toolRegistry: toolRegistry,
	}
}

func (e *AIEngine) Chat(ctx context.Context, prompt string) (string, error) {
	e.logger.Infof("AI Chat: %s", prompt)

	// If no client is configured, return placeholder
	if e.client == nil {
		return fmt.Sprintf("AI response to: %s (no API key configured)"), nil
	}

	// Prepare tools for the request
	var openaiTools []openai.Tool
	if e.toolRegistry != nil {
		for _, toolName := range e.toolRegistry.ListTools() {
			if tool, exists := e.toolRegistry.GetTool(toolName); exists {
				// Convert tool.Parameters() to JSON bytes
				paramsBytes, err := json.Marshal(tool.Parameters())
				if err != nil {
					e.logger.Errorf("Failed to marshal tool parameters: %v", err)
					continue
				}

				// Create OpenAI Tool definition
				funcDef := &openai.FunctionDefinition{
					Name:        tool.Name(),
					Description: tool.Description(),
					Parameters:  paramsBytes,
				}
				openaiTool := openai.Tool{
					Type:     "function",
					Function: funcDef,
				}
				openaiTools = append(openaiTools, openaiTool)
			}
		}
	}

	// Prepare messages
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	// Function calling loop
	for {
		req := openai.ChatCompletionRequest{
			Model:    e.model,
			Messages: messages,
		}

		// Add tools if available
		if len(openaiTools) > 0 {
			req.Tools = openaiTools
		}

		resp, err := e.client.CreateChatCompletion(ctx, req)
		if err != nil {
			e.logger.Errorf("OpenAI API error: %v", err)
			return "", fmt.Errorf("failed to get AI response: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("no response from AI")
		}

		choice := resp.Choices[0]

		// Check if the AI wants to call a function
		if choice.FinishReason == "tool_calls" || choice.Message.ToolCalls != nil {
			// Handle tool calls
			// Add the assistant's message with tool calls to history
			messages = append(messages, choice.Message)

			for _, toolCall := range choice.Message.ToolCalls {
				e.logger.Infof("AI requested tool call: %s", toolCall.Function.Name)

				// Parse arguments
				var params map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
					e.logger.Errorf("Failed to parse tool arguments: %v", err)
					// Handle error response
					messages = append(messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    "Error parsing arguments: " + err.Error(),
						ToolCallID: toolCall.ID,
					})
					continue
				}

				// Execute tool
				result, err := e.toolRegistry.Execute(toolCall.Function.Name, params)
				if err != nil {
					e.logger.Errorf("Tool execution failed: %v", err)
					result = "Error: " + err.Error()
				}

				// Add tool result to messages
				messages = append(messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    result,
					ToolCallID: toolCall.ID,
				})
			}
			// Continue the loop to get the final response
			continue
		}

		// Return the final text response
		return choice.Message.Content, nil
	}
}

func (e *AIEngine) ChatStream(ctx context.Context, prompt string, streamChan chan string) error {
	e.logger.Infof("AI Chat Stream: %s", prompt)
	// TODO: Implement streaming response

	// Placeholder streaming
	streamChan <- fmt.Sprintf("Streaming response to: %s", prompt)
	close(streamChan)
	return nil
}
