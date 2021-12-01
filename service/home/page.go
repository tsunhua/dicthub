package home

import (
	"app/infrastructure/log"
	"app/infrastructure/util"
	"app/service/dict/dao"
	"app/service/dict/model"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"sync"
)

func HandlePagePing(context *gin.Context) {
	context.String(http.StatusOK, "pong")
}

func HandlePageIndex(context *gin.Context) {
	t := template.Must(template.New("index").Funcs(util.DictHubFuncMap()).ParseFiles("static/index.gohtml",
		"static/common/dict.list.gohtml",
		"static/common/title.gohtml",
		"static/common/search.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))
	var items []*model.WordBO
	var dicts []*model.DictBO
	var err error
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		items, err = dao.FindRecentWordList(12)
		if err != nil {
			log.Error(err.Error())
			items = []*model.WordBO{}
		}
		wg.Done()
	}()

	go func() {
		dicts, err = dao.FindRecentDictList(12)
		if err != nil {
			log.Error(err.Error())
			dicts = []*model.DictBO{}
		}
		wg.Done()
	}()

	wg.Wait()

	err = t.Execute(context.Writer, struct {
		Items []*model.WordBO
		Dicts []*model.DictBO
	}{
		Items: items,
		Dicts: dicts,
	})
	if err != nil {
		log.Error(err.Error())
	}
}
