package cmd

import (
	"flag"
	"fmt"
	"os"
)

type RunMode int

const (
	ModeJSONOnly RunMode = iota
	ModeOllama
	ModeAPI
)

func (m RunMode) String() string {
	switch m {
	case ModeJSONOnly:
		return "json"
	case ModeOllama:
		return "ollama"
	case ModeAPI:
		return "api"
	default:
		return "unknown"
	}
}

type Options struct {
	ConfigPath   string
	RunMode      RunMode
	OllamaURL    string
	OllamaModel  string
	APIEndpoint  string
	APIKey       string
	ShowHelp     bool
	InputFile    string
}

var (
	flagConfig     = flag.String("config", "config.json", "配置文件路径")
	flagJSON       = flag.Bool("json", false, "JSON文件输出模式：仅生成JSON文件，不进行模型分析")
	flagOllama     = flag.Bool("ollama", false, "Ollama模式：使用本地Ollama模型进行分析")
	flagAPI        = flag.Bool("api", false, "API模式：调用远程模型API进行分析")
	flagOllamaURL  = flag.String("ollama-url", "http://localhost:11434", "Ollama服务地址")
	flagOllamaModel = flag.String("ollama-model", "qwen2.5:7b", "Ollama模型名称")
	flagAPIEndpoint = flag.String("api-endpoint", "", "远程API端点地址")
	flagAPIKey     = flag.String("api-key", "", "远程API密钥")
	flagInput      = flag.String("input", "", "输入JSON文件路径（用于分析模式）")
	flagHelp       = flag.Bool("help", false, "显示帮助信息")
)

func Parse() *Options {
	flag.Usage = printHelp
	flag.Parse()

	opts := &Options{
		ConfigPath:   *flagConfig,
		OllamaURL:    *flagOllamaURL,
		OllamaModel:  *flagOllamaModel,
		APIEndpoint:  *flagAPIEndpoint,
		APIKey:       *flagAPIKey,
		ShowHelp:     *flagHelp,
		InputFile:    *flagInput,
	}

	modeCount := 0
	if *flagJSON {
		opts.RunMode = ModeJSONOnly
		modeCount++
	}
	if *flagOllama {
		opts.RunMode = ModeOllama
		modeCount++
	}
	if *flagAPI {
		opts.RunMode = ModeAPI
		modeCount++
	}

	if modeCount > 1 {
		fmt.Println("错误：只能指定一种运行模式")
		fmt.Println("模式优先级：json > ollama > api")
		opts.RunMode = ModeJSONOnly
	}

	if opts.ShowHelp {
		printHelp()
		os.Exit(0)
	}

	return opts
}

func printHelp() {
	fmt.Println("=== B站推荐视频 Tag 分析爬虫 ===")
	fmt.Println()
	fmt.Println("用法: biliTagAnalyse [选项]")
	fmt.Println()
	fmt.Println("运行模式（互斥，优先级从高到低）：")
	fmt.Println("  -json           JSON文件输出模式：仅生成JSON格式文件，不进行模型分析或API调用")
	fmt.Println("  -ollama         Ollama模式：调用本地部署的Ollama模型进行数据分析")
	fmt.Println("  -api            API模式：调用远程模型API接口完成分析任务")
	fmt.Println()
	fmt.Println("通用选项：")
	fmt.Println("  -config string      配置文件路径 (默认: config.json)")
	fmt.Println("  -input string       输入JSON文件路径（用于ollama/api模式分析已有数据）")
	fmt.Println("  -help               显示帮助信息")
	fmt.Println()
	fmt.Println("Ollama模式选项：")
	fmt.Println("  -ollama-url string      Ollama服务地址 (默认: http://localhost:11434)")
	fmt.Println("  -ollama-model string    Ollama模型名称 (默认: qwen2.5:7b)")
	fmt.Println()
	fmt.Println("API模式选项：")
	fmt.Println("  -api-endpoint string    远程API端点地址")
	fmt.Println("  -api-key string         远程API密钥")
	fmt.Println()
	fmt.Println("示例：")
	fmt.Println("  biliTagAnalyse -json                      # 仅生成JSON文件")
	fmt.Println("  biliTagAnalyse -ollama                    # 使用Ollama分析新爬取的数据")
	fmt.Println("  biliTagAnalyse -ollama -input data.json   # 使用Ollama分析已有JSON文件")
	fmt.Println("  biliTagAnalyse -api -api-endpoint https://api.example.com/v1/chat")
	fmt.Println()
}

func (o *Options) Validate() error {
	switch o.RunMode {
	case ModeOllama:
		if o.OllamaURL == "" {
			return fmt.Errorf("Ollama模式需要指定 -ollama-url")
		}
		if o.OllamaModel == "" {
			return fmt.Errorf("Ollama模式需要指定 -ollama-model")
		}
	case ModeAPI:
		if o.APIEndpoint == "" {
			return fmt.Errorf("API模式需要指定 -api-endpoint")
		}
	}
	return nil
}

func (o *Options) ModeDescription() string {
	switch o.RunMode {
	case ModeJSONOnly:
		return "JSON文件输出模式"
	case ModeOllama:
		return fmt.Sprintf("Ollama本地模型分析模式 (模型: %s, 地址: %s)", o.OllamaModel, o.OllamaURL)
	case ModeAPI:
		return fmt.Sprintf("远程API调用模式 (端点: %s)", o.APIEndpoint)
	default:
		return "未知模式"
	}
}
