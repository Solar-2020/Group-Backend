package models

// PUT /group/membership
type InviteUserRequest struct {
	UserID []int 	`json:"userId"`
	Group int	`json:"group" validate:"required"`
	User []string		`json:"userEmail"`
	Role MemberRole	`json:"role"`
}
type InviteUserResponse InviteUserRequest

// POST /group/membership
type ChangeRoleRequest struct {
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
	UserID int 	`json:"userId"`
	Group int	`json:"group" validate:"required"`
	User string		`json:"userEmail" validate:"required,email"`
}
type ExpelUserResponse struct {
	User string `json:"userEmail"`
}

// PUT /group/invite
type AddInviteLinkRequest struct {
	Group int	`json:"group"`
}
type AddInviteLinkResponse struct {
	Group int	`json:"group"`
	Link string `json:"link"`
}

// DELETE /group/invite
type RemoveInviteLinkRequest struct {
	Group int      `json:"group" validate:"required"`
	Links []string `json:"links" validate:"required"`
}
type RemoveInviteLinkRsponse struct {
	Group int      `json:"group"`
	Links []string `json:"links"`
}

// POST /group/invite/list
type ListInviteLinkRequest struct {
	Group int `json:"group" validate:"required"`
}
type ListInviteLinkResponse struct {
	Group int `json:"group"`
	Links []GroupInviteLink	`json:"links"`
}

// POST /group/invite/resolves
type ResolveInviteLinkRequest struct {
	Link string `json:"link"`
}
type ResolveInviteLinkResponse struct {
	Group int `json:"group"`
}