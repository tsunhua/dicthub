package dao

import (
	"app/infrastructure/log"
	"app/infrastructure/util"
	"app/service/dict/db"
	"app/service/dict/model"
	"html/template"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

const (
	ID_LINKER     = "."
	NAME_LINKER   = "/"
	NUMBER_LINKER = "."
)

func FindDictById(id string) (dictBO *model.DictBO, err error) {
	var dict *model.Dict
	if dict, err = db.FindDictById(id); err != nil {
		return
	}
	dictBO, err = dictToDictBO(dict)
	return
}

func FindDictsBy(tags []string) (dictBOs []*model.DictBO, err error) {
	var dicts []*model.Dict
	dicts, err = db.FindDictsBy(tags)
	if err != nil {
		log.Error(err.Error())
		return
	}
	dictBOs, err = dictsToDictBOs(dicts)
	return
}

func FindManyDictsById(ids []string) (dictBOs []*model.DictBO, err error) {
	var dicts []*model.Dict
	dicts, err = db.FindManyDictsById(ids)
	if err != nil {
		log.Error(err.Error())
		return
	}
	dictBOs, err = dictsToDictBOs(dicts)
	return
}

func FindRecentDictList(limit int64) (dictBOs []*model.DictBO, err error) {
	var dicts []*model.Dict
	dicts, err = db.FindRecentDictList(limit, 0)
	if err != nil {
		log.Error(err.Error())
		return
	}
	dictBOs, err = dictsToDictBOs(dicts)
	return
}

func dictsToDictBOs(dicts []*model.Dict) (dictBOs []*model.DictBO, err error) {
	if len(dicts) == 0 {
		return []*model.DictBO{}, nil
	}
	dictBOs = make([]*model.DictBO, 0, len(dicts))

	for _, word := range dicts {
		var wordBO *model.DictBO
		wordBO, err = dictToDictBO(word)
		if err != nil {
			log.Error(err.Error())
			break
		}
		dictBOs = append(dictBOs, wordBO)
	}
	return
}

func dictToDictBO(dict *model.Dict) (dictBO *model.DictBO, err error) {
	if dict == nil {
		return
	}

	catalogTree := parse2TreeNodeBOs(dict.CatalogText)

	specTree := parse2TreeNodeBOs(dict.SpecText)

	var preferSpecs []*model.TreeNodeBO
	if dict.PreferSpecLinkIds != nil {
		preferSpecs = make([]*model.TreeNodeBO, len(dict.PreferSpecLinkIds))
		for i, it := range dict.PreferSpecLinkIds {
			for _, spec := range specTree {
				if spec.LinkId == it {
					preferSpecs[i] = spec
				}
			}
		}
	} else {
		preferSpecs = []*model.TreeNodeBO{}
	}

	dictBO = &model.DictBO{
		Id:            dict.Id,
		Name:          dict.Name,
		IsPublic:      dict.IsPublic,
		Cover:         dict.Cover,
		DescRaw:       dict.Desc,
		Desc:          template.HTML(util.MdToHtml([]byte(dict.Desc))),
		Contributor:   dict.Contributor,
		FeedbackEmail: dict.FeedbackEmail,
		CatalogTree:   catalogTree,
		SpecTree:      specTree,
		PreferSpecs:   preferSpecs,
		Tags:          dict.Tags,
		CreateTime:    dict.CreateTime,
		UpdateTime:    dict.UpdateTime,
	}
	return
}

func parse2TreeNodeBOs(text string) []*model.TreeNodeBO {
	lines := strings.Split(text, "\r\n")
	reg, err := regexp.Compile(`(#*) (.+?)/([^\/]*)/?(.*)?`)
	if err != nil {
		return nil
	}
	var arr = make([]*model.TreeNodeBO, 0, len(lines))
	lastLevel := 0
	lastLinkId := ""
	lastLinkName := ""
	lastNumber := ""
	count := 1
	for _, line := range lines {
		strs := reg.FindStringSubmatch(line)
		if len(strs) < 3 {
			continue
		}
		level := len(strs[1])
		id := strs[3]
		name := strs[2]
		var linkId, linkName, number string
		switch {
		case level == lastLevel:
			lastLinkId = lastLinkId[:strings.LastIndex(lastLinkId, ID_LINKER)]
			lastLinkName = lastLinkName[:strings.LastIndex(lastLinkName, NAME_LINKER)]
			count++
			lastNumber = lastNumber[:strings.LastIndex(lastNumber, NUMBER_LINKER)]
		case level < lastLevel:
			lastLinkId = lastLinkId[:strings.LastIndex(lastLinkId, ID_LINKER)]
			lastLinkId = lastLinkId[:strings.LastIndex(lastLinkId, ID_LINKER)]
			lastLinkName = lastLinkName[:strings.LastIndex(lastLinkName, NAME_LINKER)]
			lastLinkName = lastLinkName[:strings.LastIndex(lastLinkName, NAME_LINKER)]
			lastNumber = lastNumber[:strings.LastIndex(lastNumber, NUMBER_LINKER)]
			count = cast.ToInt(lastNumber[strings.LastIndex(lastNumber, NUMBER_LINKER)+1:]) + 1
			lastNumber = lastNumber[:strings.LastIndex(lastNumber, NUMBER_LINKER)]
		}
		linkId = lastLinkId + ID_LINKER + id
		linkName = lastLinkName + NAME_LINKER + name
		number = lastNumber + NUMBER_LINKER + cast.ToString(count)
		arr = append(arr, &model.TreeNodeBO{
			Level:    level,
			Name:     name,
			Id:       id,
			LinkId:   strings.TrimPrefix(linkId, ID_LINKER),
			LinkName: strings.TrimPrefix(linkName, NAME_LINKER),
			Number:   strings.TrimPrefix(number, NUMBER_LINKER),
		})
		lastLevel = level
		lastLinkId = linkId
		lastLinkName = linkName
		lastNumber = number
	}
	return arr
}
