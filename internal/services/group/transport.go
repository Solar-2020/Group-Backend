package group

import (
	"encoding/json"
	"github.com/Solar-2020/Group-Backend/internal/models"
	"github.com/go-playground/validator"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Transport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request models.Group, err error)
	CreateEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	UpdateDecode(ctx *fasthttp.RequestCtx) (request models.Group, userID int, err error)
	UpdateEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	DeleteDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error)
	DeleteEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error)
	GetEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error)

	GetListDecode(ctx *fasthttp.RequestCtx) (userID int, err error)
	GetListEncode(response []models.GroupPreview, ctx *fasthttp.RequestCtx) (err error)

	InviteDecode(ctx *fasthttp.RequestCtx) (request models.InviteUserRequest, err error)
	ChangeRoleDecode(ctx *fasthttp.RequestCtx) (request models.ChangeRoleRequest, err error)
	ExpelDecode(ctx *fasthttp.RequestCtx) (request models.ExpelUserRequest, err error)
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

func (t transport) CreateDecode(ctx *fasthttp.RequestCtx) (request models.Group, err error) {
	//userID := ctx.Value("UserID").(int)
	userID := 1
	var group models.Group
	err = json.Unmarshal(ctx.Request.Body(), &group)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	if err != nil {
		return
	}
	group.CreateBy = userID
	request = group
	return
}

func (t transport) CreateEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) UpdateDecode(ctx *fasthttp.RequestCtx) (request models.Group, userID int, err error) {
	//userID := ctx.Value("UserID").(int)
	userID = 1

	var group models.Group
	err = json.Unmarshal(ctx.Request.Body(), &group)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	if err != nil {
		return
	}

	groupIDStr := ctx.UserValue("groupID").(string)
	group.ID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}
	request = group
	return request, userID, err
}

func (t transport) UpdateEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) DeleteDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error) {
	//userID := ctx.Value("UserID").(int)
	userID = 1
	groupIDStr := ctx.UserValue("groupID").(string)
	groupID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}

	return
}

func (t transport) DeleteEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error) {
	//userID := ctx.Value("UserID").(int)
	userID = 1
	groupIDStr := ctx.UserValue("groupID").(string)
	groupID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}
	return
}

func (t transport) GetEncode(response models.Group, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetListDecode(ctx *fasthttp.RequestCtx) (userID int, err error) {
	//userID := ctx.Value("UserID").(int)
	userID = 1
	return
}

func (t transport) GetListEncode(response []models.GroupPreview, ctx *fasthttp.RequestCtx) (err error) {
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
