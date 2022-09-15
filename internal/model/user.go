package model

type User struct {
	Id       ID     `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}
