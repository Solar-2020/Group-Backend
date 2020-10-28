package groupHandler

import (
	"github.com/Solar-2020/GoUtils/context"
	httputils "github.com/Solar-2020/GoUtils/http"
	models2 "github.com/Solar-2020/Group-Backend/pkg/models"
)

type Handler interface {
	Create(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
	Get(ctx context.Context)
	GetList(ctx context.Context)
	//GetListInternal(ctx context.Context)
	Invite(ctx context.Context)
	EditRole(ctx context.Context)
	Expel(ctx context.Context)
	Resolve(ctx context.Context)
	AddLink(ctx context.Context)
	RemoveLink(ctx context.Context)
	ListLinks(ctx context.Context)
}

type handler struct {
	groupService   groupService
	groupTransport groupTransport
	errorWorker    errorWorker
}

func NewHandler(groupService groupService, groupTransport groupTransport, errorWorker errorWorker) Handler {
	return &handler{
		groupService:   groupService,
		groupTransport: groupTransport,
		errorWorker:    errorWorker,
	}
}

func (h *handler) Create(ctx context.Context) {
	group, err := h.groupTransport.CreateDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupReturn, err := h.groupService.Create(ctx, group)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.CreateEncode(groupReturn, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Update(ctx context.Context) {
	group, _, err := h.groupTransport.UpdateDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupReturn, err := h.groupService.Update(ctx, group)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.UpdateEncode(groupReturn, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Delete(ctx context.Context) {
	groupID, _, err := h.groupTransport.DeleteDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	group, err := h.groupService.Delete(ctx, groupID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.DeleteEncode(group, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Get(ctx context.Context) {
	groupID, _, err := h.groupTransport.GetDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	group, err := h.groupService.Get(ctx, groupID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.GetEncode(group, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) GetList(ctx context.Context) {
	_, groupID, err := h.groupTransport.GetListDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupList, err := h.groupService.GetList(ctx, groupID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.GetListEncode(groupList, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

//func (h *handler) GetListInternal(ctx context.Context) {
//	userID, groupID, err := h.groupTransport.GetListDecode(ctx.RequestCtx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	ctx.Session.Uid = userID
//
//	groupList, err := h.groupService.GetList(ctx, groupID)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	err = h.groupTransport.GetListEncode(groupList, ctx.RequestCtx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//}

func (h *handler) Invite(ctx context.Context) {
	request, err := h.groupTransport.InviteDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx, models2.Group{ID: request.Group}, models2.ActionEditRole)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.groupService.Invite(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) EditRole(ctx context.Context) {
	request, err := h.groupTransport.ChangeRoleDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx, models2.Group{ID: request.Group}, models2.ActionEditRole)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.ChangeRole(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
	}
}

func (h *handler) Expel(ctx context.Context) {
	request, err := h.groupTransport.ExpelDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx, models2.Group{ID: request.Group}, models2.ActionExpel)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.groupService.ExpelUser(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Resolve(ctx context.Context) {
	request, err := h.groupTransport.ResolveDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.groupService.ResolveGroup(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}
func (h *handler) AddLink(ctx context.Context) {
	request, err := h.groupTransport.AddLinkDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	err = h.groupService.CheckPermission(ctx, models2.Group{ID: request.Group}, models2.ActionExpel)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.AddGroupInviteLink(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}
func (h *handler) RemoveLink(ctx context.Context) {
	request, err := h.groupTransport.RemoveLinkDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx, models2.Group{ID: request.Group}, models2.ActionExpel)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.RemoveGroupInviteLink(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}
func (h *handler) ListLinks(ctx context.Context) {
	request, err := h.groupTransport.ListLinkDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx, models2.Group{ID: request.Group}, models2.ActionExpel)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.groupService.ListGroupInviteLink(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}


func (h *handler) handleError(err error, ctx context.Context) {
	err = h.errorWorker.ServeJSONError(ctx.RequestCtx, err)
	if err != nil {
		h.errorWorker.ServeFatalError(ctx.RequestCtx)
	}
	return
}