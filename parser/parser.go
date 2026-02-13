package parser

import (
	"regexp"
	"strings"
)

type VideoParser struct{}

func NewVideoParser() *VideoParser {
	return &VideoParser{}
}

func (p *VideoParser) ExtractBVNumber(link string) string {
	re := regexp.MustCompile(`[Bb][Vv][A-Za-z0-9]+`)
	match := re.FindString(link)
	return match
}

func (p *VideoParser) ExtractTitle(html string) string {
	patterns := []string{
		`<h1[^>]*class="video-title[^"]*"[^>]*>([^<]+)</h1>`,
		`"title":"([^"]+)"`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) > 1 {
				title := strings.TrimSpace(match[1])
				if title != "" {
					return title
				}
			}
		}
	}

	return ""
}

func (p *VideoParser) ExtractAuthor(html string) string {
	patterns := []string{
		`<a[^>]*class="up-name[^"]*"[^>]*>([^<]+)</a>`,
		`"author":"([^"]+)"`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) > 1 {
				author := strings.TrimSpace(match[1])
				if author != "" {
					return author
				}
			}
		}
	}

	return ""
}
