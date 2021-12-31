package model

import (
	dictModel "app/service/dict/model"
	"time"
)

type ProposalKind string

const (
	KindAddDict  ProposalKind = "add_dict"
	KindEditDict ProposalKind = "edit_dict"
	KindAddWord  ProposalKind = "add_word"
	KindEditWord ProposalKind = "edit_word"
)

type ProposalStatus string

const (
	StatusChecking       ProposalStatus = "checking"
	StatusApproved       ProposalStatus = "approved"
	StatusPartlyApproved ProposalStatus = "partly_approved"
	StatusRejected       ProposalStatus = "rejected"
)

type Proposal struct {
	Id         string          `bson:"id" json:"id"`
	Kind       ProposalKind    `bson:"kind" json:"kind"`
	Word       *dictModel.Word `bson:"word,omitempty" json:"word"`
	Dict       *dictModel.Dict `bson:"dict,omitempty" json:"dict"`
	Status     ProposalStatus  `bson:"status" json:"status"`
	Applicant  *Applicant      `bson:"applicant" json:"applicant"`
	CreateTime time.Time       `bson:"createTime" json:"createTime"`
	UpdateTime time.Time       `bson:"updateTime" json:"updateTime"`
}

type Applicant struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
}
