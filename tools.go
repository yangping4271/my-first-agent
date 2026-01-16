package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// ToolDefinition 定义一个工具的结构
type ToolDefinition struct {
	Name        string
	Description string
	InputSchema anthropic.ToolInputSchemaParam
	Function    func(input json.RawMessage) (string, error)
}

// 将 ToolDefinition 转换为 Anthropic API 格式
func (t ToolDefinition) ToAnthropicTool() anthropic.ToolUnionParam {
	return anthropic.ToolUnionParam{
		OfTool: &anthropic.ToolParam{
			Name:        t.Name,
			Description: param.NewOpt(t.Description),
			InputSchema: t.InputSchema,
		},
	}
}

// read_file 工具：读取文件内容
func readFileTool() ToolDefinition {
	return ToolDefinition{
		Name:        "read_file",
		Description: "Read the contents of a file at the specified path",
		InputSchema: anthropic.ToolInputSchemaParam{
			Type: "object",
			Properties: map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The path to the file to read",
				},
			},
			Required: []string{"path"},
		},
		Function: func(input json.RawMessage) (string, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("failed to parse input: %w", err)
			}

			content, err := os.ReadFile(params.Path)
			if err != nil {
				return "", fmt.Errorf("failed to read file: %w", err)
			}

			return string(content), nil
		},
	}
}

// list_files 工具：列出目录中的文件
func listFilesTool() ToolDefinition {
	return ToolDefinition{
		Name:        "list_files",
		Description: "List all files and directories at the specified path",
		InputSchema: anthropic.ToolInputSchemaParam{
			Type: "object",
			Properties: map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The directory path to list files from (default: current directory)",
				},
			},
		},
		Function: func(input json.RawMessage) (string, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("failed to parse input: %w", err)
			}

			if params.Path == "" {
				params.Path = "."
			}

			var files []string
			err := filepath.Walk(params.Path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				files = append(files, path)
				return nil
			})

			if err != nil {
				return "", fmt.Errorf("failed to list files: %w", err)
			}

			result, err := json.Marshal(files)
			if err != nil {
				return "", fmt.Errorf("failed to marshal result: %w", err)
			}

			return string(result), nil
		},
	}
}

// edit_file 工具：编辑文件内容
func editFileTool() ToolDefinition {
	return ToolDefinition{
		Name:        "edit_file",
		Description: "Edit a file by replacing old_text with new_text. Creates the file if it doesn't exist.",
		InputSchema: anthropic.ToolInputSchemaParam{
			Type: "object",
			Properties: map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The path to the file to edit",
				},
				"old_text": map[string]interface{}{
					"type":        "string",
					"description": "The text to replace",
				},
				"new_text": map[string]interface{}{
					"type":        "string",
					"description": "The replacement text",
				},
			},
			Required: []string{"path", "old_text", "new_text"},
		},
		Function: func(input json.RawMessage) (string, error) {
			var params struct {
				Path    string `json:"path"`
				OldText string `json:"old_text"`
				NewText string `json:"new_text"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("failed to parse input: %w", err)
			}

			var content string
			if data, err := os.ReadFile(params.Path); err == nil {
				content = string(data)
			}

			newContent := strings.ReplaceAll(content, params.OldText, params.NewText)

			if err := os.WriteFile(params.Path, []byte(newContent), 0644); err != nil {
				return "", fmt.Errorf("failed to write file: %w", err)
			}

			return fmt.Sprintf("Successfully edited %s", params.Path), nil
		},
	}
}

// GetAllTools 返回所有可用的工具
func GetAllTools() []ToolDefinition {
	return []ToolDefinition{
		readFileTool(),
		listFilesTool(),
		editFileTool(),
	}
}
