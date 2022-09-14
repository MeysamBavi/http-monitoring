package model

import "time"

type Alert struct {
	Id       ID        `json:"-"`
	UserId   ID        `json:"-"`
	UrlId    ID        `json:"url_id"`
	Url      string    `json:"url"`
	IssuedAt time.Time `json:"issued_at"`
}
