package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func cleanText(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
}

func extractLink(s *goquery.Selection) string {
	link, exists := s.Find("a.poly-component__title").Attr("href")
	if !exists {
		return ""
	}
	return link
}
