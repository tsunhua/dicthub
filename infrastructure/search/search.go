package search

import (
	"app/infrastructure/cc"
	"app/infrastructure/config"
	"app/infrastructure/log"
	"app/service/dict/db"
	"app/service/dict/model"
	"fmt"
	"strings"

	"github.com/expectedsh/go-sonic/sonic"
	"github.com/spf13/cast"
)

const (
	COL_DICTS         = "dicts"
	COL_WORDS         = "words"
	BUCKET_GENERAL    = "general"
	BUCKET_PLUS       = "plus"
	LANG              = sonic.LangCmn // Mandarin
	PARALLEL_ROUTINES = 6
)

var sonicConfig config.Sonic

func init() {
	sonicConfig = config.Get().Sonic
}

type Record struct {
	Id   string
	Text string
}

func RefreshWords(words []*model.Word) {
	objects := make([]string, len(words))
	for index, word := range words {
		objects[index] = word.Id
	}
	BulkFlushObject(COL_WORDS, BUCKET_GENERAL, objects)
	BulkFlushObject(COL_WORDS, BUCKET_PLUS, objects)
	BulkPushWord(words)
}

func RefreshDicts(dicts []*model.Dict) {
	objects := make([]string, len(dicts))
	for index, dict := range dicts {
		objects[index] = dict.Id
	}
	BulkFlushObject(COL_DICTS, BUCKET_GENERAL, objects)
	BulkFlushObject(COL_DICTS, BUCKET_PLUS, objects)
	BulkPushDict(dicts)
}

func BulkPushWord(words []*model.Word) {
	records := make([]*Record, len(words))
	plusRecords := make([]*Record, len(words))
	for index, word := range words {
		// 同時存入簡繁體
		text := strings.Join(cc.Convert2All(word.Writing), " ")
		// 取屬性，並去除聲調符號
		if len(word.Specs) > 0 {
			for _, spec := range word.Specs {
				text = text + " " + clearTone(spec.Value)
			}
		}
		records[index] = &Record{
			Id:   word.Id,
			Text: text,
		}

		if word.Meaning != "" {
			runes := []rune(word.Meaning)
			if len(runes) > 20 {
				runes = runes[:20] // 取含義的前20個字符
			}
			text = strings.Join(cc.Convert2All(string(runes)), " ")
			plusRecords[index] = &Record{
				Id:   word.Id,
				Text: text,
			}
		}

	}
	BulkPush(COL_WORDS, BUCKET_GENERAL, records)
	BulkPush(COL_WORDS, BUCKET_PLUS, plusRecords)
}

func BulkPushDict(dicts []*model.Dict) {
	records := make([]*Record, len(dicts))
	plusRecords := make([]*Record, len(dicts))
	for index, dict := range dicts {
		// 同時存入簡繁體
		text := strings.Join(cc.Convert2All(dict.Name), " ")
		records[index] = &Record{
			Id:   dict.Id,
			Text: text,
		}

		text = ""
		// 取辭典標籤
		if len(dict.Tags) > 0 {
			for _, tag := range dict.Tags {
				text = text + " " + strings.Join(cc.Convert2All(tag), " ")
			}
		}
		if dict.Desc != "" {
			text = text + " " + strings.Join(cc.Convert2All(dict.Desc), " ")
		}
		if strings.TrimSpace(text) != "" {
			plusRecords[index] = &Record{
				Id:   dict.Id,
				Text: text,
			}
		}
	}
	BulkPush(COL_DICTS, BUCKET_GENERAL, records)
	BulkPush(COL_DICTS, BUCKET_PLUS, plusRecords)
}

func BulkPush(col string, bucket string, records []*Record) {
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
	_ = ingester.BulkPush(col, bucket, PARALLEL_ROUTINES, ibrs, LANG)
}

func Query(col string, bucket string, keyword string, limit int, offset int) (results []string, err error) {
	searcher, err := sonic.NewSearch(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return searcher.Query(col, bucket, keyword, limit, offset, LANG)
}

func BulkFlushObject(col string, bucket string, objects []string) {
	ingester, err := sonic.NewIngester(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, object := range objects {
		_ = ingester.FlushObject(col, bucket, object)
	}
}

func FlushObject(col string, bucket string, object string) {
	ingester, err := sonic.NewIngester(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		log.Error(err.Error())
		return
	}
	_ = ingester.FlushObject(col, bucket, object)
}

func FlushCollection(col string) (err error) {
	ingester, err := sonic.NewIngester(sonicConfig.Host, cast.ToInt(sonicConfig.Port), sonicConfig.Pwd)
	if err != nil {
		return
	}
	err = ingester.FlushCollection(COL_DICTS)
	if err != nil {
		return
	}
	return
}

func Sync() (err error) {
	if err = FlushCollection(COL_DICTS); err != nil {
		return
	}
	syncDicts()
	if err = FlushCollection(COL_WORDS); err != nil {
		return
	}
	syncWords()
	return
}

func syncDicts() {
	const limit = 100
	stop := false
	offset := int64(0)
	log.Info("start sync dicts")
	for !stop {
		dicts, err := db.FindRecentDictList(limit, offset)
		if err != nil {
			log.Error(err.Error())
			return
		}
		BulkPushDict(dicts)
		log.Info(fmt.Sprintf("bulk pushed dicts: %d", int64(len(dicts))+offset))

		if len(dicts) < limit {
			stop = true
		}
		offset = offset + limit
	}
	log.Info("sync dicts done!")
}

func syncWords() {
	const limit = 100
	stop := false
	offset := int64(0)
	log.Info("start sync words")
	for !stop {
		words, err := db.FindRecentWordList(limit, offset)
		if err != nil {
			log.Error(err.Error())
			return
		}
		BulkPushWord(words)
		log.Info(fmt.Sprintf("bulk pushed words: %d", int64(len(words))+offset))

		if len(words) < limit {
			stop = true
		}
		offset = offset + limit
	}
	log.Info("sync words done!")
}
