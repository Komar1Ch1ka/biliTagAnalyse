package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"biliTagAnalyse/config"
	"biliTagAnalyse/crawler"
	"biliTagAnalyse/statistics"
)

var configPath = flag.String("config", "config.json", "配置文件路径")

func main() {
	flag.Parse()

	fmt.Println("=== B站推荐视频 Tag 分析爬虫 ===")
	log.Printf("配置文件: %s", *configPath)

	cfg, err := config.LoadConfig(*configPath)
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
	log.Printf("  - 运行模式: %s", cfg.RunMode)

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
		log.Println("未获取到任何视频数据，请检查网络连接或 Cookie 是否有效")
		os.Exit(1)
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

	if err := statistics.SaveResults(result, cfg.OutputFile); err != nil {
		log.Printf("保存结果失败: %v", err)
		os.Exit(1)
	}

	fmt.Printf("\n结果已保存到: %s\n", cfg.OutputFile)
	log.Println("=== 爬虫运行完成 ===")
}
