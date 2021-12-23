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
	ID_LINKER     = "~"
	NAME_LINKER   = "·"
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

	// catalogTree := traversal(dict.CatalogTree, &traversalContext{
	// 	lastLevel: -1,
	// })()
	catalogTree := parse2TreeNodeBOs(dict.CatalogText)

	// specTree := traversal(dict.SpecTree, &traversalContext{
	// 	lastLevel: -1,
	// })()
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
	lines := strings.Split(text, "\n")
	reg, err := regexp.Compile("(#*) (.*)/(.*)")
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
		if len(strs) <3{
			continue
		}
		level := len(strs[1])
		id := strs[3]
		name := strs[2]
		var linkId, linkName, number string
		switch {
		case level == lastLevel:
			lastLinkId = lastLinkId[:strings.LastIndex(lastLinkId, "~")]
			lastLinkName = lastLinkName[:strings.LastIndex(lastLinkName, "~")]
			count++
			lastNumber = lastNumber[:strings.LastIndex(lastNumber, ".")]
		case level < lastLevel:
			lastLinkId = lastLinkId[:strings.LastIndex(lastLinkId, "~")]
			lastLinkId = lastLinkId[:strings.LastIndex(lastLinkId, "~")]
			lastLinkName = lastLinkName[:strings.LastIndex(lastLinkName, "~")]
			lastLinkName = lastLinkName[:strings.LastIndex(lastLinkName, "~")]
			lastNumber = lastNumber[:strings.LastIndex(lastNumber, ".")]
			count =  cast.ToInt(lastNumber[strings.LastIndex(lastNumber, ".")+1:])+1
			lastNumber = lastNumber[:strings.LastIndex(lastNumber, ".")]
		}
		linkId = lastLinkId + "~" + id
		linkName = lastLinkName + "~" + name
		number = lastNumber + "." + cast.ToString(count)
		arr = append(arr, &model.TreeNodeBO{
			Level:    level,
			Name:     name,
			Id:       id,
			LinkId:   strings.TrimPrefix(linkId, "~"),
			LinkName: strings.TrimPrefix(linkName, "~"),
			Number:   strings.TrimPrefix(number, "."),
		})
		lastLevel = level
		lastLinkId = linkId
		lastLinkName = linkName
		lastNumber = number
	}
	return arr
}

// type traversalContext struct {
// 	lastLevel    int
// 	lastNumber   string
// 	lastLinkId   string
// 	lastLinkName string
// }

// func traversal(nodes []*model.TreeNode, ctx *traversalContext) func() []*model.TreeNodeBO {
// 	var arr = make([]*model.TreeNodeBO, 0, 10)
// 	return func() []*model.TreeNodeBO {
// 		level := ctx.lastLevel + 1
// 		number := ""
// 		for index, node := range nodes {
// 			linkId := ctx.lastLinkId
// 			if linkId == "" {
// 				linkId = node.Id
// 			} else {
// 				linkId = linkId + ID_LINKER + node.Id
// 			}

// 			linkName := ctx.lastLinkName
// 			if linkName == "" {
// 				linkName = node.Name
// 			} else {
// 				linkName = linkName + NAME_LINKER + node.Name
// 			}

// 			if ctx.lastNumber == "" {
// 				number = cast.ToString(index + 1)
// 			} else {
// 				number = ctx.lastNumber + NUMBER_LINKER + cast.ToString(index+1)
// 			}

// 			nodeBO := &model.TreeNodeBO{
// 				Id:          node.Id,
// 				Name:        node.Name,
// 				LinkId:      linkId,
// 				LinkName:    linkName,
// 				Level:       level,
// 				IsLastLevel: node.Next == nil,
// 				Number:      number,
// 			}

// 			arr = append(arr, nodeBO)
// 			if node.Next != nil {
// 				arr = append(arr, traversal(node.Next, &traversalContext{
// 					lastLinkId:   linkId,
// 					lastLinkName: linkName,
// 					lastLevel:    level,
// 					lastNumber:   number,
// 				})()...)
// 			}
// 		}
// 		return arr
// 	}
// }
