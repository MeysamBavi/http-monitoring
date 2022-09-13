package model

import "time"

type Alert struct {
	Id       ID
	UserId   ID
	UrlId    ID
	Url      string
	IssuedAt time.Time
}
