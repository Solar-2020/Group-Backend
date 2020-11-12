package models

import (
	"time"
)

type MemberRole int

const (
	roleCreator MemberRole = 1
	roleAdmin              = 2
	roleDweller            = 3
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
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	URL         string    `json:"URL" validate:"required"`
	CreateBy    int       `json:"createBy"`
	CreatAt     time.Time `json:"creatAt"`
	AvatarURL   string    `json:"avatarURL"`
	StatusID    int       `json:"-"`
	Count       int       `json:"count"`
	UserRole    UserRole  `json:"userRole"`
}

type Membership struct {
	UserID    int    `json:"userID"`
	GroupID   int    `json:"groupID"`
	RoleID    int    `json:"roleID"`
	RoleName  string `json:"roleName"`
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Surname   string `json:"surname"`
	AvatarURL string `json:"avatarURL"`
}

type UserRole struct {
	UserID   int    `json:"userID"`
	GroupID  int    `json:"groupID"`
	RoleID   int    `json:"roleID"`
	RoleName string `json:"roleName"`
}

type GroupPreview struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"URL"`
	AvatarURL   string     `json:"avatarURL"`
	UserID      int        `json:"userID"`
	UserRole `json:"userRole"`
	//UserRoleID  MemberRole `json:"userRoleID"`
	//UserRole    string     `json:"userRole"`
	Status      int        `json:"status"`
	Count       int        `json:"count"`
}

type AuthorPack struct {
	Login string	`json:"login"`
	ID int			`json:"id"`
}

type GroupInviteLink struct {
	Link   string    `json:"link"`
	Added  time.Time `json:"added"`
	Author AuthorPack    `json:"author"`
}
