package statistics

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"biliTagAnalyse/crawler"
)

type TagStat struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

type StatsResult struct {
	CrawlTime   string    `json:"crawl_time"`
	TotalVideos int       `json:"total_videos"`
	TotalTags   int       `json:"total_tags"`
	TagStats    []TagStat `json:"tag_stats"`
}

func CountTags(videos []*crawler.VideoInfo) *StatsResult {
	tagCount := make(map[string]int)

	for _, video := range videos {
		for _, tag := range video.Tags {
			tagCount[tag]++
		}
	}

	var tagStats []TagStat
	for tag, count := range tagCount {
		tagStats = append(tagStats, TagStat{Tag: tag, Count: count})
	}

	sort.Slice(tagStats, func(i, j int) bool {
		return tagStats[i].Count > tagStats[j].Count
	})

	totalVideos := len(videos)
	totalTags := len(tagStats)

	return &StatsResult{
		CrawlTime:   time.Now().Format("2006-01-02 15:04:05"),
		TotalVideos: totalVideos,
		TotalTags:   totalTags,
		TagStats:    tagStats,
	}
}

func CountTagsMultipleRounds(allVideos [][]*crawler.VideoInfo) *StatsResult {
	tagCount := make(map[string]int)

	for _, roundVideos := range allVideos {
		for _, video := range roundVideos {
			for _, tag := range video.Tags {
				tagCount[tag]++
			}
		}
	}

	var tagStats []TagStat
	for tag, count := range tagCount {
		tagStats = append(tagStats, TagStat{Tag: tag, Count: count})
	}

	sort.Slice(tagStats, func(i, j int) bool {
		return tagStats[i].Count > tagStats[j].Count
	})

	totalVideos := 0
	for _, roundVideos := range allVideos {
		totalVideos += len(roundVideos)
	}

	return &StatsResult{
		CrawlTime:   time.Now().Format("2006-01-02 15:04:05"),
		TotalVideos: totalVideos,
		TotalTags:   len(tagCount),
		TagStats:    tagStats,
	}
}

func SaveResults(result *StatsResult, outputPath string) error {
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
