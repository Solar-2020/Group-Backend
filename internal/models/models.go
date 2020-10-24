package models

import (
	httputils "github.com/Solar-2020/GoUtils/http"
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

type CreateRequest struct {
	httputils.Authorized
	Group
}

type CreateResponse struct {
	Group
}

type UpdateRequest struct {
	httputils.Authorized
	Group
	UserID int
}

type UpdateResponse struct {
	Group
}

type DeleteRequest struct {
	httputils.Authorized
	GroupID int
	UserID int
}
type DeleteResponse struct {
	Group
}

type GetRequest struct {
	httputils.Authorized
	GroupID int
	UserID int
}

type GetResponse struct {
	Group
}

type GetListRequest struct {
	httputils.Authorized
	UserID int
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

type GetListResponse struct {
	UserID int				`json:"user"`
	Groups []GroupPreview		`json:"groups"`
}


type InviteUserRequest struct {
	httputils.Authorized
	UserID []int 	`json:"userId"`
	Group int	`json:"group" validate:"required"`
	User []string		`json:"userEmail" validate:"required"`
	Role MemberRole	`json:"role"`
}

type InviteUserResponse InviteUserRequest

type ChangeRoleRequest struct {
	httputils.Authorized
	UserID int 	`json:"userId"`
	Group int	`json:"group" validate:"required"`
	User string		`json:"userEmail" validate:"required,email"`
	Role MemberRole	`json:"role"`
}

type ChangeRoleResponse struct {
	Role MemberRole	`json:"role"`
}

type ExpelUserRequest struct {
	httputils.Authorized
	UserID int 	`json:"userId"`
	Group int	`json:"group" validate:"required"`
	User string		`json:"userEmail" validate:"required,email"`
}

type ExpelUserResponse struct {
	User string	`json:"userEmail"`
}

type GroupInviteLink struct {
	Link string	`json:"link"`
	Added time.Time		`json:"added"`
	Author int	`json:"author"`
}

type AddInviteLinkRequest struct {
	httputils.Authorized
	Group int	`json:"group" validate:"required"`
}

type AddInviteLinkResponse struct {
	Group int	`json:"group"`
	Link string `json:"link"`
}


type RemoveInviteLinkRequest struct {
	httputils.Authorized
	Group int      `json:"group" validate:"required"`
	Links []string `json:"links" validate:"required"`
}

type RemoveInviteLinkRsponse struct {
	Group int      `json:"group"`
	Links []string `json:"links"`
}

type ListInviteLinkRequest struct {
	httputils.Authorized
	Group int `json:"group" validate:"required"`
}

type ListInviteLinkResponse struct {
	Group int `json:"group"`
	Links []GroupInviteLink	`json:"links"`
}

type ResolveInviteLinkRequest struct {
	httputils.Authorized
	Link string `json:"link"`
}

type ResolveInviteLinkResponse struct {
	Group int `json:"group"`
}