package user

import (
	"app/infrastructure/log"
	model2 "app/service/user/db"
	model3 "app/service/user/model"
)

func HadPermission(email string, password string) (ok bool, user *model3.User) {
	if email == "" {
		log.Error("email is empty")
		return
	}
	if password == "" {
		log.Error("password is empty")
		return
	}
	user, err := model2.FindUser(email)
	if err != nil {
		log.Error(err.Error())
		return
	}
	ok = user.Password == password
	return
}
