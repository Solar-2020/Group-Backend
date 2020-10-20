package groupHandler

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	Update(ctx *fasthttp.RequestCtx)
	Delete(ctx *fasthttp.RequestCtx)
	Get(ctx *fasthttp.RequestCtx)
	GetList(ctx *fasthttp.RequestCtx)
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
	group, err := h.groupTransport.CreateDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	groupReturn, err := h.groupService.Create(group)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.groupTransport.CreateEncode(groupReturn, ctx)
	if err != nil {
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
