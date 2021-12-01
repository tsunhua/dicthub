package search

import (
	"app/infrastructure/cc"
	"app/infrastructure/config"
	"app/infrastructure/log"
	"app/service/dict/db"
	"app/service/dict/model"
	"strings"

	"github.com/expectedsh/go-sonic/sonic"
	"github.com/spf13/cast"
)

const (
	COL_DICTS         = "dicts"
	COL_WORDS         = "words"
	BUCKET            = "general"
	LANG              = sonic.LangAutoDetect
	PARALLEL_ROUTINES = 6
)

var sonicConfig config.Sonic

type Record struct {
	Id   string
	Text string
}

func RefreshWords(words []*model.Word) {
	objects := make([]string, len(words))
	for index, word := range words {
		objects[index] = word.Id
	}
	BulkFlushObject(COL_WORDS, objects)
	BulkPushWord(words)
}

func RefreshDicts(dicts []*model.Dict) {
	objects := make([]string, len(dicts))
	for index, dict := range dicts {
		objects[index] = dict.Id
	}
	BulkFlushObject(COL_DICTS, objects)
	BulkPushDict(dicts)
}

func BulkPushWord(words []*model.Word) {
	records := make([]*Record, len(words))
	for index, word := range words {
		// 同時存入簡繁體
		text := strings.Join(cc.Convert2All(word.Writing), " ")
		// 取屬性，並去除聲調符號
		if len(word.Specs) > 0 {
			for _, spec := range word.Specs {
				text = text + " " + clearTone(spec.Value)
			}
		}
		// 取含義的前20個字符
		if word.Meaning != "" {
			runes := []rune(word.Meaning)
			if len(runes) > 20 {
				runes = runes[:20]
			}
			text = text + " " + string(runes)
		}
		records[index] = &Record{
			Id:   word.Id,
			Text: text,
		}
	}
	BulkPush(COL_WORDS, records)
}

func BulkPushDict(dicts []*model.Dict) {
	records := make([]*Record, len(dicts))
	for index, dict := range dicts {
		// 同時存入簡繁體
		text := strings.Join(cc.Convert2All(dict.Name), " ")
		// 取辭典標籤
		if len(dict.Tags) > 0 {
			for _, tag := range dict.Tags {
				text = text + " " + tag
			}
		}
		records[index] = &Record{
			Id:   dict.Id,
			Text: text,
		}
	}
	BulkPush(COL_DICTS, records)
}

func BulkPush(col string, records []*Record) {
	ibrs := make([]sonic.IngestBulkRecord, len(records))
	for index, record := range records {
		ibrs[index] = sonic.IngestBulkRecord{
			Object: record.Id,
			Text:   record.Text,
		}
	}
	ingester, err := sonic.NewIngester(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		log.Error(err.Error())
		return
	}
	_ = ingester.BulkPush(col, BUCKET, PARALLEL_ROUTINES, ibrs, LANG)
}

func Query(col string, keyword string, limit int, offset int) (results []string, err error) {
	searcher, err := sonic.NewSearch(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return searcher.Query(col, BUCKET, keyword, limit, offset, LANG)
}

func BulkFlushObject(col string, objects []string) {
	ingester, err := sonic.NewIngester(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, object := range objects {
		_ = ingester.FlushObject(col, BUCKET, object)
	}
}

func FlushObject(col string, object string) {
	ingester, err := sonic.NewIngester(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		log.Error(err.Error())
		return
	}
	_ = ingester.FlushObject(col, BUCKET, object)
}

func init() {
	sonicConfig = config.Get().Sonic
	if sonicConfig.Repush {
		go repushDicts()
		go repushWords()
	}
}

func repushDicts() {
	const limit = 100
	stop := false
	offset := int64(0)
	log.Info("start repush dicts")
	for !stop {
		dicts, err := db.FindRecentDictList(limit, offset)
		if err != nil {
			log.Error(err.Error())
			return
		}
		BulkPushDict(dicts)
		log.Debug("bulk push dicts: " + cast.ToString(len(dicts)))

		if len(dicts) < limit {
			stop = true
		}
		offset = offset + limit
	}
	log.Info("repush dicts done!")
}

func repushWords() {
	const limit = 100
	stop := false
	offset := int64(0)
	log.Info("start repush words")
	for !stop {
		words, err := db.FindRecentWordList(limit, offset)
		if err != nil {
			log.Error(err.Error())
			return
		}
		BulkPushWord(words)
		log.Debug("bulk push words: " + cast.ToString(len(words)))

		if len(words) < limit {
			stop = true
		}
		offset = offset + limit
	}
	log.Info("repush words done!")
}
