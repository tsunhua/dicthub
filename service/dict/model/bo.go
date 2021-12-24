package model

import (
	"html/template"
	"time"
)

type WordBO struct {
	Id             string        `bson:"id" json:"id"`
	Dict           *DictBO       `bson:"dict" json:"dict"`
	CatalogLinkIds []string      `bson:"catalogLinkIds" json:"catalogLinkIds"`
	Catalogs       []*CatalogBO  `bson:"catalogs" json:"catalogs"`
	Location       []string      `bson:"location" json:"location"`
	Writing        string        `bson:"writing" json:"writing"`
	Meaning        template.HTML `bson:"meaning" json:"meaning"`
	MeaningRaw     string        `bson:"meaningRaw" json:"meaningRaw"`
	Specs          []*SpecBO     `bson:"specs" json:"specs"`
	SourceUrl      string        `bson:"sourceUrl" json:"sourceUrl"`
	Completion     *CompletionBO `bson:"completion" json:"completion"`
	CreateTime     time.Time     `bson:"createTime" json:"createTime"`
	UpdateTime     time.Time     `bson:"updateTime" json:"updateTime"`
}

type DictBO struct {
	Id            string        `bson:"id" json:"id"`
	Name          string        `bson:"name" json:"name"`
	IsPublic      bool          `bson:"isPublic" json:"isPublic"`
	Cover         string        `bson:"cover" json:"cover"`
	DescRaw       string        `bson:"desc" json:"descRaw"`
	Desc          template.HTML `bson:"desc" json:"desc"`
	Contributor   string        `bson:"contributor" json:"contributor"`
	FeedbackEmail string        `bson:"feedbackEmail" json:"feedbackEmail"`
	CatalogTree   []*TreeNodeBO `bson:"catalogTree" json:"catalogTree"`
	SpecTree      []*TreeNodeBO `bson:"specTree" json:"specTree"`
	PreferSpecs   []*TreeNodeBO `bson:"preferSpecs" json:"preferSpecs"`
	Tags          []string      `bson:"tags" json:"tags"`
	CreateTime    time.Time     `bson:"createTime" json:"createTime"`
	UpdateTime    time.Time     `bson:"updateTime" json:"updateTime"`
}

// 場景： 編輯條目（扁平結構即可）、查看辭書（深度優先遍歷）
type TreeNodeBO struct {
	Id          string `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	LinkId      string `bson:"linkId" json:"linkId"`     // 拼接的ID
	LinkName    string `bson:"linkName" json:"linkName"` // 拼接的Name
	Number      string `bson:"number" json:"number"`     // 编号
	Level       int    `bson:"level" json:"level"`
	IsLastLevel bool   `bson:"isLastLevel" json:"isLastLevel"` // 是否是最後一級
}

type CompletionBO struct {
	Name  string
	Value Completion
}

var CompletionBOMap = map[Completion]string{
	Draft:   "草稿",
	Release: "可發佈",
	Perfect: "完美",
}

var CompletionBOs = []*CompletionBO{
	{Name: CompletionBOMap[Draft], Value: Draft},
	{Name: CompletionBOMap[Release], Value: Release},
	{Name: CompletionBOMap[Perfect], Value: Perfect},
}

type SpecBO struct {
	LinkId   string `bson:"linkId" json:"linkId"`
	LinkName string `bson:"linkName" json:"linkName"`
	Value    string `bson:"value" json:"value"`
	Note     string `bson:"note" json:"note"`
}

type CatalogBO struct {
	Id       string `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	LinkId   string `bson:"linkId" json:"linkId"`
	LinkName string `bson:"linkName" json:"linkName"`
}
