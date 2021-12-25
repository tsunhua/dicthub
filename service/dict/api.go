package dict

import (
	"app/infrastructure/cache"
	"app/infrastructure/log"
	"app/infrastructure/search"
	"app/infrastructure/util"
	"app/service/dict/db"
	"app/service/dict/model"
	userService "app/service/user"
	userModel "app/service/user/model"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func HandleAPIGetWord(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		log.Error("id required")
		context.String(http.StatusBadRequest, "param id required")
		return
	}
	word, err := db.FindWordById(id)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusBadRequest, "parse request body fail")
		return
	}
	context.JSON(http.StatusOK, word)
}

func HandleAPICreateWord(context *gin.Context) {
	ok, _ := hasPermission(context)
	if !ok {
		context.String(http.StatusUnauthorized, "permission required, login first")
		return
	}
	var word model.Word
	err := json.NewDecoder(context.Request.Body).Decode(&word)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusBadRequest, "parse request body fail")
		return
	}
	word.Id = uuid.NewV4().String()
	word.CreateTime = util.GetCurrentShanghaiTime()
	word.UpdateTime = word.CreateTime
	err = db.InsertWord(&word)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusInternalServerError, "insert word fail")
		return
	}
	// 更新搜索引擎數據
	go search.RefreshWords([]*model.Word{&word})
	context.String(http.StatusOK, "add word success")
}

func HandleAPIUpdateWord(context *gin.Context) {
	ok, _ := hasPermission(context)
	if !ok {
		context.String(http.StatusUnauthorized, "permission required, login first")
		return
	}
	var word model.Word
	err := json.NewDecoder(context.Request.Body).Decode(&word)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusBadRequest, "parse request body fail")
		return
	}
	word.UpdateTime = util.GetCurrentShanghaiTime()
	err = db.UpdateWord(&word)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusInternalServerError, "update word fail")
		return
	}
	// 更新搜索引擎數據
	go search.RefreshWords([]*model.Word{&word})
	// 移除舊的緩存
	go func() {
		cache.Cache().Remove("/word/" + word.Id)
	}()
	context.String(http.StatusOK, "update word success")
}

func HandleAPIGetDict(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		log.Error("id required")
		context.String(http.StatusBadRequest, "param id required")
		return
	}
	dict, err := db.FindDictById(id)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusBadRequest, "parse request body fail")
		return
	}
	context.JSON(http.StatusOK, dict)
}

func HandleAPICreateDict(context *gin.Context) {
	ok, _ := hasPermission(context)
	if !ok {
		context.String(http.StatusUnauthorized, "permission required, login first")
		return
	}
	var dict model.Dict
	err := json.NewDecoder(context.Request.Body).Decode(&dict)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusBadRequest, "parse request body fail")
		return
	}
	dict.Id = uuid.NewV4().String()
	dict.CreateTime = util.GetCurrentShanghaiTime()
	dict.UpdateTime = dict.CreateTime
	err = db.InsertDict(&dict)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusInternalServerError, "insert dict fail")
		return
	}
	// 更新搜索引擎數據
	go search.RefreshDicts([]*model.Dict{&dict})
	context.String(http.StatusOK, "add dict success")
}

func HandleAPIUpdateDict(context *gin.Context) {
	ok, _ := hasPermission(context)
	if !ok {
		context.String(http.StatusUnauthorized, "permission required, login first")
		return
	}
	var dict model.Dict
	err := json.NewDecoder(context.Request.Body).Decode(&dict)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusBadRequest, "parse request body fail")
		return
	}
	dict.UpdateTime = util.GetCurrentShanghaiTime()
	err = db.UpdateDict(&dict)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusInternalServerError, "update word fail")
		return
	}
	// 更新搜索引擎數據
	go search.RefreshDicts([]*model.Dict{&dict})
	// 移除舊的緩存
	go func() {
		cache.Cache().Remove("/dict/" + dict.Id)
	}()
	context.String(http.StatusOK, "update word success")
}

func hasPermission(context *gin.Context) (ok bool, user *userModel.User) {
	email, err := context.Cookie("email")
	if err != nil {
		log.Error(err.Error())
		return
	}
	password, err := context.Cookie("password")
	if err != nil {
		log.Error(err.Error())
		return
	}
	return userService.HadPermission(email, password)
}
