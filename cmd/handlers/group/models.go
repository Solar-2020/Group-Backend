package groupHandler

import (
	"github.com/Solar-2020/GoUtils/context"
	"github.com/Solar-2020/Group-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type groupService interface {
	Create(ctx context.Context, request models.CreateRequest) (response models.CreateResponse, err error)
	Update(ctx context.Context, request models.UpdateRequest) (response models.UpdateResponse, err error)
	Delete(ctx context.Context, request models.DeleteRequest) (response models.DeleteResponse, err error)
	Get(ctx context.Context, request models.GetRequest) (response models.GetResponse, err error)
	GetList(ctx context.Context, request models.GetListRequest) (response models.GetListResponse, err error)

	Invite(ctx context.Context, request models.InviteUserRequest) (response models.InviteUserResponse, err error)
	ChangeRole(ctx context.Context, request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error)
	ExpelUser(ctx context.Context, request models.ExpelUserRequest) (response models.ExpelUserResponse, err error)

	ResolveGroup(ctx context.Context, request models.ResolveInviteLinkRequest) (response models.ResolveInviteLinkResponse, err error)
	AddGroupInviteLink(ctx context.Context, request models.AddInviteLinkRequest) (response models.AddInviteLinkResponse, err error)
	RemoveGroupInviteLink(ctx context.Context, request models.RemoveInviteLinkRequest) (response models.RemoveInviteLinkRsponse, err error)
	ListGroupInviteLink(ctx context.Context, request models.ListInviteLinkRequest) (response models.ListInviteLinkResponse, err error)


	CheckPermission(ctx context.Context, group models.Group, action models.GroupAction) error
}

type groupTransport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request models.CreateRequest, err error)
	CreateEncode(response models.CreateResponse, ctx *fasthttp.RequestCtx) (err error)

	UpdateDecode(ctx *fasthttp.RequestCtx) (request models.UpdateRequest, err error)
	UpdateEncode(response models.UpdateResponse, ctx *fasthttp.RequestCtx) (err error)

	DeleteDecode(ctx *fasthttp.RequestCtx) (request models.DeleteRequest, err error)
	DeleteEncode(response models.DeleteResponse, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (request models.GetRequest, err error)
	GetEncode(response models.GetResponse, ctx *fasthttp.RequestCtx) (err error)

	GetListDecode(ctx *fasthttp.RequestCtx) (request models.GetListRequest, err error)
	GetListEncode(response models.GetListResponse, ctx *fasthttp.RequestCtx) (err error)

	InviteDecode(ctx *fasthttp.RequestCtx) (request models.InviteUserRequest, err error)
	ChangeRoleDecode(ctx *fasthttp.RequestCtx) (request models.ChangeRoleRequest, err error)
	ExpelDecode(ctx *fasthttp.RequestCtx) (request models.ExpelUserRequest, err error)

	ResolveDecode(ctx *fasthttp.RequestCtx) (request models.ResolveInviteLinkRequest, err error)
	AddLinkDecode(ctx *fasthttp.RequestCtx) (request models.AddInviteLinkRequest, err error)
	RemoveLinkDecode(ctx *fasthttp.RequestCtx) (request models.RemoveInviteLinkRequest, err error)
	ListLinkDecode(ctx *fasthttp.RequestCtx) (request models.ListInviteLinkRequest, err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error) (err error)
	ServeFatalError(ctx *fasthttp.RequestCtx)
}
