package db

import (
	"app/infrastructure/config"
	"app/infrastructure/db"
	"app/service/proposal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindProposalById(id string) (proposal *model.Proposal, err error) {
	err = getProposalTable().FindOne(context.TODO(), bson.M{"id": id}).Decode(&proposal)
	return
}

func FindProposalByStatus(status string) (proposal *model.Proposal, err error) {
	err = getProposalTable().FindOne(context.TODO(), bson.M{"status": status}).Decode(&proposal)
	return
}

func InsertProposal(proposal *model.Proposal) (err error) {
	_, err = getProposalTable().InsertOne(context.TODO(), proposal)
	return
}

func UpdateProposal(proposal *model.Proposal) (err error) {
	update := bson.M{
		"updateTime": proposal.UpdateTime,
	}
	if proposal.Status != "" {
		update["status"] = proposal.Status
	}
	_, err = getProposalTable().UpdateOne(context.TODO(), bson.M{"id": proposal.Id}, bson.M{"$set": update})
	return
}

func getProposalTable() *mongo.Collection {
	table, err := db.GetTable(config.Get().DB.ProposalTable.Name)
	if err != nil {
		panic(err)
	}
	_, err = table.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"status": 1},
			Options: options.Index().SetUnique(false),
		},
	})
	if err != nil {
		panic(err)
	}
	return table
}
