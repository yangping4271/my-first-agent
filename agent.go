package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Agent 代表一个代码编辑代理
type Agent struct {
	client   anthropic.Client
	tools    []ToolDefinition
	messages []anthropic.MessageParam
	model    string
}

// NewAgent 创建一个新的 Agent 实例
func NewAgent(apiKey string, baseURL string, model string) *Agent {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}

	// 如果提供了自定义 baseURL，则添加该选项
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	// 如果没有提供模型，使用默认模型
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	client := anthropic.NewClient(opts...)

	return &Agent{
		client:   client,
		tools:    GetAllTools(),
		messages: []anthropic.MessageParam{},
		model:    model,
	}
}

// SendMessage 向 Agent 发送用户消息并获取响应
func (a *Agent) SendMessage(userMessage string) (string, error) {
	// 添加用户消息到历史
	a.messages = append(a.messages, anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)))

	// 转换工具定义为 API 格式
	tools := make([]anthropic.ToolUnionParam, len(a.tools))
	for i, tool := range a.tools {
		tools[i] = tool.ToAnthropicTool()
	}

	// 推理循环
	for {
		// 调用 Claude API
		response, err := a.client.Messages.New(context.Background(), anthropic.MessageNewParams{
			Model:     anthropic.Model(a.model),
			MaxTokens: 4096,
			Messages:  a.messages,
			Tools:     tools,
		})

		if err != nil {
			return "", fmt.Errorf("API call failed: %w", err)
		}

		// 添加 assistant 响应到历史
		assistantBlocks := make([]anthropic.ContentBlockParamUnion, len(response.Content))
		for i, block := range response.Content {
			assistantBlocks[i] = block.ToParam()
		}
		a.messages = append(a.messages, anthropic.NewAssistantMessage(assistantBlocks...))

		// 检查是否需要执行工具
		needsToolExecution := false
		var toolResults []anthropic.ContentBlockParamUnion

		for _, block := range response.Content {
			if block.Type == "tool_use" {
				needsToolExecution = true

				// 执行工具
				result, err := a.executeTool(block.ID, block.Name, block.Input)
				if err != nil {
					result = fmt.Sprintf("Error: %v", err)
				}

				// 添加工具结果
				toolResults = append(toolResults, anthropic.NewToolResultBlock(
					block.ID,
					result,
					false,
				))
			}
		}

		// 如果没有工具需要执行，返回文本响应
		if !needsToolExecution {
			var textResponse string
			for _, block := range response.Content {
				if block.Type == "text" {
					textResponse += block.Text
				}
			}
			return textResponse, nil
		}

		// 添加工具结果到消息历史，继续循环
		a.messages = append(a.messages, anthropic.NewUserMessage(toolResults...))
	}
}

// executeTool 执行指定的工具
func (a *Agent) executeTool(toolID string, toolName string, input interface{}) (string, error) {
	// 查找工具
	var tool *ToolDefinition
	for _, t := range a.tools {
		if t.Name == toolName {
			tool = &t
			break
		}
	}

	if tool == nil {
		return "", fmt.Errorf("tool not found: %s", toolName)
	}

	// 将 input 转换为 JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}

	// 执行工具函数
	result, err := tool.Function(inputJSON)
	if err != nil {
		return "", err
	}

	return result, nil
}
