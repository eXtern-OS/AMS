package AMS

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"sync"
	"time"
)

var URI = ""
var SQL_URI = ""

func random(low, hi int) int {
	return low + rand.Intn(hi-low)
}

type DBConn struct {
	Mutex  sync.Mutex
	Client *mongo.Client
}

var dbc DBConn

func initDb() {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	dbc.Mutex.Lock()
	dbc.Client = client
	dbc.Mutex.Unlock()
}

func (d *DBConn) newCollection(collname string) *mongo.Collection {
	d.Mutex.Lock()
	coll := d.Client.Database("Users").Collection(collname)
	d.Mutex.Unlock()
	return coll
}

func newCollection(collname string) *mongo.Collection {
	return dbc.newCollection(collname)
}
