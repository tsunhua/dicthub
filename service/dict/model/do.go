package model

import (
	"time"
)

type Completion string

const (
	Draft   Completion = "draft"
	Release Completion = "release"
	Perfect Completion = "perfect"
)

type Dict struct {
	Id                string      `bson:"id" json:"id"`
	Name              string      `bson:"name" json:"name"`
	IsPublic          bool        `bson:"isPublic" json:"isPublic"`
	Desc              string      `bson:"desc" json:"desc"`
	Cover             string      `bson:"cover" json:"cover"`
	FeedbackEmail     string      `bson:"feedbackEmail" json:"feedbackEmail"`
	CatalogTree       []*TreeNode `bson:"catalogTree,omitempty" json:"catalogTree,omitempty"`
	CatalogText       string      `bson:"catalogText" json:"catalogText"`
	SpecTree          []*TreeNode `bson:"specTree,omitempty" json:"specTree,omitempty"`
	SpecText          string      `bson:"specText" json:"specText"`
	PreferSpecLinkIds []string    `bson:"preferSpecLinkIds" json:"preferSpecLinkIds"`
	Tags              []string    `bson:"tags" json:"tags"`
	CreateTime        time.Time   `bson:"createTime,omitempty" json:"createTime"`
	UpdateTime        time.Time   `bson:"updateTime,omitempty" json:"updateTime"`
	Creator           string      `bson:"creater,omitempty" json:"creater"`   // 創建者郵箱
	Updators          []string    `bson:"updators,omitempty" json:"updators"` // 更新者郵箱列表
}

type TreeNode struct {
	Id   string      `bson:"id" json:"id"`
	Name string      `bson:"name" json:"name"`
	Next []*TreeNode `bson:"next" json:"next"`
}

type Word struct {
	Id             string     `bson:"id" json:"id"`
	DictId         string     `bson:"dictId,omitempty" json:"dictId"`
	Writing        string     `bson:"writing" json:"writing"`
	CatalogLinkIds []string   `bson:"catalogLinkIds" json:"catalogLinkIds"`
	Meaning        string     `bson:"meaning" json:"meaning"`
	Specs          []*Spec    `bson:"specs" json:"specs"`
	Completion     Completion `bson:"completion" json:"completion"`
	SourceUrl      string     `bson:"sourceUrl,omitempty" json:"sourceUrl"`
	CreateTime     time.Time  `bson:"createTime,omitempty" json:"createTime"`
	UpdateTime     time.Time  `bson:"updateTime,omitempty" json:"updateTime"`
	Creator        string     `bson:"creater,omitempty" json:"creater"`   // 創建者郵箱
	Updators       []string   `bson:"updators,omitempty" json:"updators"` // 更新者郵箱列表
}

// 這個不只可以表示發音，理論上可以表示任意資料項，就像商品的規格信息一樣。
type Spec struct {
	LinkId string `bson:"linkId" json:"linkId"`
	Value  string `bson:"value" json:"value"`
	Note   string `bson:"note,omitempty" json:"note,omitempty"`
}
