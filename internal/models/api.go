package models

import (
	httputils "github.com/Solar-2020/GoUtils/http"
)

// POST /group/group
type CreateRequest struct {
	httputils.Authorized
	Group
}
type CreateResponse struct {
	Group
}

// DELETE /group/group/:groupID
type UpdateRequest struct {
	httputils.Authorized
	Group
	UserID int
}
type UpdateResponse struct {
	Group
}

// PUT /group/group/:groupID
type DeleteRequest struct {
	httputils.Authorized
	GroupID int
	UserID int
}
type DeleteResponse struct {
	Group
}

// GET /group/group/:groupID
type GetRequest struct {
	httputils.Authorized
	GroupID int
	UserID int
}
type GetResponse struct {
	Group
}

// GET /group/list
type GetListRequest struct {
	httputils.Authorized
	UserID int
}
type GetListResponse struct {
	UserID int				`json:"user"`
	Groups []GroupPreview		`json:"groups"`
}

// PUT /group/membership
type InviteUserRequest struct {
	httputils.Authorized
	UserID []int 	`json:"userId"`
	Group int	`json:"group" validate:"required"`
	User []string		`json:"userEmail" validate:"required"`
	Role MemberRole	`json:"role"`
}
type InviteUserResponse InviteUserRequest

// POST /group/membership
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

// DELETE /group/membership
type ExpelUserRequest struct {
	httputils.Authorized
	UserID int 	`json:"userId"`
	Group int	`json:"group" validate:"required"`
	User string		`json:"userEmail" validate:"required,email"`
}
type ExpelUserResponse struct {
	User string `json:"userEmail"`
}

// PUT /group/invite
type AddInviteLinkRequest struct {
	httputils.Authorized
	Group int	`json:"group" validate:"required"`
}
type AddInviteLinkResponse struct {
	Group int	`json:"group"`
	Link string `json:"link"`
}

// DELETE /group/invite
type RemoveInviteLinkRequest struct {
	httputils.Authorized
	Group int      `json:"group" validate:"required"`
	Links []string `json:"links" validate:"required"`
}
type RemoveInviteLinkRsponse struct {
	Group int      `json:"group"`
	Links []string `json:"links"`
}

// POST /group/invite/list
type ListInviteLinkRequest struct {
	httputils.Authorized
	Group int `json:"group" validate:"required"`
}
type ListInviteLinkResponse struct {
	Group int `json:"group"`
	Links []GroupInviteLink	`json:"links"`
}

// POST /group/invite/resolves
type ResolveInviteLinkRequest struct {
	httputils.Authorized
	Link string `json:"link"`
}
type ResolveInviteLinkResponse struct {
	Group int `json:"group"`
}