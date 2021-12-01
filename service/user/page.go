package user

import (
	"app/infrastructure/log"
	"app/infrastructure/util"
	"html/template"

	"github.com/gin-gonic/gin"
)

func HandlePageLogin(context *gin.Context) {
	t := template.Must(template.New("login").Funcs(util.DictHubFuncMap()).ParseFiles(
		"static/user/login.gohtml",
		"static/common/title.gohtml",
		"static/common/search.gohtml",
		"static/common/head.gohtml",
		"static/common/footer.gohtml"))
	err := t.Execute(context.Writer, struct{}{})
	if err != nil {
		log.Error(err.Error())
	}
}
