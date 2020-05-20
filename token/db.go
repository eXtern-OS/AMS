package token

import (
	"context"
	beatrix "github.com/eXtern-OS/Beatrix"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var URI = ""

func Init(mongoUri string) {
	URI = mongoUri
}

func PutToken(token Token) int {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error creating new mongo client", "PUTTOKEN")
		return http.StatusInternalServerError
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error connecting with new mongo client", "PUTTOKEN")
		return http.StatusInternalServerError
	}

	collection := client.Database("Users").Collection("token")
	_, err = collection.InsertOne(context.Background(), token)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error inserting token", "PUTTOKEN")
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
