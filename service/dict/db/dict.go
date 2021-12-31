package db

import (
	"app/infrastructure/config"
	"app/infrastructure/db"
	"app/service/dict/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindRecentDictList(limit int64, offset int64) (dicts []*model.Dict, err error) {
	cur, err := getDictTable().Find(context.TODO(), bson.M{"isPublic": true},
		options.Find().SetSort(bson.M{"updateTime": -1}).SetLimit(limit).SetSkip(offset))
	if err != nil {
		return
	}
	dicts = make([]*model.Dict, 0, limit)
	for cur.Next(context.TODO()) {
		item := &model.Dict{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		dicts = append(dicts, item)
	}
	return
}

func FindDictById(id string) (dict *model.Dict, err error) {
	err = getDictTable().FindOne(context.TODO(), bson.M{"id": id}).Decode(&dict)
	return
}

func FindDictsBy(tags []string) (dicts []*model.Dict, err error) {
	filter := bson.M{}
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$all": tags}
	}
	cur, err := getDictTable().Find(context.TODO(), filter,
		options.Find().SetSort(bson.M{"updateTime": 1}).SetLimit(100))
	if err != nil {
		return
	}
	dicts = make([]*model.Dict, 0, 10)
	for cur.Next(context.TODO()) {
		item := &model.Dict{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		dicts = append(dicts, item)
	}
	return
}

func FindManyDictsById(ids []string) (dicts []*model.Dict, err error) {
	cur, err := getDictTable().Find(context.TODO(), bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return
	}
	dicts = make([]*model.Dict, 0, 3)
	for cur.Next(context.TODO()) {
		item := &model.Dict{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		dicts = append(dicts, item)
	}
	return
}

func InsertDict(dict *model.Dict) (err error) {
	_, err = getDictTable().InsertOne(context.TODO(), dict)
	return
}

func UpdateDict(dict *model.Dict) (err error) {
	update := bson.M{
		"updateTime": dict.UpdateTime,
	}
	if dict.Name != "" {
		update["name"] = dict.Name
	}
	if dict.Desc != "" {
		update["desc"] = dict.Desc
	}
	if dict.FeedbackEmail != "" {
		update["feedbackEmail"] = dict.FeedbackEmail
	}
	if dict.Tags != nil {
		update["tags"] = dict.Tags
	}
	if dict.CatalogText != "" {
		update["catalogText"] = dict.CatalogText
	}
	if dict.SpecText != "" {
		update["specText"] = dict.SpecText
	}
	if dict.PreferSpecLinkIds != nil {
		update["preferSpecLinkIds"] = dict.PreferSpecLinkIds
	}
	_, err = getDictTable().UpdateOne(context.TODO(), bson.M{"id": dict.Id}, bson.M{"$set": update})
	return
}

func getDictTable() *mongo.Collection {
	table, err := db.GetTable(config.Get().DB.DictTable.Name)
	if err != nil {
		panic(err)
	}
	_, err = table.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"isPublic": 1},
			Options: options.Index().SetUnique(false),
		},
	})
	if err != nil {
		panic(err)
	}
	return table
}
