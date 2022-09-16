package model

import "go.mongodb.org/mongo-driver/bson"

type User struct {
	Id       ID     `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func (u *User) NoId() bson.M {
	return bson.M{
		"username": u.Username,
		"password": u.Password,
	}
}
