package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"biliTagAnalyse/analyzer"
	"biliTagAnalyse/cmd"
	"biliTagAnalyse/config"
	"biliTagAnalyse/crawler"
	"biliTagAnalyse/statistics"
)

func main() {
	opts := cmd.Parse()

	if err := opts.Validate(); err != nil {
		log.Fatalf("参数验证失败: %v", err)
	}

	fmt.Println("=== B站推荐视频 Tag 分析爬虫 ===")
	log.Printf("配置文件: %s", opts.ConfigPath)
	log.Printf("运行模式: %s", opts.ModeDescription())

	cfg, err := config.LoadConfig(opts.ConfigPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	log.Printf("配置信息:")
	log.Printf("  - 爬取次数: %d", cfg.CrawlCount)
	log.Printf("  - 爬取间隔: %d 秒", cfg.CrawlInterval)
	log.Printf("  - 请求间隔: %d 毫秒", cfg.RequestInterval)
	log.Printf("  - 最大并发: %d", cfg.MaxConcurrent)
	log.Printf("  - 重试次数: %d", cfg.RetryCount)
	log.Printf("  - 输出文件: %s", cfg.OutputFile)

	var statsResult *statistics.StatsResult

	if opts.InputFile != "" && (opts.RunMode == cmd.ModeOllama || opts.RunMode == cmd.ModeAPI) {
		log.Printf("从文件加载数据: %s", opts.InputFile)
		statsResult, err = analyzer.LoadStatsFromFile(opts.InputFile)
		if err != nil {
			log.Fatalf("加载输入文件失败: %v", err)
		}
	} else {
		statsResult, err = runCrawler(cfg)
		if err != nil {
			log.Fatalf("爬取失败: %v", err)
		}
	}

	log.Println("\n=== 执行分析 ===")
	az := analyzer.NewAnalyzer(opts)
	analysisResult, err := az.Analyze(statsResult)
	if err != nil {
		log.Fatalf("分析失败: %v", err)
	}

	outputPath := cfg.OutputFile
	if opts.RunMode != cmd.ModeJSONOnly {
		outputPath = "results/analysis_result.json"
	}

	if err := saveResults(analysisResult, outputPath, opts.RunMode); err != nil {
		log.Fatalf("保存结果失败: %v", err)
	}

	fmt.Printf("\n结果已保存到: %s\n", outputPath)
	printAnalysisSummary(analysisResult, opts.RunMode)
	log.Println("=== 程序运行完成 ===")
}

func runCrawler(cfg *config.Config) (*statistics.StatsResult, error) {
	homepageCrawler := crawler.NewHomepageCrawler(
		cfg.Cookie,
		cfg.RetryCount,
		cfg.RetryDelay,
		cfg.RequestInterval,
		cfg.MaxConcurrent,
	)

	videoCrawler := crawler.NewVideoCrawler(
		cfg.Cookie,
		cfg.RetryCount,
		cfg.RetryDelay,
		cfg.RequestInterval,
		cfg.MaxConcurrent,
	)

	var allVideos [][]*crawler.VideoInfo

	for i := 0; i < cfg.CrawlCount; i++ {
		log.Printf("\n--- 第 %d/%d 轮爬取 ---", i+1, cfg.CrawlCount)

		links, err := homepageCrawler.CrawlHomepage()
		if err != nil {
			log.Printf("爬取首页失败: %v", err)
			continue
		}

		if len(links) == 0 {
			log.Println("未获取到视频链接，可能需要使用无头浏览器")
			continue
		}

		videos := videoCrawler.CrawlVideosConcurrently(links)
		allVideos = append(allVideos, videos)

		log.Printf("第 %d 轮爬取完成，获取到 %d 个视频", i+1, len(videos))

		if i < cfg.CrawlCount-1 {
			log.Printf("等待 %d 秒后进行下一轮爬取...", cfg.CrawlInterval)
			time.Sleep(time.Duration(cfg.CrawlInterval) * time.Second)
		}
	}

	if len(allVideos) == 0 {
		return nil, fmt.Errorf("未获取到任何视频数据，请检查网络连接或 Cookie 是否有效")
	}

	log.Println("\n=== 统计 Tag ===")
	result := statistics.CountTagsMultipleRounds(allVideos)

	log.Printf("统计结果:")
	log.Printf("  - 总视频数: %d", result.TotalVideos)
	log.Printf("  - 总 Tag 数: %d", result.TotalTags)
	log.Printf("  - Top 10 Tags:")
	for i := 0; i < len(result.TagStats) && i < 10; i++ {
		log.Printf("    %d. %s (%d)", i+1, result.TagStats[i].Tag, result.TagStats[i].Count)
	}

	return result, nil
}

func saveResults(result *analyzer.AnalysisResult, outputPath string, mode cmd.RunMode) error {
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

	if mode == cmd.ModeJSONOnly {
		return statistics.SaveResults(result.RawStats, outputPath)
	}

	return analyzer.SaveAnalysisResult(result, outputPath)
}

func printAnalysisSummary(result *analyzer.AnalysisResult, mode cmd.RunMode) {
	fmt.Println("\n=== 分析摘要 ===")
	fmt.Printf("总视频数: %d\n", result.RawStats.TotalVideos)
	fmt.Printf("总Tag数: %d\n", result.RawStats.TotalTags)
	
	fmt.Println("\nTop 5 Tags:")
	for i := 0; i < len(result.TopTags) && i < 5; i++ {
		fmt.Printf("  %d. %s (次数: %d)\n", i+1, result.TopTags[i].Tag, result.TopTags[i].Count)
	}

	if mode != cmd.ModeJSONOnly && result.Summary != "" {
		fmt.Println("\n模型分析结果:")
		fmt.Println(result.Summary)
	}
}
