package groupHandler

import (
	"fmt"
	"github.com/Solar-2020/GoUtils/context"
	httputils "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/GoUtils/session"
	"github.com/Solar-2020/Group-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	Update(ctx *fasthttp.RequestCtx)
	Delete(ctx *fasthttp.RequestCtx)
	Get(ctx *fasthttp.RequestCtx)
	GetList(ctx *fasthttp.RequestCtx)
	Invite(ctx *fasthttp.RequestCtx)
	EditRole(ctx *fasthttp.RequestCtx)
	Expel(ctx *fasthttp.RequestCtx)
	Resolve(ctx *fasthttp.RequestCtx)
	AddLink(ctx *fasthttp.RequestCtx)
	RemoveLink(ctx *fasthttp.RequestCtx)
	ListLinks(ctx *fasthttp.RequestCtx)
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

func (h *handler) Create(ctx *fasthttp.RequestCtx) {
	req, err := h.groupTransport.CreateDecode(ctx)
	if err != nil {
		fmt.Println("Create: cannot decode request")
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, req)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupReturn, err := h.groupService.Create(ctx_, req)
	if err != nil {
		fmt.Println("Create: bad usecase: ", err)
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.groupTransport.CreateEncode(groupReturn, ctx)
	if err != nil {
		fmt.Println("Create: cannot encode response: ", err)
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Update(ctx *fasthttp.RequestCtx) {
	req, err := h.groupTransport.UpdateDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, req)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx_, req.Group, models.ActionEditRole)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupReturn, err := h.groupService.Update(ctx_, req)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.groupTransport.UpdateEncode(groupReturn, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Delete(ctx *fasthttp.RequestCtx) {
	req, err := h.groupTransport.DeleteDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, req)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx_, models.Group{ID: req.GroupID}, models.ActionEditRole)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	group, err := h.groupService.Delete(ctx_, req)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.groupTransport.DeleteEncode(group, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Get(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.GetDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx_, models.Group{ID: request.GroupID}, models.ActionEditRole)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	group, err := h.groupService.Get(ctx_, request)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.groupTransport.GetEncode(group, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) GetList(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.GetListDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	groupList, err := h.groupService.GetList(ctx_, request)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.groupTransport.GetListEncode(groupList, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Invite(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.InviteDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx_, models.Group{ID: request.Group}, models.ActionEditRole)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.groupService.Invite(ctx_, request)
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
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.groupService.CheckPermission(ctx_, models.Group{ID: request.Group}, models.ActionEditRole)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.ChangeRole(ctx_, request)
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
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	err = h.groupService.CheckPermission(ctx_, models.Group{ID: request.Group}, models.ActionExpel)

	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.ExpelUser(ctx_, request)
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

func (h *handler) Resolve(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.ResolveDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.groupService.ResolveGroup(ctx_, request)
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
func (h *handler) AddLink(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.AddLinkDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	err = h.groupService.CheckPermission(ctx_, models.Group{ID: request.Group}, models.ActionExpel)

	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.AddGroupInviteLink(ctx_, request)
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
func (h *handler) RemoveLink(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.RemoveLinkDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	err = h.groupService.CheckPermission(ctx_, models.Group{ID: request.Group}, models.ActionExpel)

	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.RemoveGroupInviteLink(ctx_, request)
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
func (h *handler) ListLinks(ctx *fasthttp.RequestCtx) {
	request, err := h.groupTransport.ListLinkDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	ctx_ := context.Context{
		RequestCtx: ctx,
		Session:    &session.Session{},
	}
	err = ctx_.Session.Authorise(ctx, request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
	err = h.groupService.CheckPermission(ctx_, models.Group{ID: request.Group}, models.ActionExpel)

	if err != nil {
		h.handleError(err, ctx)
		return
	}
	response, err := h.groupService.ListGroupInviteLink(ctx_, request)
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


func (h *handler) handleError(err error, ctx *fasthttp.RequestCtx) {
	err = h.errorWorker.ServeJSONError(ctx, err)
	if err != nil {
		h.errorWorker.ServeFatalError(ctx)
	}
	return
}