package main

import (
	"context"
	"log"
	"os"

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
	cc := c.Clone()

	c.OnHTML(selector.URL, func(e *colly.HTMLElement) {
		cc.Visit(e.Attr("href"))
	})

	cc.OnHTML(selector.ArticleContainer, func(e *colly.HTMLElement) {
		article := collectArticle(e, selector)
		db.Upsert(article, collection)
	})

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

func collectArticle(e *colly.HTMLElement, s config.SelectorConfig) news.Article {
	articleURL := e.Request.URL.String()
	var timeString string
	if s.PublishedDate.Attr != "" {
		timeString = e.ChildAttr(s.PublishedDate.Css, s.PublishedDate.Attr)
	} else {
		timeString = e.ChildText(s.PublishedDate.Css)
	}
	published, _ := news.ParseDateTime(timeString, s.PublishedDate.TimeFormat)
	return news.Article{
		ID:        db.CreateObjectID(articleURL, published),
		Title:     e.ChildText(s.Title),
		URL:       articleURL,
		Published: published,
		Content:   news.CleanContent(e, s.Content),
	}
}
