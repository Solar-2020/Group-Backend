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

	router.Handle("POST", "/api/group/group", clientside(group.Create))
	router.Handle("DELETE", "/api/group/group/:groupID", clientside(group.Delete))
	router.Handle("PUT", "/api/group/group/:groupID", clientside(group.Update))
	router.Handle("GET", "/api/group/group/:groupID", clientside(group.Get))
	router.Handle("GET", "/api/group/list", clientside(group.GetList))

	router.Handle("PUT", "/api/group/membership/:groupID", clientside(group.Invite))
	router.Handle("POST", "/api/group/membership", clientside(group.EditRole))
	router.Handle("DELETE", "/api/group/membership", clientside(group.Expel))

	router.Handle("PUT", "/api/group/invite/:groupID", clientside(group.AddLink))
	router.Handle("DELETE", "/api/group/invite", clientside(group.RemoveLink))
	router.Handle("POST", "/api/group/invite/list", clientside(group.ListLinks))
	router.Handle("POST", "/api/group/invite/resolve", clientside(group.Resolve))

	serverside := httputils.ServersideChain(middleware)
	router.Handle("GET", "/api/internal/group/list", serverside(group.GetList))

	return router
}
