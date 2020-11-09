package groupHandler

import (
	"errors"
	httputils "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/Group-Backend/internal/services/group"
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	Update(ctx *fasthttp.RequestCtx)
	Delete(ctx *fasthttp.RequestCtx)
	Get(ctx *fasthttp.RequestCtx)
	GetList(ctx *fasthttp.RequestCtx)
	InternalGetList(ctx *fasthttp.RequestCtx)
	InternalGetPermission(ctx *fasthttp.RequestCtx)
	GetMembershipList(ctx *fasthttp.RequestCtx)
	Invite(ctx *fasthttp.RequestCtx)
	EditRole(ctx *fasthttp.RequestCtx)
	Expel(ctx *fasthttp.RequestCtx)
	//Resolve(ctx *fasthttp.RequestCtx)
	//AddLink(ctx *fasthttp.RequestCtx)
	//RemoveLink(ctx *fasthttp.RequestCtx)
	//ListLinks(ctx *fasthttp.RequestCtx)
}

type handler struct {
	groupService   group.Service
	groupTransport group.Transport
	errorWorker    errorWorker
}

func NewHandler(groupService group.Service, groupTransport group.Transport, errorWorker errorWorker) Handler {
	return &handler{
		groupService:   groupService,
		groupTransport: groupTransport,
		errorWorker:    errorWorker,
	}
}

func (h *handler) Create(ctx *fasthttp.RequestCtx) {
	group, err := h.groupTransport.CreateDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupReturn, err := h.groupService.Create(group)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.CreateEncode(groupReturn, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Update(ctx *fasthttp.RequestCtx) {
	group, userID, err := h.groupTransport.UpdateDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupReturn, err := h.groupService.Update(group, userID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.UpdateEncode(groupReturn, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Delete(ctx *fasthttp.RequestCtx) {
	groupID, userID, err := h.groupTransport.DeleteDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	group, err := h.groupService.Delete(groupID, userID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.DeleteEncode(group, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Get(ctx *fasthttp.RequestCtx) {
	groupID, userID, err := h.groupTransport.GetDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	group, err := h.groupService.Get(groupID, userID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.GetEncode(group, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) GetList(ctx *fasthttp.RequestCtx) {
	userID, groupID, err := h.groupTransport.GetListDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupList, err := h.groupService.GetList(groupID, userID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.GetListEncode(groupList, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) InternalGetList(ctx *fasthttp.RequestCtx) {
	userID, groupID, err := h.groupTransport.InternalGetListDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupList, err := h.groupService.GetList(groupID, userID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.InternalGetListEncode(groupList, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) InternalGetPermission(ctx *fasthttp.RequestCtx) {
	userID, groupID, err := h.groupTransport.InternalGetPermissionDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupList, err := h.groupService.GetUserRole(groupID, userID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.InternalGetPermissionEncode(groupList, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

//func (h *handler) GetListInternal(ctx *fasthttp.RequestCtx) {
//	userID, groupID, err := h.groupTransport.GetListDecode(ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	ctx.Session.Uid = userID
//
//	groupList, err := h.groupService.GetList(groupID)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	err = h.groupTransport.GetListEncode(groupList, ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//}

func (h *handler) Invite(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.InviteDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	role, err := h.groupService.GetUserRole(request.Group, request.CreatorID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	if !(role.RoleID == 1 || role.RoleID == 2) {
		h.handleError(errors.New("access denied"), ctx)
		return
	}

	response, err := h.groupService.Invite(request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) EditRole(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.ChangeRoleDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	//err = h.groupService.CheckPermission(models2.Group{ID: request.Group}, models2.ActionEditRole)
	//if err != nil {
	//	h.handleError(err, ctx)
	//	return
	//}
	response, err := h.groupService.ChangeRole(request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx)
	if err != nil {
		h.handleError(err, ctx)
	}
}

func (h *handler) Expel(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.ExpelDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	//err = h.groupService.CheckPermission(models2.Group{ID: request.Group}, models2.ActionExpel)
	//if err != nil {
	//	h.handleError(err, ctx)
	//	return
	//}

	response, err := h.groupService.ExpelUser(request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = httputils.EncodeDefault(response, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) GetMembershipList(ctx *fasthttp.RequestCtx) {
	userID, groupID, err := h.groupTransport.GetMembershipListDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.groupService.GetMembershipList(userID, groupID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupTransport.GetMembershipListEncode(response, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

//func (h *handler) Resolve(ctx *fasthttp.RequestCtx) {
//	request, err := h.groupTransport.ResolveDecode(ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	response, err := h.groupService.ResolveGroup(request)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	err = httputils.EncodeDefault(response, ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//}
//func (h *handler) AddLink(ctx *fasthttp.RequestCtx) {
//	request, err := h.groupTransport.AddLinkDecode(ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//	//err = h.groupService.CheckPermission(models2.Group{ID: request.Group}, models2.ActionExpel)
//	//if err != nil {
//	//	h.handleError(err, ctx)
//	//	return
//	//}
//	response, err := h.groupService.AddGroupInviteLink(request)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	err = httputils.EncodeDefault(response, ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//}
//func (h *handler) RemoveLink(ctx *fasthttp.RequestCtx) {
//	request, err := h.groupTransport.RemoveLinkDecode(ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	//err = h.groupService.CheckPermission(models2.Group{ID: request.Group}, models2.ActionExpel)
//	//if err != nil {
//	//	h.handleError(err, ctx)
//	//	return
//	//}
//	response, err := h.groupService.RemoveGroupInviteLink(request)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	err = httputils.EncodeDefault(response, ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//}
//func (h *handler) ListLinks(ctx *fasthttp.RequestCtx) {
//	request, err := h.groupTransport.ListLinkDecode(ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	//err = h.groupService.CheckPermission(models2.Group{ID: request.Group}, models2.ActionExpel)
//	//if err != nil {
//	//	h.handleError(err, ctx)
//	//	return
//	//}
//
//	response, err := h.groupService.ListGroupInviteLink(request)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//
//	err = httputils.EncodeDefault(response, ctx)
//	if err != nil {
//		h.handleError(err, ctx)
//		return
//	}
//}

func (h *handler) handleError(err error, ctx *fasthttp.RequestCtx) {
	err = h.errorWorker.ServeJSONError(ctx, err)
	if err != nil {
		h.errorWorker.ServeFatalError(ctx)
	}
	return
}
