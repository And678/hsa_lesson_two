package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func connectToMongo() (*mongo.Client, context.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:password@127.0.0.1:27017/"))
	if err != nil {
		panic(err)
	}
	return client, ctx
}
func main() {
	// ElasticSearch
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	// Gin
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Endpoints
	r.GET("/mongo", func(c *gin.Context) {
		mongo, mongoCtx := connectToMongo()
		defer mongo.Disconnect(mongoCtx)
		quickstartDatabase := mongo.Database("lessonTwo")
		podcastsCollection := quickstartDatabase.Collection("podcasts")
		cursor, err := podcastsCollection.Find(mongoCtx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		var podcasts []bson.M
		if err = cursor.All(mongoCtx, &podcasts); err != nil {
			log.Fatal(err)
		}
		fmt.Println(podcasts)
		c.JSON(200, podcasts)
	})


	r.GET("/elastic", func(c *gin.Context) {
		var data map[string]interface{}
		// Perform the search request.
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex("business-data"),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
		)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		c.JSON(200, data)
	})
	r.Run("0.0.0.0:4222") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
