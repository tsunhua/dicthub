package model

type User struct {
	Id          string `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Email       string `bson:"email" json:"email"`
	Password    string `bson:"password" json:"password"`
	NotifyEmail string `bson:"notifyEmail" json:"notifyEmail"` // 提醒專用郵箱，默認取 Email 值
}
