package models

import (
	"fmt"
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
	ID          int   `json:"id"`
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
	UserRoleID  MemberRole    `json:"userRoleID"`
	UserRole    string `json:"userRole"`
}

type Authorized struct {
	Uid int `json:"uid"`
}

func (r Authorized) GetUid() (int, error) {
	if r.Uid == 0 {
		return 0, fmt.Errorf("No uid")
	}
	return r.Uid, nil
}

type UserIdMock struct {
	UserID int 	`json:"userId"`
}

type InviteUserRequest struct {
	Authorized
	UserIdMock
	Group int	`json:"group" validate:"required"`
	User string		`json:"userEmail" validate:"required,email"`
	Role MemberRole	`json:"role"`
}

type InviteUserResponse InviteUserRequest

type ChangeRoleRequest InviteUserRequest

type ChangeRoleResponse struct {
	Role MemberRole	`json:"role"`
}

type ExpelUserRequest struct {
	Authorized
	UserIdMock
	Group int	`json:"group" validate:"required"`
	User string		`json:"userEmail" validate:"required,email"`
}

type ExpelUserResponse struct {
	User string	`json:"userEmail"`
}