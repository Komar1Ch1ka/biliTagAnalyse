package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"biliTagAnalyse/cmd"
)

type Config struct {
	Cookie          string `json:"cookie"`
	CrawlCount      int    `json:"crawl_count"`
	CrawlInterval   int    `json:"crawl_interval"`
	RequestInterval int    `json:"request_interval"`
	MaxConcurrent   int    `json:"max_concurrent"`
	RetryCount      int    `json:"retry_count"`
	RetryDelay      int    `json:"retry_delay"`
	OutputFile      string `json:"output_file"`
	RunMode         string `json:"run_mode"`
	OllamaURL       string `json:"ollama_url"`
	OllamaModel     string `json:"ollama_model"`
	APIEndpoint     string `json:"api_endpoint"`
	APIKey          string `json:"api_key"`
}

func ResolveConfigPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	if _, err := os.Stat(path); err == nil {
		absPath, _ := filepath.Abs(path)
		return absPath
	}

	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		execConfigPath := filepath.Join(execDir, path)
		if _, err := os.Stat(execConfigPath); err == nil {
			return execConfigPath
		}
	}

	absPath, _ := filepath.Abs(path)
	return absPath
}

func LoadConfig(path string) (*Config, error) {
	resolvedPath := ResolveConfigPath(path)
	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if cfg.Cookie == "" || cfg.Cookie == "你的B站Cookie" {
		return nil, fmt.Errorf("请在 config.json 中设置有效的 B站 Cookie")
	}

	if cfg.CrawlCount <= 0 {
		cfg.CrawlCount = 1
	}
	if cfg.CrawlInterval <= 0 {
		cfg.CrawlInterval = 300
	}
	if cfg.RequestInterval <= 0 {
		cfg.RequestInterval = 500
	}
	if cfg.MaxConcurrent <= 0 {
		cfg.MaxConcurrent = 3
	}
	if cfg.RetryCount <= 0 {
		cfg.RetryCount = 3
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = 2
	}
	if cfg.OutputFile == "" {
		cfg.OutputFile = "results/tags_stats.json"
	}
	if cfg.RunMode == "" {
		cfg.RunMode = "once"
	}
	if cfg.OllamaURL == "" {
		cfg.OllamaURL = cmd.DefaultOllamaURL
	}
	if cfg.OllamaModel == "" {
		cfg.OllamaModel = cmd.DefaultOllamaModel
	}

	return &cfg, nil
}
