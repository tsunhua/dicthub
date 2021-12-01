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

func FindRecentWordList(limit int64, offset int64) (words []*model.Word, err error) {
	cur, err := getWordTable().Find(context.TODO(), bson.M{},
		options.Find().SetSort(bson.M{"updateTime": -1}).SetLimit(limit).SetSkip(offset))
	if err != nil {
		return
	}
	words = make([]*model.Word, 0, limit)
	for cur.Next(context.TODO()) {
		item := &model.Word{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		words = append(words, item)
	}
	return
}

func FindWordById(id string) (word *model.Word, err error) {
	err = getWordTable().FindOne(context.TODO(), bson.M{"id": id}).Decode(&word)
	return
}

func FindManyWordsById(ids []string) (words []*model.Word, err error) {
	cur, err := getWordTable().Find(context.TODO(), bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return
	}
	words = make([]*model.Word, 0, 3)
	for cur.Next(context.TODO()) {
		item := &model.Word{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		words = append(words, item)
	}
	return
}

func FindWordByWriting(writing string) (words []*model.Word, err error) {
	cur, err := getWordTable().Find(context.TODO(), bson.M{"writing": writing})
	if err != nil {
		return
	}
	words = make([]*model.Word, 0, 3)
	for cur.Next(context.TODO()) {
		item := &model.Word{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		words = append(words, item)
	}
	return
}

func FindWordsBy(catalogLinkId, dictId string) (words []*model.Word, err error) {
	filter := bson.M{}
	if catalogLinkId != "" {
		filter["catalogLinkIds"] = bson.M{"$all": []string{catalogLinkId}}
	}
	if dictId != "" {
		filter["dictId"] = dictId
	}
	cur, err := getWordTable().Find(context.TODO(), filter,
		options.Find().SetSort(bson.M{"writing": 1}).SetLimit(100))
	if err != nil {
		return
	}
	words = make([]*model.Word, 0, 10)
	for cur.Next(context.TODO()) {
		item := &model.Word{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		words = append(words, item)
	}
	return
}

func FindBriefWordsBy(catalogLinkId, dictId string) (words []*model.Word, err error) {
	filter := bson.M{}
	if catalogLinkId != "" {
		filter["catalogLinkIds"] = bson.M{"$all": []string{catalogLinkId}}
	}
	if dictId != "" {
		filter["dictId"] = dictId
	}
	cur, err := getWordTable().Find(context.TODO(), filter,
		options.Find().SetSort(bson.M{"writing": 1}).SetProjection(bson.M{
			"id": 1, "dictId": 1, "writing": 1, "tags": 1,
			"catalogLinkIds": 1, "completion": 1, "meaning": 1,
			"specs": bson.M{"$slice": 1},
		}))
	if err != nil {
		return
	}
	words = make([]*model.Word, 0, 10)
	for cur.Next(context.TODO()) {
		item := &model.Word{}
		err = cur.Decode(&item)
		if err != nil {
			return
		}
		words = append(words, item)
	}
	return
}

func InsertWord(word *model.Word) (err error) {
	_, err = getWordTable().InsertOne(context.TODO(), word)
	return
}

func UpdateWord(word *model.Word) (err error) {
	_, err = getWordTable().UpdateOne(context.TODO(), bson.M{
		"id": word.Id,
	}, bson.M{"$set": bson.M{
		"writing":        word.Writing,
		"specs":          word.Specs,
		"completion":     word.Completion,
		"catalogLinkIds": word.CatalogLinkIds,
		"meaning":        word.Meaning,
		"updateTime":     word.UpdateTime,
	}})
	return
}

func getWordTable() *mongo.Collection {
	table, err := db.GetTable(config.Get().DB.WordTable.Name)
	if err != nil {
		panic(err)
	}
	_, err = table.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"writing": 1},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.M{"tags": 1},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.M{"catalogLinkIds": 1},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.M{"completion": 1},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.M{"updateTime": -1},
			Options: options.Index().SetUnique(false),
		},
	})
	if err != nil {
		panic(err)
	}
	return table
}
