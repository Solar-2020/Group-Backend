package groupHandler

import (
	"github.com/Solar-2020/Group-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type groupService interface {
	Create(request models.Group) (response models.Group, err error)
	Update(request models.Group, userID int) (response models.Group, err error)
	Delete(groupID, userID int) (response models.Group, err error)
	Get(groupID, userID int) (response models.Group, err error)
	GetList(userID int) (response []models.GroupPreview, err error)

	Invite(request models.InviteUserRequest) (response models.InviteUserResponse, err error)
	ChangeRole(request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error)
	ExpelUser(request models.ExpelUserRequest) (response models.ExpelUserResponse, err error)
}

type groupTransport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request models.Group, err error)
	CreateEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	UpdateDecode(ctx *fasthttp.RequestCtx) (request models.Group, userID int, err error)
	UpdateEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	DeleteDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error)
	DeleteEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error)
	GetEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	GetListDecode(ctx *fasthttp.RequestCtx) (userID int, err error)
	GetListEncode(response []models.GroupPreview, ctx *fasthttp.RequestCtx) (err error)

	InviteDecode(ctx *fasthttp.RequestCtx) (request models.InviteUserRequest, err error)
	ChangeRoleDecode(ctx *fasthttp.RequestCtx) (request models.ChangeRoleRequest, err error)
	ExpelDecode(ctx *fasthttp.RequestCtx) (request models.ExpelUserRequest, err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error) (err error)
	ServeFatalError(ctx *fasthttp.RequestCtx)
}
