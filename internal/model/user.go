package model

import "go.mongodb.org/mongo-driver/bson"

type User struct {
	Id       ID     `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func (u *User) NoId() bson.M {
	return bson.M{
		"username": u.Username,
		"password": u.Password,
	}
}
