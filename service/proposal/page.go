package proposal

import (
	"app/infrastructure/log"
	"app/infrastructure/util"
	"app/service/dict/dao"
	"app/service/dict/model"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProposalVO struct {
	AddWord  *AddWordVO
	EditWord *EditWordVO
}

type AddWordVO struct {
	Dict        *model.DictBO
	Completions []*model.CompletionBO
}

type EditWordVO struct {
	Word        *model.WordBO
	Completions []*model.CompletionBO
}

func HandlePageApplyUpdateWord(context *gin.Context) {
	dictId := context.Query("dictId")
	wordId := context.Query("wordId")

	var word *model.WordBO
	var err error
	var dict *model.DictBO
	t := template.Must(template.New("proposal").Funcs(util.DictHubFuncMap()).ParseFiles(
		"static/proposal/proposal.gohtml",
		"static/proposal/word.add.gohtml",
		"static/proposal/word.edit.gohtml",
		"static/common/title.gohtml",
		"static/common/search.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))
	if wordId == "" {
		dict, err = dao.FindDictById(dictId)
		if err != nil {
			log.Error(err.Error())
			context.String(http.StatusNotFound, "dict not found")
			return
		}

		err = t.Execute(context.Writer, ProposalVO{
			AddWord: &AddWordVO{
				Dict:        dict,
				Completions: model.CompletionBOs,
			},
		})

		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	word, err = dao.FindWordById(wordId)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusNotFound, "word not found")
		return
	}

	err = t.Execute(context.Writer, ProposalVO{
		EditWord: &EditWordVO{
			Word:        word,
			Completions: model.CompletionBOs,
		},
	})
	if err != nil {
		log.Error(err.Error())
	}
}
