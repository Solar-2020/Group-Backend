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
	clientside := httputils.ClientsideChain(middleware)
	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	router.Handle("POST", "/group/group", clientside(group.Create))
	router.Handle("DELETE", "/group/group/:groupID", clientside(group.Delete))
	router.Handle("PUT", "/group/group/:groupID", clientside(group.Update))
	router.Handle("GET", "/group/group/:groupID", clientside(group.Get))
	router.Handle("GET", "/group/list", clientside(group.GetList))

	router.Handle("PUT", "/group/membership/:groupID", clientside(group.Invite))
	router.Handle("POST", "/group/membership", clientside(group.EditRole))
	router.Handle("DELETE", "/group/membership", clientside(group.Expel))

	router.Handle("PUT", "/group/invite/:groupID", clientside(group.AddLink))
	router.Handle("DELETE", "/group/invite", clientside(group.RemoveLink))
	router.Handle("POST", "/group/invite/list", clientside(group.ListLinks))
	router.Handle("POST", "/group/invite/resolve", clientside(group.Resolve))

	serverside := httputils.ServersideChain(middleware)
	router.Handle("GET", "/internal/group/list", serverside(group.GetList))

	return router
}
