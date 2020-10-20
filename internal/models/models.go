package models

import "time"

type Group struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"URL"`
	CreateBy    int       `json:"createBy"`
	CreatAt     time.Time `json:"creatAt"`
	AvatarURL   string    `json:"avatarURL"`
	StatusID    int       `json:"-"`
}
