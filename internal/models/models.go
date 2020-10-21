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

type GroupPreview struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"URL"`
	AvatarURL   string `json:"avatarURL"`
	UserID      int    `json:"userID"`
	UserRoleID  int    `json:"userRoleID"`
	UserRole    string `json:"userRole"`
}
