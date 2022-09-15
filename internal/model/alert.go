package model

import "time"

type Alert struct {
	Id       ID        `json:"-" bson:"_id"`
	UserId   ID        `json:"-" bson:"user_id"`
	UrlId    ID        `json:"url_id" bson:"url_id"`
	Url      string    `json:"url" bson:"url"`
	IssuedAt time.Time `json:"issued_at" bson:"issued_at"`
}
