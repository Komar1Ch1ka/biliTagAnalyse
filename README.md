# B站推荐视频 Tag 统计爬虫

爬取 B站 首页推荐视频，提取视频 Tag 并统计出现频次，支持多种分析模式。

## 功能特性

- 爬取 B站 首页推荐视频 Tag
- 支持多轮爬取和统计
- 三种运行模式：JSON输出、Ollama本地模型分析、远程API分析
- 并发爬取，支持自定义并发数和请求间隔

## 快速开始

### 1. 获取 B站 Cookie

1. 登录 B站网页版
2. 按 F12 打开开发者工具
3. 切换到 Network -> 选择任意 bilibili.com 请求 -> 复制 Cookie 值
4. 将 Cookie 填入 config.json

### 2. 编译程序

```bash
go build -o biliTagAnalyse.exe .
```

### 3. 运行程序

```bash
# 显示帮助信息
./biliTagAnalyse.exe -help

# JSON文件输出模式（默认）
./biliTagAnalyse.exe -json

# Ollama 本地模型分析模式
./biliTagAnalyse.exe -ollama

# 远程 API 分析模式
./biliTagAnalyse.exe -api -api-endpoint https://api.example.com/v1/chat -api-key your-key
```

## 运行模式

程序支持三种互斥的运行模式，优先级从高到低：

| 模式 | 参数 | 说明 |
|------|------|------|
| JSON输出模式 | `-json` | 仅生成JSON格式文件，不进行模型分析或API调用 |
| Ollama模式 | `-ollama` | 调用本地部署的Ollama模型进行数据分析 |
| API模式 | `-api` | 调用远程模型API接口完成分析任务 |

### JSON输出模式

仅爬取数据并保存为JSON文件，不调用任何模型：

```bash
./biliTagAnalyse.exe -json
```

输出文件：`results/tags_stats.json`

### Ollama 模式

使用本地 Ollama 模型分析数据：

```bash
# 分析新爬取的数据
./biliTagAnalyse.exe -ollama

# 分析已有的JSON文件
./biliTagAnalyse.exe -ollama -input results/tags_stats.json

# 指定Ollama服务和模型
./biliTagAnalyse.exe -ollama -ollama-url http://localhost:11434 -ollama-model qwen2.5:7b
```

输出文件：`results/analysis_result.json`

### API 模式

调用远程模型API进行分析：

```bash
./biliTagAnalyse.exe -api -api-endpoint https://api.openai.com/v1/chat/completions -api-key sk-xxx
```

输出文件：`results/analysis_result.json`

## 命令行参数

### 通用参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-config` | 配置文件路径 | config.json |
| `-input` | 输入JSON文件路径（用于分析已有数据） | - |
| `-help` | 显示帮助信息 | - |

### Ollama 模式参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-ollama-url` | Ollama服务地址 | http://localhost:11434 |
| `-ollama-model` | Ollama模型名称 | qwen2.5:7b |

### API 模式参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-api-endpoint` | 远程API端点地址 | - |
| `-api-key` | 远程API密钥 | - |

## 环境变量

程序支持通过环境变量配置参数，命令行参数优先级高于环境变量：

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `BILI_CONFIG_PATH` | 配置文件路径 | config.json |
| `BILI_OLLAMA_URL` | Ollama服务地址 | http://localhost:11434 |
| `BILI_OLLAMA_MODEL` | Ollama模型名称 | qwen2.5:7b |
| `BILI_API_ENDPOINT` | 远程API端点地址 | - |
| `BILI_API_KEY` | 远程API密钥 | - |

使用示例：

```bash
# Linux/macOS
export BILI_OLLAMA_URL="http://192.168.1.100:11434"
export BILI_OLLAMA_MODEL="llama3:8b"
./biliTagAnalyse -ollama

# Windows PowerShell
$env:BILI_OLLAMA_URL="http://192.168.1.100:11434"
$env:BILI_OLLAMA_MODEL="llama3:8b"
.\biliTagAnalyse.exe -ollama
```

## 配置文件

config.json 配置文件说明：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| cookie | B站登录 Cookie，必填 | - |
| crawl_count | 爬取轮数 | 5 |
| crawl_interval | 每轮间隔（秒） | 300 |
| request_interval | 请求间隔（毫秒） | 500 |
| max_concurrent | 最大并发数 | 3 |
| retry_count | 失败重试次数 | 3 |
| retry_delay | 重试延迟（秒） | 2 |
| output_file | 结果输出路径 | results/tags_stats.json |
| ollama_url | Ollama服务地址 | http://localhost:11434 |
| ollama_model | Ollama模型名称 | qwen2.5:7b |
| api_endpoint | 远程API端点 | - |
| api_key | 远程API密钥 | - |

## 输出格式

### JSON输出模式 (tags_stats.json)

```json
{
  "crawl_time": "2024-01-01 12:00:00",
  "total_videos": 100,
  "total_tags": 500,
  "tag_stats": [
    {"tag": "游戏", "count": 50},
    {"tag": "科技", "count": 30}
  ]
}
```

### 分析模式 (analysis_result.json)

```json
{
  "summary": "模型生成的分析摘要...",
  "top_tags_insights": [
    {"tag": "游戏", "count": 50, "description": "由模型分析生成"}
  ],
  "trends": ["趋势分析..."],
  "suggestions": ["创作建议..."],
  "raw_stats": { /* 原始统计数据 */ }
}
```

## 项目结构

```
biliTagAnalyse/
├── cmd/
│   ├── cmd.go           # 命令行参数解析
│   └── defaults.go      # 默认值和常量定义
├── config/
│   └── config.go        # 配置文件加载
├── crawler/
│   └── crawler.go       # 爬虫核心逻辑
├── parser/
│   └── parser.go        # HTML解析
├── statistics/
│   └── statistics.go    # 统计计算
├── analyzer/
│   └── analyzer.go      # 分析模式处理
├── utils/
│   └── http.go          # HTTP工具
├── main.go              # 程序入口
├── config.json          # 配置文件
└── results/             # 输出目录
```

## 注意事项

- 合理设置爬取间隔，避免被封禁
- Cookie 有效期较短，需要定期更新
- Ollama 模式需要本地已安装并运行 Ollama 服务
- API 模式需要确保网络可访问远程 API 端点
- 多模式参数同时使用时，优先级为：json > ollama > api

## 依赖

- Go 1.24+
- golang.org/x/net
