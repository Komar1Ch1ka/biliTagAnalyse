package cmd

const (
	DefaultConfigPath  = "config.json"
	DefaultOllamaURL   = "http://localhost:11434"
	DefaultOllamaModel = "qwen2.5:7b"
)

const (
	EnvOllamaURL   = "BILI_OLLAMA_URL"
	EnvOllamaModel = "BILI_OLLAMA_MODEL"
	EnvAPIEndpoint = "BILI_API_ENDPOINT"
	EnvAPIKey      = "BILI_API_KEY"
	EnvConfigPath  = "BILI_CONFIG_PATH"
)

const (
	HelpHeader = "=== B站推荐视频 Tag 分析爬虫 ==="
	
	HelpUsage = `用法: biliTagAnalyse [选项]`

	HelpModeSection = `运行模式（互斥，优先级从高到低）：
  -json           JSON文件输出模式：仅生成JSON格式文件，不进行模型分析或API调用
  -ollama         Ollama模式：调用本地部署的Ollama模型进行数据分析
  -api            API模式：调用远程模型API接口完成分析任务`

	HelpCommonSection = `通用选项：
  -config string      配置文件路径 (默认: %s)
  -input string       输入JSON文件路径（用于ollama/api模式分析已有数据）
  -help               显示帮助信息`

	HelpOllamaSection = `Ollama模式选项：
  -ollama-url string      Ollama服务地址 (默认: %s)
  -ollama-model string    Ollama模型名称 (默认: %s)`

	HelpAPISection = `API模式选项：
  -api-endpoint string    远程API端点地址
  -api-key string         远程API密钥`

	HelpExamples = `示例：
  biliTagAnalyse -json                      # 仅生成JSON文件
  biliTagAnalyse -ollama                    # 使用Ollama分析新爬取的数据
  biliTagAnalyse -ollama -input data.json   # 使用Ollama分析已有JSON文件
  biliTagAnalyse -api -api-endpoint https://api.example.com/v1/chat`
)

const (
	ErrMultipleModes    = "错误：只能指定一种运行模式"
	ErrModePriority     = "模式优先级：json > ollama > api"
	ErrOllamaURL        = "Ollama模式需要指定 -ollama-url"
	ErrOllamaModel      = "Ollama模式需要指定 -ollama-model"
	ErrAPIEndpoint      = "API模式需要指定 -api-endpoint"
)

const (
	ModeDescJSON   = "JSON文件输出模式"
	ModeDescOllama = "Ollama本地模型分析模式 (模型: %s, 地址: %s)"
	ModeDescAPI    = "远程API调用模式 (端点: %s)"
	ModeDescUnknown = "未知模式"
)
