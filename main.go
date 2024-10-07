package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"scraper/config"
	"scraper/db"
	"scraper/news"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := db.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(cfg.MongoDBName).Collection(cfg.MongoCollection)
	scrape("detik", collection, cfg)
	scrape("liputan6", collection, cfg)
}

func scrape(source string, collection *mongo.Collection, cfg *config.Config) {
	selector := cfg.Selectors[source]
	allowedDomain := cfg.Sources[source].Domain

	c := colly.NewCollector(colly.AllowedDomains(allowedDomain))
	c.UserAgent = cfg.Scraper.UserAgent
	c.SetRequestTimeout(time.Duration(cfg.Scraper.Timeout) * time.Second)
	cc := c.Clone()

	c.OnHTML(selector.URL, func(e *colly.HTMLElement) {
		cc.Visit(e.Attr("href"))
	})

	if selector.Jsonld != "" {
		// JSON-LD content parsing
		cc.OnHTML(selector.Jsonld, func(e *colly.HTMLElement) {
			object, err := news.FromJsonLdString(e.Text)
			if len(object) > 0 {
				article := collectArticleJsonLD(object[0])
				db.Upsert(article, collection)
			} else {
				log.Println("Empty JSON-LD:", err)
			}
		})
	} else {
		// HTML content parsing
		cc.OnHTML(selector.ArticleContainer, func(e *colly.HTMLElement) {
			article := collectArticleHTML(e, selector)
			isMultiplePages := news.DetectPagination(e, selector.PageIndex)
			if isMultiplePages {
				log.Println("Detect pagination:", article.URL)
				updatedURL := e.Request.URL.String() + "?single=1"
				cc.Visit(updatedURL)
			} else {
				// save if all-pages link accessed
				db.Upsert(article, collection)
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

	err := c.Visit(cfg.Sources[source].StartURL)
	if err != nil {
		log.Fatal(err)
	}
}

func collectArticleHTML(e *colly.HTMLElement, s config.SelectorConfig) news.Article {
	timeFormat := s.PublishedDate.TimeFormat
	articleURL := e.Request.URL.String()
	articleID := news.GetIDFromURL(articleURL)
	dtString := news.ParseDatePublished(e, s)
	published, _ := news.ConvertDateTime(dtString, timeFormat)
	return news.Article{
		ID:        db.DefineObjectID(articleID),
		Title:     e.ChildText(s.Title),
		URL:       articleURL,
		Published: published,
		Content:   news.CleanContentHTML(e, s.Content),
	}
}

func collectArticleJsonLD(j news.JsonLD) news.Article {
	timeFormat := "2006-01-02T15:04:05-07:00"
	articleID := news.GetIDFromURL(j.URL)
	published, _ := news.ConvertDateTime(j.Published, timeFormat)
	return news.Article{
		ID:        db.DefineObjectID(articleID),
		Title:     j.Title,
		URL:       j.URL,
		Published: published,
		Content:   news.CleanContentLiputan6(j.Content),
	}
}
