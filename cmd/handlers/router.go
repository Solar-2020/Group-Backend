package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	groupHandler "github.com/Solar-2020/Group-Backend/cmd/handlers/group"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(group groupHandler.Handler, middleware Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	////router.Handle("GET", "/health", check)

	router.PanicHandler = httputils.PanicHandler

	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	router.Handle("POST", "/api/group/group", middleware.Log(middleware.ExternalAuth(group.Create)))
	router.Handle("DELETE", "/api/group/group/:groupID", middleware.Log(middleware.ExternalAuth(group.Delete)))
	router.Handle("PUT", "/api/group/group/:groupID", middleware.Log(middleware.ExternalAuth(group.Update)))
	router.Handle("GET", "/api/group/group/:groupID", middleware.Log(middleware.ExternalAuth(group.Get)))
	router.Handle("GET", "/api/group/list", middleware.Log(middleware.ExternalAuth(group.GetList)))

	//router.Handle("PUT", "/api/group/membership/:groupID", middleware.Log(middleware.InternalAuth(group.Invite)))
	//router.Handle("POST", "/api/group/membership", middleware.Log(middleware.InternalAuth(group.EditRole)))
	//router.Handle("DELETE", "/api/group/membership", middleware.Log(middleware.InternalAuth(group.Expel)))
	//
	//router.Handle("PUT", "/api/group/invite/:groupID", middleware.Log(middleware.InternalAuth(group.AddLink)))
	//router.Handle("DELETE", "/api/group/invite", middleware.Log(middleware.InternalAuth(group.RemoveLink)))
	//router.Handle("POST", "/api/group/invite/list", middleware.Log(middleware.InternalAuth(group.ListLinks)))
	//router.Handle("POST", "/api/group/invite/resolve", middleware.Log(middleware.InternalAuth(group.Resolve)))

	router.Handle("GET", "/api/internal/group/list", middleware.Log(middleware.InternalAuth(group.InternalGetList)))

	return router
}
