package handlers

import (
	"fmt"
	groupHandler "github.com/Solar-2020/Group-Backend/cmd/handlers/group"
	"github.com/Solar-2020/Group-Backend/internal/errorWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"runtime/debug"
)

func NewFastHttpRouter(group groupHandler.Handler, middleware Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	//router.Handle("GET", "/health", check)

	router.PanicHandler = panicHandler

	router.Handle("POST", "/group/group", middleware.CORS(group.Create))
	router.Handle("DELETE", "/group/group/:groupID", middleware.CORS(group.Delete))
	router.Handle("PUT", "/group/group/:groupID", middleware.CORS(group.Update))
	router.Handle("GET", "/group/group/:groupID", middleware.CORS(group.Get))
	router.Handle("GET", "/group/list", middleware.CORS(group.GetList))

	return router
}

func panicHandler(ctx *fasthttp.RequestCtx, err interface{}) {
	fmt.Printf("Request falied with panic: %s, error: %v\nTrace:\n", string(ctx.Request.RequestURI()), err)
	fmt.Println(string(debug.Stack()))
	errorWorker.NewErrorWorker().ServeFatalError(ctx)
}
