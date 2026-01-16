package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 从环境变量读取配置
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Println("错误: 请设置 ANTHROPIC_API_KEY 环境变量")
		os.Exit(1)
	}

	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if baseURL != "" {
		fmt.Printf("使用自定义 BaseURL: %s\n", baseURL)
	}

	model := os.Getenv("ANTHROPIC_MODEL")
	if model != "" {
		fmt.Printf("使用自定义模型: %s\n", model)
	} else {
		fmt.Println("使用默认模型: claude-3-5-sonnet-20241022")
	}

	// 创建 Agent
	agent := NewAgent(apiKey, baseURL, model)

	fmt.Println("代码编辑 Agent 已启动！")
	fmt.Println("输入 'exit' 或 'quit' 退出")
	fmt.Println()

	// 交互循环
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("你: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("再见！")
			break
		}

		// 发送消息并获取响应
		response, err := agent.SendMessage(input)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}

		fmt.Printf("\nAgent: %s\n\n", response)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取输入错误: %v\n", err)
	}
}
