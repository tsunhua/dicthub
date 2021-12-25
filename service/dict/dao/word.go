package dao

import (
	"app/infrastructure/cache"
	"app/infrastructure/log"
	"app/infrastructure/util"
	"app/service/dict/db"
	"app/service/dict/model"
	"html/template"
	"time"
)

func FindRecentWordList(limit int64) (wordBOs []*model.WordBO, err error) {
	var words []*model.Word
	words, err = db.FindRecentWordList(limit, 0)

	wordBOs = make([]*model.WordBO, 0, len(words))
	for _, word := range words {
		var wordBO = model.WordBO{
			Id:      word.Id,
			Writing: word.Writing,
			Completion: &model.CompletionBO{
				Name:  model.CompletionBOMap[word.Completion],
				Value: word.Completion,
			},
			UpdateTime: word.UpdateTime,
		}
		wordBOs = append(wordBOs, &wordBO)
	}
	return
}

func FindWordById(id string) (wordBO *model.WordBO, err error) {
	ckey := "/word/" + id
	var word *model.Word
	cword, err := cache.Cache().Get(ckey)
	if err == nil {
		wordBO = cword.(*model.WordBO)
		return
	}
	log.Debug(err.Error())

	word, err = db.FindWordById(id)
	if err != nil {
		log.Error(err.Error())
		return
	}
	wordBO, err = wordToWordBO(word)
	if err != nil {
		return
	}
	err = cache.Cache().SetWithExpire(ckey, wordBO, 30*time.Second)
	if err != nil {
		log.Warn(err.Error())
		err = nil
	}
	return
}

func FindManyWordsById(ids []string) (wordBOs []*model.WordBO, err error) {
	var words []*model.Word
	words, err = db.FindManyWordsById(ids)
	if err != nil {
		log.Error(err.Error())
		return
	}
	wordBOs, err = wordsToWordBOs(words)
	return
}

func FindWordsBy(catalogLinkId, dictId string) (wordBOs []*model.WordBO, err error) {
	var words []*model.Word
	words, err = db.FindWordsBy(catalogLinkId, dictId)
	if err != nil {
		log.Error(err.Error())
		return
	}
	wordBOs, err = wordsToWordBOs(words)
	return
}

func FindBriefWordsBy(catalogLinkId, dictId string) (wordBOs []*model.WordBO, err error) {
	var words []*model.Word
	words, err = db.FindBriefWordsBy(catalogLinkId, dictId)
	if err != nil {
		log.Error(err.Error())
		return
	}
	wordBOs, err = wordsToWordBOs(words)
	return
}

func FindWordByWriting(writing string) (wordBOs []*model.WordBO, err error) {
	var words []*model.Word
	words, err = db.FindWordByWriting(writing)
	if err != nil {
		log.Error(err.Error())
		return
	}
	wordBOs, err = wordsToWordBOs(words)
	return
}

func wordsToWordBOs(words []*model.Word) (wordBOs []*model.WordBO, err error) {
	if len(words) == 0 {
		return []*model.WordBO{}, nil
	}
	wordBOs = make([]*model.WordBO, 0, len(words))

	for _, word := range words {
		var wordBO *model.WordBO
		wordBO, err = wordToWordBO(word)
		if err != nil {
			log.Error(err.Error())
			break
		}
		wordBOs = append(wordBOs, wordBO)
	}
	return
}

func wordToWordBO(word *model.Word) (wordBO *model.WordBO, err error) {
	if word == nil {
		return
	}

	dictBO, err := FindDictById(word.DictId)

	if err != nil || dictBO == nil {
		return
	}

	var pbos = make([]*model.SpecBO, 0, len(word.Specs))
	if len(word.Specs) > 0 {
		for _, p := range word.Specs {
			linkName := ""
			for _, node := range dictBO.SpecTree {
				if node.LinkId == p.LinkId {
					linkName = node.LinkName
					break
				}
			}

			pbos = append(pbos, &model.SpecBO{
				LinkId:   p.LinkId,
				LinkName: linkName,
				Value:    p.Value,
				Note:     p.Note,
			})
		}
	}
	var catalogs = make([]*model.CatalogBO, 0, len(word.CatalogLinkIds))
	for _, cid := range word.CatalogLinkIds {
		var catalogBO *model.CatalogBO
		for _, dcnode := range dictBO.CatalogTree {
			if dcnode.LinkId == cid {
				catalogBO = &model.CatalogBO{
					Id:       dcnode.Id,
					LinkId:   dcnode.LinkId,
					Name:     dcnode.Name,
					LinkName: dcnode.LinkName,
				}
			}
		}
		catalogs = append(catalogs, catalogBO)
	}

	// log.Debug("", log.Reflect("catalogs", catalogs))

	return &model.WordBO{
		Id:             word.Id,
		Dict:           dictBO,
		Writing:        word.Writing,
		Meaning:        template.HTML(util.MdToHtml([]byte(word.Meaning))),
		MeaningRaw:     word.Meaning,
		CatalogLinkIds: word.CatalogLinkIds,
		Catalogs:       catalogs,
		Specs:          pbos,
		SourceUrl:      word.SourceUrl,
		Completion: &model.CompletionBO{
			Name:  model.CompletionBOMap[word.Completion],
			Value: word.Completion,
		},
		CreateTime: word.CreateTime,
		UpdateTime: word.UpdateTime,
	}, nil
}
