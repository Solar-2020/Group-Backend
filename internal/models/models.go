package models

import (
	"time"
)

type MemberRole int
const (
	roleCreator MemberRole = 1
	roleAdmin = 2
	roleDweller = 3
)

type GroupAction int
const (
	ActionCreate GroupAction = iota
	ActionEdit
	ActionRemove
	ActionGet
	ActionInvite
	ActionEditRole
	ActionExpel
)

type Group struct {
	ID          int   	  `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	URL         string    `json:"URL" validate:"required"`
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
	UserRoleID  MemberRole    `json:"userRoleID"`
	UserRole    string `json:"userRole"`
	Status		int 	`json:"status"`
}

type GroupInviteLink struct {
	Link string	`json:"link"`
	Added time.Time		`json:"added"`
	Author int	`json:"author"`
}
