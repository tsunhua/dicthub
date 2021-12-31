package db

import (
	"app/infrastructure/cache"
	"app/infrastructure/config"
	"app/infrastructure/db"
	"app/infrastructure/log"
	"app/service/user/model"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllUser() (users []*model.User, err error) {
	cur, err := getUserTable().Find(context.TODO(), bson.M{})
	if err != nil {
		return
	}
	users = make([]*model.User, 0, 100)
	for cur.Next(context.TODO()) {
		user := &model.User{}
		err = cur.Decode(user)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

func FindUser(email string) (user *model.User, err error) {
	ckey := "/user/" + email
	cword, err := cache.Cache().Get(ckey)
	if err == nil {
		user = cword.(*model.User)
		return
	}
	log.Debug(err.Error())
	result := getUserTable().FindOne(context.TODO(), bson.M{"email": email})
	if result == nil {
		err = errors.New("user not found")
		return
	}
	user = &model.User{}
	err = result.Decode(user)
	if err == nil {
		cache.Cache().Set(ckey, user)
	}
	return
}

func getUserTable() *mongo.Collection {
	table, err := db.GetTable(config.Get().DB.UserTable.Name)
	if err != nil {
		panic(err)
	}
	_, err = table.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(false),
		},
	})
	if err != nil {
		panic(err)
	}
	return table
}
