package dict

import (
	"app/infrastructure/log"
	"app/infrastructure/search"
	"app/infrastructure/util"
	"app/service/dict/dao"
	"app/service/dict/db"
	"app/service/dict/model"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DICT pages

func HandlePageQueryDict(context *gin.Context) {
	tags := context.Query("tags")

	if tags == "" {
		context.String(http.StatusBadRequest, "tags required")
		return
	}
	t := template.Must(template.New("dict.query").Funcs(util.DictHubFuncMap()).ParseFiles(
		"static/dict/dict.query.gohtml",
		"static/common/dict.list.gohtml",
		"static/common/title.gohtml",
		"static/common/search.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))

	var dicts []*model.DictBO
	var err error
	tagArr := strings.Split(tags, ",")
	dicts, err = dao.FindDictsBy(tagArr)

	if dicts == nil || err != nil {
		context.String(http.StatusNotFound, "dict not found")
		return
	}

	err = t.Execute(context.Writer, struct {
		Dicts []*model.DictBO
		Tags  []string
		Kw    string
	}{
		Dicts: dicts,
		Tags:  tagArr,
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func HandlePageDict(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		context.String(http.StatusBadRequest, "id required")
		return
	}
	categoryLinkId := context.Param("categoryLinkId")
	t := template.Must(template.New("dict").Funcs(util.DictHubFuncMap()).ParseFiles("static/dict/dict.gohtml",
		"static/common/title.gohtml",
		"static/common/search.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))

	var err error
	var dict *model.DictBO

	dict, err = dao.FindDictById(id)

	if dict == nil || err != nil {
		context.String(http.StatusNotFound, "dict not found")
		return
	}

	var words []*model.WordBO

	if categoryLinkId != "" {
		if words, err = dao.FindBriefWordsBy(categoryLinkId, id); err != nil {
			log.Error(err.Error())
			context.String(http.StatusNotFound, "find words error")
			return
		}
	}

	err = t.Execute(context.Writer, struct {
		CategoryLinkId string
		Dict           *model.DictBO
		Words          []*model.WordBO
	}{
		CategoryLinkId: categoryLinkId,
		Dict:           dict,
		Words:          words,
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func HandlePageEditDict(context *gin.Context) {
	dictId := context.Query("dictId")
	ok, _ := hasPermission(context)

	var err error
	var dict *model.Dict

	t := template.Must(template.New("dict.edit").Funcs(util.DictHubFuncMap()).ParseFiles("static/dict/dict.edit.gohtml",
		"static/common/title.gohtml",
		"static/common/head.gohtml",
		"static/common/search.gohtml",
		"static/common/footer.gohtml"))

	if ok {
		dict, err = db.FindDictById(dictId)
		if err != nil {
			log.Error(err.Error())
			context.String(http.StatusNotFound, "dict not found")
			return
		}
	}

	err = t.Execute(context.Writer, struct {
		HasPermission bool
		Dict          *model.Dict
	}{
		ok,
		dict,
	})
	if err != nil {
		log.Error(err.Error())
	}
}

// WORD pages

func HandlePageQueryWord(context *gin.Context) {
	dictId := context.Query("dictId")
	catalogLinkId := context.Query("catalogLinkId")

	if dictId == "" && catalogLinkId == "" {
		context.String(http.StatusBadRequest, "dictId and catalogLinkId required")
		return
	}
	t := template.Must(template.New("word.query").Funcs(util.DictHubFuncMap()).ParseFiles("static/dict/word.query.gohtml",
		"static/common/search.gohtml",
		"static/common/title.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))

	var words []*model.WordBO
	var err error
	words, err = dao.FindWordsBy(catalogLinkId, dictId)

	if words == nil || err != nil {
		context.String(http.StatusNotFound, "word not found")
		return
	}

	err = t.Execute(context.Writer, struct {
		Words []*model.WordBO
	}{
		Words: words,
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func HandlePageWord(context *gin.Context) {
	writing := context.Param("writing")
	idTrunc := context.Param("id")
	log.Debug("writing: " + writing + " id: " + idTrunc)
	if writing == "" {
		context.String(http.StatusBadRequest, "writing required")
		return
	}
	t := template.Must(template.New("word").Funcs(util.DictHubFuncMap()).ParseFiles("static/dict/word.gohtml",
		"static/common/title.gohtml",
		"static/common/search.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))

	words, err := dao.FindWordByWriting(writing)

	if words == nil || err != nil {
		context.String(http.StatusNotFound, "word not found")
		return
	}

	if idTrunc != "" {
		var word *model.WordBO
		for _, item := range words {
			if strings.Contains(item.Id, idTrunc) {
				word = item
			}
		}
		if word != nil {
			words = []*model.WordBO{word}
		} else {
			words = []*model.WordBO{}
		}
	}

	err = t.Execute(context.Writer, struct {
		Words []*model.WordBO
	}{
		Words: words,
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func HandlePageEditWord(context *gin.Context) {
	dictId := context.Query("dictId")
	wordId := context.Query("wordId")
	ok, _ := hasPermission(context)

	var word *model.WordBO
	var err error
	var dict *model.DictBO

	if wordId == "" {
		t := template.Must(template.New("word.add").Funcs(util.DictHubFuncMap()).ParseFiles("static/dict/word.add.gohtml",
			"static/common/title.gohtml",
			"static/common/search.gohtml",
			"static/common/head.gohtml",
			"static/common/footer.gohtml"))
		if ok {
			dict, err = dao.FindDictById(dictId)
			if err != nil {
				log.Error(err.Error())
				context.String(http.StatusNotFound, "dict not found")
				return
			}
		}

		err = t.Execute(context.Writer, struct {
			Dict          *model.DictBO
			HasPermission bool
			Completions   []*model.CompletionBO
		}{
			dict, ok, model.CompletionBOs,
		})
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	t := template.Must(template.New("word.edit").Funcs(util.DictHubFuncMap()).ParseFiles("static/dict/word.edit.gohtml", "static/common/title.gohtml",
		"static/common/head.gohtml",
		"static/common/search.gohtml",
		"static/common/footer.gohtml"))

	if ok {
		word, err = dao.FindWordById(wordId)
		if err != nil {
			log.Error(err.Error())
			context.String(http.StatusNotFound, "word not found")
			return
		}
	}

	err = t.Execute(context.Writer, struct {
		HasPermission bool
		Completions   []*model.CompletionBO
		Word          *model.WordBO
	}{
		ok,
		model.CompletionBOs,
		word,
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func HandlePageSearchWords(context *gin.Context) {
	kw := context.Query("kw")

	if kw == "" {
		context.String(http.StatusBadRequest, "kw required")
		return
	}
	t := template.Must(template.New("word.query").Funcs(util.DictHubFuncMap()).ParseFiles("static/dict/word.query.gohtml",
		"static/common/search.gohtml",
		"static/common/title.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))

	var results []string
	var words []*model.WordBO
	var err error
	results, err = search.Query(search.COL_WORDS, kw, 20, 0)
	if err != nil || len(results) == 0 {
		log.Error(err.Error())
		context.String(http.StatusNotFound, "words not found")
		return
	}

	words, err = dao.FindManyWordsById(results)

	if words == nil || err != nil {
		log.Error(err.Error())
		context.String(http.StatusNotFound, "words not found")
		return
	}

	err = t.Execute(context.Writer, struct {
		Words []*model.WordBO
		Kw    string
	}{
		Words: words,
		Kw:    kw,
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func HandlePageSearchDicts(context *gin.Context) {
	kw := context.Query("kw")

	if kw == "" {
		context.String(http.StatusBadRequest, "kw required")
		return
	}
	t := template.Must(template.New("dict.query").Funcs(util.DictHubFuncMap()).ParseFiles(
		"static/dict/dict.query.gohtml",
		"static/common/dict.list.gohtml",
		"static/common/search.gohtml",
		"static/common/title.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))

	var results []string
	var dicts []*model.DictBO
	var err error
	results, err = search.Query(search.COL_DICTS, kw, 20, 0)
	if err != nil || len(results) == 0 {
		log.Error(err.Error())
		context.String(http.StatusNotFound, "dicts not found")
		return
	}

	dicts, err = dao.FindManyDictsById(results)

	if dicts == nil || err != nil {
		log.Error(err.Error())
		context.String(http.StatusNotFound, "dicts not found")
		return
	}

	err = t.Execute(context.Writer, struct {
		Tags  []string
		Kw    string
		Dicts []*model.DictBO
	}{
		Kw:    kw,
		Dicts: dicts,
	})
	if err != nil {
		log.Error(err.Error())
	}
}
