package user

import (
	"app/infrastructure/config"
	"app/infrastructure/log"
	"app/service/user/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleAPILogin(context *gin.Context) {
	email := context.PostForm("email")
	u, err := db.FindUser(email)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusNotFound, "error: %s", err.Error())
		return
	}
	password := context.PostForm("password")
	if u.Password != password {
		log.Debug(fmt.Sprintf("password incorrect. expected: %s, actual: %s", u.Password, password))
		context.String(http.StatusUnauthorized, "email or password incorrect")
		return
	}
	// set cookie 1 month
	context.SetCookie("email", email, 2592000, "/", config.Get().Domain, false, true)
	context.SetCookie("password", password, 2592000, "/", config.Get().Domain, false, true)
	context.String(http.StatusOK, "ok")
}
