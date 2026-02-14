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
	ConfigPath  string
	RunMode     RunMode
	OllamaURL   string
	OllamaModel string
	APIEndpoint string
	APIKey      string
	ShowHelp    bool
	InputFile   string
}

var (
	flagConfig      = flag.String("config", "", "配置文件路径")
	flagJSON        = flag.Bool("json", false, "JSON文件输出模式：仅生成JSON文件，不进行模型分析")
	flagOllama      = flag.Bool("ollama", false, "Ollama模式：使用本地Ollama模型进行分析")
	flagAPI         = flag.Bool("api", false, "API模式：调用远程模型API进行分析")
	flagOllamaURL   = flag.String("ollama-url", "", "Ollama服务地址")
	flagOllamaModel = flag.String("ollama-model", "", "Ollama模型名称")
	flagAPIEndpoint = flag.String("api-endpoint", "", "远程API端点地址")
	flagAPIKey      = flag.String("api-key", "", "远程API密钥")
	flagInput       = flag.String("input", "", "输入JSON文件路径（用于分析模式）")
	flagHelp        = flag.Bool("help", false, "显示帮助信息")
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Parse() *Options {
	flag.Usage = printHelp
	flag.Parse()

	configPath := *flagConfig
	if configPath == "" {
		configPath = getEnvOrDefault(EnvConfigPath, DefaultConfigPath)
	}

	ollamaURL := *flagOllamaURL
	if ollamaURL == "" {
		ollamaURL = getEnvOrDefault(EnvOllamaURL, DefaultOllamaURL)
	}

	ollamaModel := *flagOllamaModel
	if ollamaModel == "" {
		ollamaModel = getEnvOrDefault(EnvOllamaModel, DefaultOllamaModel)
	}

	apiEndpoint := *flagAPIEndpoint
	if apiEndpoint == "" {
		apiEndpoint = os.Getenv(EnvAPIEndpoint)
	}

	apiKey := *flagAPIKey
	if apiKey == "" {
		apiKey = os.Getenv(EnvAPIKey)
	}

	opts := &Options{
		ConfigPath:  configPath,
		OllamaURL:   ollamaURL,
		OllamaModel: ollamaModel,
		APIEndpoint: apiEndpoint,
		APIKey:      apiKey,
		ShowHelp:    *flagHelp,
		InputFile:   *flagInput,
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
		fmt.Println(ErrMultipleModes)
		fmt.Println(ErrModePriority)
		opts.RunMode = ModeJSONOnly
	}

	if opts.ShowHelp {
		printHelp()
		os.Exit(0)
	}

	return opts
}

func printHelp() {
	fmt.Println(HelpHeader)
	fmt.Println()
	fmt.Println(HelpUsage)
	fmt.Println()
	fmt.Println(HelpModeSection)
	fmt.Println()
	fmt.Printf(HelpCommonSection+"\n", DefaultConfigPath)
	fmt.Println()
	fmt.Printf(HelpOllamaSection+"\n", DefaultOllamaURL, DefaultOllamaModel)
	fmt.Println()
	fmt.Println(HelpAPISection)
	fmt.Println()
	fmt.Println(HelpExamples)
	fmt.Println()
}

func (o *Options) Validate() error {
	switch o.RunMode {
	case ModeOllama:
		if o.OllamaURL == "" {
			return fmt.Errorf(ErrOllamaURL)
		}
		if o.OllamaModel == "" {
			return fmt.Errorf(ErrOllamaModel)
		}
	case ModeAPI:
		if o.APIEndpoint == "" {
			return fmt.Errorf(ErrAPIEndpoint)
		}
	}
	return nil
}

func (o *Options) ModeDescription() string {
	switch o.RunMode {
	case ModeJSONOnly:
		return ModeDescJSON
	case ModeOllama:
		return fmt.Sprintf(ModeDescOllama, o.OllamaModel, o.OllamaURL)
	case ModeAPI:
		return fmt.Sprintf(ModeDescAPI, o.APIEndpoint)
	default:
		return ModeDescUnknown
	}
}
