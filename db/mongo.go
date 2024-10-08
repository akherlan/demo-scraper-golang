package db

import (
	"context"
	"crypto/md5"
	"log"
	"scraper/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Upsert(item model.NewsArticle, c *mongo.Collection) {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": item.ID}
	update := bson.M{"$set": item}
	_, err := c.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Error insert: %v", err)
	} else {
		log.Printf("Update news - %s", item.Title)
	}
}

func DefineObjectID(id string) primitive.ObjectID {
	oid := [12]byte{}
	idHash := md5.Sum([]byte(id))
	copy(oid[:], idHash[:12])
	return primitive.ObjectID(oid)
}
