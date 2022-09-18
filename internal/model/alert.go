package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Alert struct {
	Id       ID        `json:"-" bson:"_id"`
	UserId   ID        `json:"-" bson:"user_id"`
	UrlId    ID        `json:"url_id" bson:"url_id"`
	Url      string    `json:"url" bson:"url"`
	IssuedAt time.Time `json:"issued_at" bson:"issued_at"`
}

func (a *Alert) NoId() bson.M {
	return bson.M{
		"user_id":   a.UserId,
		"url_id":    a.UrlId,
		"url":       a.Url,
		"issued_at": a.IssuedAt,
	}
}
