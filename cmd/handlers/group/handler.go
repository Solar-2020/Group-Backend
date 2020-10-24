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
	fmt.Println("New incoming request: POST /group/group")
	group, err := h.groupTransport.CreateDecode(ctx)
	if err != nil {
		fmt.Println("Create: cannot decode request")
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	groupReturn, err := h.groupService.Create(group)
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
	group, userID, err := h.groupTransport.UpdateDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	groupReturn, err := h.groupService.Update(group, userID)
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
	groupID, userID, err := h.groupTransport.DeleteDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	group, err := h.groupService.Delete(groupID, userID)
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
	groupID, userID, err := h.groupTransport.GetDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	group, err := h.groupService.Get(groupID, userID)
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
	userID, err := h.groupTransport.GetListDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	groupList, err := h.groupService.GetList(userID)
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

func (h *handler) handleError(err error, ctx *fasthttp.RequestCtx) {
	err = h.errorWorker.ServeJSONError(ctx, err)
	if err != nil {
		h.errorWorker.ServeFatalError(ctx)
	}
	return
}