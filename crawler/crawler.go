package crawler

import (
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"biliTagAnalyse/utils"

	"golang.org/x/net/html"
)

type HomepageCrawler struct {
	client          *utils.HTTPClient
	retryCount      int
	retryDelay      int
	requestInterval int
	maxConcurrent   int
}

func NewHomepageCrawler(cookie string, retryCount, retryDelay, requestInterval, maxConcurrent int) *HomepageCrawler {
	return &HomepageCrawler{
		client:          utils.NewHTTPClient(cookie),
		retryCount:      retryCount,
		retryDelay:      retryDelay,
		requestInterval: requestInterval,
		maxConcurrent:   maxConcurrent,
	}
}

func (c *HomepageCrawler) CrawlHomepage() ([]string, error) {
	log.Println("正在爬取 B站 首页...")

	body, err := utils.RetryGet(c.client, "https://www.bilibili.com", c.retryCount, c.retryDelay)
	if err != nil {
		return nil, err
	}

	videos := extractVideoLinks(string(body))
	log.Printf("从首页提取到 %d 个视频链接", len(videos))

	return videos, nil
}

func extractVideoLinks(htmlContent string) []string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("解析HTML失败: %v", err)
		return nil
	}

	var videos []string
	visited := make(map[string]bool)
	var mu sync.Mutex

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			href := ""
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}

			if href != "" && strings.Contains(href, "/video/BV") {
				absURL, err := url.Parse(href)
				if err != nil {
					return
				}

				if !absURL.IsAbs() {
					base, _ := url.Parse("https://www.bilibili.com")
					absURL = base.ResolveReference(absURL)
				}
				absURL.Scheme = "https"
				cleanURL := strings.Split(strings.Split(absURL.String(), "?")[0], "#")[0]

				mu.Lock()
				if !visited[cleanURL] && strings.Contains(cleanURL, "/video/BV") {
					visited[cleanURL] = true
					videos = append(videos, cleanURL)
				}
				mu.Unlock()
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return videos
}

type VideoCrawler struct {
	client          *utils.HTTPClient
	retryCount      int
	retryDelay      int
	requestInterval int
	maxConcurrent   int
}

func NewVideoCrawler(cookie string, retryCount, retryDelay, requestInterval, maxConcurrent int) *VideoCrawler {
	return &VideoCrawler{
		client:          utils.NewHTTPClient(cookie),
		retryCount:      retryCount,
		retryDelay:      retryDelay,
		requestInterval: requestInterval,
		maxConcurrent:   maxConcurrent,
	}
}

type VideoInfo struct {
	Link   string
	Title  string
	Author string
	Tags   []string
}

func extractTagsFromHTML(htmlContent string) []string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	var tags []string
	seen := make(map[string]bool)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					classVal := attr.Val
					if strings.Contains(classVal, "tag-link") || strings.Contains(classVal, "tag") {
						text := strings.TrimSpace(getText(n))
						if text != "" && !seen[text] && len(text) > 0 && len(text) < 20 {
							seen[text] = true
							tags = append(tags, text)
						}
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return tags
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var s string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s += getText(c)
	}
	return s
}

func (c *VideoCrawler) CrawlVideosConcurrently(links []string) []*VideoInfo {
	log.Printf("开始并发爬取 %d 个视频...", len(links))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, c.maxConcurrent)
	var mu sync.Mutex
	results := make([]*VideoInfo, 0, len(links))

	for i, link := range links {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(id int, videoLink string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			time.Sleep(time.Duration(id) * 800 * time.Millisecond)

			body, err := utils.RetryGet(c.client, videoLink, c.retryCount, c.retryDelay)
			if err != nil {
				log.Printf("获取视频页面失败 %s: %v", videoLink, err)
				return
			}

			tags := extractTagsFromHTML(string(body))

			mu.Lock()
			results = append(results, &VideoInfo{
				Link: videoLink,
				Tags: tags,
			})
			mu.Unlock()
		}(i, link)
	}

	wg.Wait()

	return results
}
