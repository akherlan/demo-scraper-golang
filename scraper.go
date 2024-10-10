package main

import (
	"log"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"

	"scraper/config"
	"scraper/db"
	"scraper/model"
)

type Source interface {
	Scrape(coll *mongo.Collection, cfg *config.Config)
}

type Website struct {
	Name     string `mapstructure:"name"`
	Domain   string `mapstructure:"domain"`
	StartURL string `mapstructure:"startURL"`
	Method   string `mapstructure:"method"`
}

func LoadSource() ([]Website, error) {
	viper.SetConfigFile("./source.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var target []Website
	err = viper.UnmarshalKey("source", &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (site Website) Scrape(coll *mongo.Collection, cfg *config.Config) {
	scrape_news(site, coll, cfg)
}

func scrape_news(site Website, coll *mongo.Collection, cfg *config.Config) {
	selector := cfg.Selectors[site.Name]
	c := colly.NewCollector(colly.AllowedDomains(site.Domain))
	c.UserAgent = cfg.Scraper.UserAgent
	c.SetRequestTimeout(time.Duration(cfg.Scraper.Timeout) * time.Second)
	cc := c.Clone()

	c.OnHTML(selector.URL, func(e *colly.HTMLElement) {
		cc.Visit(e.Attr("href"))
	})

	switch site.Method {
	case "jsonld":
		cc.OnHTML(selector.Jsonld, func(e *colly.HTMLElement) {
			object, err := FromJsonLdString(e.Text)
			if len(object) > 0 {
				article := collectArticleJsonLD(object[0])
				db.Upsert(article, coll)
			} else {
				log.Println("Empty JSON-LD:", err)
			}
		})
	default:
		cc.OnHTML(selector.ArticleContainer, func(e *colly.HTMLElement) {
			article := collectArticleHTML(e, selector)
			isMultiplePages := DetectPagination(e, selector.PageIndex)
			if isMultiplePages {
				log.Println("Detect pagination:", article.URL)
				updatedURL := e.Request.URL.String() + "?single=1"
				cc.Visit(updatedURL)
			} else {
				// save if all-pages link accessed
				db.Upsert(article, coll)
			}
		})
	}

	c.OnResponse(func(r *colly.Response) {
		log.Println(r.StatusCode, r.Request.URL)
	})

	cc.OnResponse(func(r *colly.Response) {
		log.Println(r.StatusCode, r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatalf("URL failed %s with response %v\nError: %v", r.Request.URL, r, err)
	})

	cc.OnError(func(r *colly.Response, err error) {
		log.Fatalf("URL failed %s with response %v\nError: %v", r.Request.URL, r, err)
	})

	err := c.Visit(site.StartURL)
	if err != nil {
		log.Fatal(err)
	}
}

func collectArticleHTML(e *colly.HTMLElement, s config.Selector) model.NewsArticle {
	timeFormat := s.PublishedDate.TimeFormat
	articleURL := e.Request.URL.String()
	articleID := GetIDFromURL(articleURL)
	dtString := ParseDatePublished(e, s)
	published, _ := ConvertDateTime(dtString, timeFormat)
	return model.NewsArticle{
		ID:        db.DefineObjectID(articleID),
		Title:     e.ChildText(s.Title),
		URL:       articleURL,
		Published: published,
		Content:   CleanContentHTML(e, s.Content),
	}
}

func collectArticleJsonLD(j model.NewsArticleJsonLD) model.NewsArticle {
	timeFormat := "2006-01-02T15:04:05-07:00"
	published, _ := ConvertDateTime(j.Published, timeFormat)
	return model.NewsArticle{
		ID:        db.DefineObjectID(GetIDFromURL(j.URL)),
		Title:     j.Title,
		URL:       j.URL,
		Published: published,
		Content:   CleanContentLiputan6(j.Content),
	}
}
