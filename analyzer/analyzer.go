package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"biliTagAnalyse/cmd"
	"biliTagAnalyse/statistics"
)

type Analyzer struct {
	opts *cmd.Options
}

func NewAnalyzer(opts *cmd.Options) *Analyzer {
	return &Analyzer{opts: opts}
}

type AnalysisResult struct {
	Summary     string                 `json:"summary"`
	TopTags     []TagInsight           `json:"top_tags_insights"`
	Trends      []string               `json:"trends"`
	Suggestions []string               `json:"suggestions"`
	RawStats    *statistics.StatsResult `json:"raw_stats"`
}

type TagInsight struct {
	Tag         string `json:"tag"`
	Count       int    `json:"count"`
	Description string `json:"description"`
}

func (a *Analyzer) Analyze(stats *statistics.StatsResult) (*AnalysisResult, error) {
	switch a.opts.RunMode {
	case cmd.ModeJSONOnly:
		return a.analyzeJSONOnly(stats)
	case cmd.ModeOllama:
		return a.analyzeWithOllama(stats)
	case cmd.ModeAPI:
		return a.analyzeWithAPI(stats)
	default:
		return a.analyzeJSONOnly(stats)
	}
}

func (a *Analyzer) analyzeJSONOnly(stats *statistics.StatsResult) (*AnalysisResult, error) {
	log.Println("运行模式：JSON文件输出模式（不进行模型分析）")
	
	result := &AnalysisResult{
		Summary:     "仅JSON输出模式，未进行模型分析",
		TopTags:     make([]TagInsight, 0),
		Trends:      []string{"需要模型分析才能生成趋势洞察"},
		Suggestions: []string{"使用 -ollama 或 -api 参数进行深度分析"},
		RawStats:    stats,
	}

	for i := 0; i < len(stats.TagStats) && i < 10; i++ {
		result.TopTags = append(result.TopTags, TagInsight{
			Tag:         stats.TagStats[i].Tag,
			Count:       stats.TagStats[i].Count,
			Description: "需要模型分析生成描述",
		})
	}

	return result, nil
}

func (a *Analyzer) analyzeWithOllama(stats *statistics.StatsResult) (*AnalysisResult, error) {
	log.Printf("运行模式：Ollama本地模型分析 (模型: %s, 地址: %s)", a.opts.OllamaModel, a.opts.OllamaURL)

	prompt := a.buildAnalysisPrompt(stats)
	
	response, err := a.callOllamaAPI(prompt)
	if err != nil {
		return nil, fmt.Errorf("Ollama分析失败: %w", err)
	}

	result := &AnalysisResult{
		Summary:     response,
		TopTags:     make([]TagInsight, 0),
		Trends:      []string{},
		Suggestions: []string{},
		RawStats:    stats,
	}

	for i := 0; i < len(stats.TagStats) && i < 10; i++ {
		result.TopTags = append(result.TopTags, TagInsight{
			Tag:         stats.TagStats[i].Tag,
			Count:       stats.TagStats[i].Count,
			Description: "由Ollama模型分析",
		})
	}

	return result, nil
}

func (a *Analyzer) analyzeWithAPI(stats *statistics.StatsResult) (*AnalysisResult, error) {
	log.Printf("运行模式：远程API调用分析 (端点: %s)", a.opts.APIEndpoint)

	prompt := a.buildAnalysisPrompt(stats)
	
	response, err := a.callRemoteAPI(prompt)
	if err != nil {
		return nil, fmt.Errorf("API分析失败: %w", err)
	}

	result := &AnalysisResult{
		Summary:     response,
		TopTags:     make([]TagInsight, 0),
		Trends:      []string{},
		Suggestions: []string{},
		RawStats:    stats,
	}

	for i := 0; i < len(stats.TagStats) && i < 10; i++ {
		result.TopTags = append(result.TopTags, TagInsight{
			Tag:         stats.TagStats[i].Tag,
			Count:       stats.TagStats[i].Count,
			Description: "由远程API分析",
		})
	}

	return result, nil
}

func (a *Analyzer) buildAnalysisPrompt(stats *statistics.StatsResult) string {
	topTags := ""
	for i := 0; i < len(stats.TagStats) && i < 20; i++ {
		topTags += fmt.Sprintf("%d. %s (出现次数: %d)\n", i+1, stats.TagStats[i].Tag, stats.TagStats[i].Count)
	}

	prompt := fmt.Sprintf(`你是一个B站视频内容分析专家。请分析以下B站推荐视频的Tag统计数据，提供专业的内容洞察。

统计信息：
- 爬取时间: %s
- 总视频数: %d
- 不同Tag数: %d

Top 20 Tags:
%s

请从以下角度进行分析：
1. 内容趋势：这些Tag反映了什么样的内容趋势？
2. 用户偏好：用户对什么类型的内容更感兴趣？
3. 创作建议：对于UP主有什么创作建议？

请用中文回答，保持简洁专业。`, 
		stats.CrawlTime,
		stats.TotalVideos,
		stats.TotalTags,
		topTags,
	)

	return prompt
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func (a *Analyzer) callOllamaAPI(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  a.opts.OllamaModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	url := a.opts.OllamaURL + "/api/generate"
	
	client := &http.Client{Timeout: 120 * time.Second}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("正在调用Ollama API: %s", url)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求Ollama失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return ollamaResp.Response, nil
}

type APIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type APIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (a *Analyzer) callRemoteAPI(prompt string) (string, error) {
	reqBody := APIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	
	req, err := http.NewRequest("POST", a.opts.APIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if a.opts.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+a.opts.APIKey)
	}

	log.Printf("正在调用远程API: %s", a.opts.APIEndpoint)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求远程API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if apiResp.Error != nil {
		return "", fmt.Errorf("API错误: %s", apiResp.Error.Message)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return apiResp.Choices[0].Message.Content, nil
}

func LoadStatsFromFile(path string) (*statistics.StatsResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	var stats statistics.StatsResult
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return &stats, nil
}

func SaveAnalysisResult(result *AnalysisResult, outputPath string) error {
	dir := outputPath
	lastSlash := -1
	for i := len(outputPath) - 1; i >= 0; i-- {
		if outputPath[i] == '/' || outputPath[i] == '\\' {
			lastSlash = i
			break
		}
	}
	if lastSlash > 0 {
		dir = outputPath[:lastSlash]
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化结果失败: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
