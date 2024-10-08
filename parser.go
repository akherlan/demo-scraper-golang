package main

import (
	"encoding/json"
	"html"
	"regexp"
	"scraper/config"
	"scraper/model"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

var monthsAbbr = map[string]string{
	"Jan": "Jan", "Feb": "Feb", "Mar": "Mar", "Apr": "Apr",
	"Mei": "May", "Jun": "Jun", "Jul": "Jul", "Agu": "Aug",
	"Sep": "Sep", "Okt": "Oct", "Nov": "Nov", "Des": "Dec",
}

func FromJsonLdString(text string) ([]model.NewsArticleJsonLD, error) {
	var objects []model.NewsArticleJsonLD
	var data []model.NewsArticleJsonLD
	err := json.Unmarshal([]byte(text), &objects)
	if err != nil {
		return nil, err
	}
	for _, item := range objects {
		if item.Type == "NewsArticle" || item.Type == "Article" {
			data = append(data, item)
		}
	}
	return data, nil
}

func CleanContentHTML(e *colly.HTMLElement, selector string) string {
	cfg, _ := config.Load()
	clone := e.DOM.Clone()
	content := clone.Find(selector)
	cssExclude := strings.Join(cfg.CssContentExclude, ", ")
	content.Find(cssExclude).Remove()
	return trimAllString(content.Text())
}

func CleanContentLiputan6(text string) string {
	text = removeRecommendation(text)
	text = removeMedia(text)
	text = html.UnescapeString(text)
	text = trimAllString(text)
	return text
}

func trimAllString(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	text = ReplacePattern(text, `\s+`, " ")
	text = strings.TrimSpace(text)
	return text
}

func removeRecommendation(text string) string {
	// Liputan6
	// [bacajuga:Baca Juga](digits)
	articlePttrn := `\[bacajuga:Baca Juga\]\(\d+(?:\s+\d+)*\)`
	return ReplacePattern(text, articlePttrn, " ")
}

func removeMedia(text string) string {
	// Liputan6
	// [vidio:Judul](https://link)
	mediaPttrn := `(?:Simak Video Pilihan Ini:)?\[vidio:[^\]]+\]\(https?://[^\)]+\)`
	return ReplacePattern(text, mediaPttrn, " ")
}

func ReplacePattern(text string, pattern string, repl string) string {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(text) {
		return text
	}
	return re.ReplaceAllString(text, repl)
}

func DetectPagination(e *colly.HTMLElement, selector string) bool {
	if e.DOM.Find(selector).Length() > 0 {
		return true
	} else {
		return false
	}
}

func ConvertDateTime(timeStr string, layout string) (time.Time, error) {
	tz := "Asia/Jakarta"
	parts := strings.SplitN(timeStr, ", ", 2)
	if len(parts) > 1 {
		timeStr = strings.TrimSpace(parts[1])
	}
	for ind, eng := range monthsAbbr {
		timeStr = strings.Replace(timeStr, ind, eng, 1)
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ParseDatePublished(e *colly.HTMLElement, s config.Selector) string {
	var timeString string
	if s.PublishedDate.Attr != "" {
		timeString = e.ChildAttr(s.PublishedDate.Css, s.PublishedDate.Attr)
	} else {
		timeString = e.ChildText(s.PublishedDate.Css)
	}
	return timeString
}

func GetIDFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return ""
}
