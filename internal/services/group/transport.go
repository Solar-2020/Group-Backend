package group

import (
	"encoding/json"
	"github.com/Solar-2020/Group-Backend/internal/models"
	"github.com/go-playground/validator"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Transport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request models.CreateRequest, err error)
	CreateEncode(response models.CreateResponse, ctx *fasthttp.RequestCtx) (err error)

	UpdateDecode(ctx *fasthttp.RequestCtx) (request models.UpdateRequest, err error)
	UpdateEncode(response models.UpdateResponse, ctx *fasthttp.RequestCtx) (err error)

	DeleteDecode(ctx *fasthttp.RequestCtx) (request models.DeleteRequest, err error)
	DeleteEncode(response models.DeleteResponse, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (request models.GetRequest, err error)
	GetEncode(response models.GetResponse, ctx *fasthttp.RequestCtx) (err error)

	GetListDecode(ctx *fasthttp.RequestCtx) (request models.GetListRequest, err error)
	GetListEncode(response models.GetListResponse, ctx *fasthttp.RequestCtx) (err error)

	InviteDecode(ctx *fasthttp.RequestCtx) (request models.InviteUserRequest, err error)
	ChangeRoleDecode(ctx *fasthttp.RequestCtx) (request models.ChangeRoleRequest, err error)
	ExpelDecode(ctx *fasthttp.RequestCtx) (request models.ExpelUserRequest, err error)

	ResolveDecode(ctx *fasthttp.RequestCtx) (request models.ResolveInviteLinkRequest, err error)
	AddLinkDecode(ctx *fasthttp.RequestCtx) (request models.AddInviteLinkRequest, err error)
	RemoveLinkDecode(ctx *fasthttp.RequestCtx) (request models.RemoveInviteLinkRequest, err error)
	ListLinkDecode(ctx *fasthttp.RequestCtx) (request models.ListInviteLinkRequest, err error)
}

const (
	DefaultUid = 1
)

type transport struct {
	validator *validator.Validate
}

func NewTransport() Transport {
	return &transport{
		validator: validator.New(),
	}
}

func (t transport) CreateDecode(ctx *fasthttp.RequestCtx) (request models.CreateRequest, err error) {
	//userID := ctx.Value("UserID").(int)
	//userID := 1
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	if err != nil {
		return
	}
	tmp := ctx.Value("UserID")
	if tmp != nil && request.CreateBy != 0{
		request.CreateBy = tmp.(int)
	}
	//request.CreateBy = userID
	return
}

func (t transport) CreateEncode(response models.CreateResponse, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) UpdateDecode(ctx *fasthttp.RequestCtx) (request models.UpdateRequest, err error) {
	//userID = 1

	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	if err != nil {
		return
	}

	tmp := ctx.Value("UserID")
	if tmp != nil && request.UserID != 0{
		request.UserID = tmp.(int)
	}

	groupIDStr := ctx.UserValue("groupID").(string)
	request.ID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}
	return request, err
}

func (t transport) UpdateEncode(response models.UpdateResponse, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) DeleteDecode(ctx *fasthttp.RequestCtx) (request models.DeleteRequest, err error) {
	//userID := ctx.Value("UserID").(int)
	//userID = 1
	tmp := ctx.Value("UserID")
	if tmp != nil {
		request.UserID = tmp.(int)
	}
	groupIDStr := ctx.UserValue("groupID").(string)
	request.GroupID, err = strconv.Atoi(groupIDStr)
	return
}

func (t transport) DeleteEncode(response models.DeleteResponse, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetDecode(ctx *fasthttp.RequestCtx) (request models.GetRequest, err error) {
	//userID := ctx.Value("UserID").(int)
	//userID = 1
	tmp := ctx.Value("UserID")
	if tmp != nil {
		request.UserID = tmp.(int)
	}
	groupIDStr := ctx.UserValue("groupID").(string)
	request.GroupID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}
	return
}

func (t transport) GetEncode(response models.GetResponse, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetListDecode(ctx *fasthttp.RequestCtx) (request models.GetListRequest, err error) {
	//userID := ctx.Value("UserID").(int)
	//userID = 1
	tmp := ctx.Value("UserID")
	if tmp != nil {
		request.UserID = tmp.(int)
	}
	return
}

func (t transport) GetListEncode(response models.GetListResponse, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) InviteDecode(ctx *fasthttp.RequestCtx) (request models.InviteUserRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	if request.Uid == 0 {
		request.Uid = DefaultUid
	}
	return
}

func (t transport) ChangeRoleDecode(ctx *fasthttp.RequestCtx) (request models.ChangeRoleRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	if request.Uid == 0 {
		request.Uid = DefaultUid
	}
	return
}

func (t transport) ExpelDecode(ctx *fasthttp.RequestCtx) (request models.ExpelUserRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	if request.Uid == 0 {
		request.Uid = DefaultUid
	}
	return
}

func (t transport) ResolveDecode(ctx *fasthttp.RequestCtx) (request models.ResolveInviteLinkRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) AddLinkDecode(ctx *fasthttp.RequestCtx) (request models.AddInviteLinkRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}
func (t transport) RemoveLinkDecode(ctx *fasthttp.RequestCtx) (request models.RemoveInviteLinkRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}
func (t transport) ListLinkDecode(ctx *fasthttp.RequestCtx) (request models.ListInviteLinkRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}
