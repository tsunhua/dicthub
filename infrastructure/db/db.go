package db

import (
	"app/infrastructure/config"
	"app/infrastructure/log"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var dbConfig = config.Get().DB

var mongoInstance = struct {
	client *mongo.Client
	err    error
}{}

var once = &sync.Once{}

func init() {
	go initMongoOnce()
}

func initMongoOnce() {
	once.Do(func() {
		log.Debug("start connect mongo")
		ctx, cancel := context.WithTimeout(context.Background(), dbConfig.ConnectionTimeout)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.URI))
		if err != nil {
			mongoInstance.err = err
			mongoInstance.client = nil
			log.Error(fmt.Sprintf("cannot connect to mongo, err:%v", mongoInstance.err))
		} else {
			mongoInstance.client = client
		}
		log.Debug(fmt.Sprintf("finish connect mongo, err:%v", mongoInstance.err))
	})
}

func mongoClient() (*mongo.Client, error) {
	initMongoOnce()
	return mongoInstance.client, mongoInstance.err
}

func GetTable(name string) (*mongo.Collection, error) {
	client, err := mongoClient()
	if err != nil {
		return nil, err
	}
	return client.Database(dbConfig.Name).Collection(name), nil
}
