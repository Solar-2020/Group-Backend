package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	groupHandler "github.com/Solar-2020/Group-Backend/cmd/handlers/group"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(group groupHandler.Handler, middleware httputils.Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	//router.Handle("GET", "/health", check)

	router.PanicHandler = httputils.PanicHandler
	middlewareChain := httputils.NewLogCorsChain(middleware)
	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	router.Handle("POST", "/group/group", middlewareChain(group.Create))
	router.Handle("DELETE", "/group/group/:groupID", middlewareChain(group.Delete))
	router.Handle("PUT", "/group/group/:groupID", middlewareChain(group.Update))
	router.Handle("GET", "/group/group/:groupID", middlewareChain(group.Get))
	router.Handle("GET", "/group/list", middlewareChain(group.GetList))

	router.Handle("PUT", "/group/membership/:groupID", middlewareChain(group.Invite))
	router.Handle("POST", "/group/membership", middlewareChain(group.EditRole))
	router.Handle("DELETE", "/group/membership", middlewareChain(group.Expel))

	router.Handle("PUT", "/group/invite/:groupID", middlewareChain(group.AddLink))
	router.Handle("DELETE", "/group/invite", middlewareChain(group.RemoveLink))
	router.Handle("POST", "/group/invite/list", middlewareChain(group.ListLinks))
	router.Handle("POST", "/group/invite/resolve", middlewareChain(group.Resolve))

	return router
}
