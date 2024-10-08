package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsArticle struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	URL       string             `json:"url" bson:"url"`
	Published time.Time          `json:"published" bson:"published"`
	Content   string             `json:"content" bson:"content"`
}

type NewsArticleJsonLD struct {
	Type      string `json:"@type"`
	Title     string `json:"headline"`
	URL       string `json:"mainEntityOfPage"`
	Published string `json:"datePublished"`
	Content   string `json:"articleBody"`
}
