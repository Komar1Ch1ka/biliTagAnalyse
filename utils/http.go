package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HTTPClient struct {
	Client *http.Client
	Cookie string
}

func NewHTTPClient(cookie string) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("重定向次数过多")
				}
				return nil
			},
		},
		Cookie: cookie,
	}
}

func (c *HTTPClient) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	if c.Cookie != "" {
		req.Header.Set("Cookie", c.Cookie)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}

func RetryGet(client *HTTPClient, url string, retryCount int, retryDelay int) ([]byte, error) {
	var err error
	var body []byte

	for i := 0; i < retryCount; i++ {
		body, err = client.Get(url)
		if err == nil {
			return body, nil
		}

		log.Printf("请求失败 (尝试 %d/%d): %v", i+1, retryCount, err)

		if i < retryCount-1 {
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}
	}

	return nil, fmt.Errorf("重试 %d 次后仍然失败: %w", retryCount, err)
}
