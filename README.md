# 代码编辑 Agent

基于 [ampcode 文章](https://ampcode.com/how-to-build-an-agent) 实现的代码编辑 Agent，使用 Go 和 Anthropic Claude API。

## 功能

该 Agent 通过对话式界面帮助你编辑代码，支持以下工具：

- **read_file**: 读取文件内容
- **list_files**: 列出目录中的文件和文件夹
- **edit_file**: 通过字符串替换编辑文件

## 配置

通过环境变量配置：

### 必需
- `ANTHROPIC_API_KEY`: Anthropic API 密钥

### 可选
- `ANTHROPIC_BASE_URL`: 自定义 API endpoint（支持代理或第三方 Claude 兼容服务）
- `ANTHROPIC_MODEL`: 自定义模型（默认：`claude-3-5-sonnet-20241022`）

## 安装

```bash
# 克隆或下载代码
cd /path/to/project

# 安装依赖
go mod download

# 构建
go build -o agent
```

## 使用

### 基本使用

```bash
export ANTHROPIC_API_KEY="your-api-key"
./agent
```

### 使用自定义 BaseURL

```bash
export ANTHROPIC_API_KEY="your-api-key"
export ANTHROPIC_BASE_URL="https://your-proxy.com/v1"
./agent
```

### 使用自定义模型

```bash
export ANTHROPIC_API_KEY="your-api-key"
export ANTHROPIC_MODEL="claude-3-opus-20240229"
./agent
```

## 示例对话

```
代码编辑 Agent 已启动！
输入 'exit' 或 'quit' 退出

你: 列出当前目录的文件
Agent: [列出文件列表]

你: 读取 main.go 的内容
Agent: [显示 main.go 内容]

你: 将 main.go 中的 "Hello" 替换为 "Hi"
Agent: [执行替换并确认]

你: exit
再见！
```

## 架构

项目由三个主要文件组成：

- **main.go**: 入口点，处理环境变量和交互循环
- **agent.go**: Agent 核心逻辑，包括推理循环和工具执行
- **tools.go**: 工具定义和实现

## 技术细节

- 使用 `github.com/anthropics/anthropic-sdk-go` SDK
- 实现了完整的工具调用循环
- 支持多轮对话和上下文保持
- 无状态服务器设计，所有历史记录在客户端维护

## 许可

根据原文章实现的教育项目。
