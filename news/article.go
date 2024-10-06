package news

import (
	"regexp"
	"scraper/config"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	URL       string             `json:"url" bson:"url"`
	Published time.Time          `json:"published" bson:"published"`
	Content   string             `json:"content" bson:"content"`
}

func CleanContent(e *colly.HTMLElement, selector string) string {
	cfg, _ := config.Load()
	clone := e.DOM.Clone()
	content := clone.Find(selector)
	cssExclude := strings.Join(cfg.CssContentExclude, ", ")
	content.Find(cssExclude).Remove()
	return TrimAllString(content.Text())
}

func TrimAllString(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	return text
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
	months := map[string]string{
		"Jan": "Jan", "Feb": "Feb", "Mar": "Mar", "Apr": "Apr",
		"Mei": "May", "Jun": "Jun", "Jul": "Jul", "Agu": "Aug",
		"Sep": "Sep", "Okt": "Oct", "Nov": "Nov", "Des": "Dec",
	}
	parts := strings.SplitN(timeStr, ", ", 2)
	if len(parts) > 1 {
		timeStr = strings.TrimSpace(parts[1])
	}
	for ind, eng := range months {
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
