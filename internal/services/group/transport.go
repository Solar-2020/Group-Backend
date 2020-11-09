package group

import (
	"encoding/json"
	"errors"
	"github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/Group-Backend/internal/models"
	models2 "github.com/Solar-2020/Group-Backend/pkg/models"
	"github.com/go-playground/validator"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Transport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request models2.Group, err error)
	CreateEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error)

	UpdateDecode(ctx *fasthttp.RequestCtx) (request models2.Group, userID int, err error)
	UpdateEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error)

	DeleteDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error)
	DeleteEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (groupID, userID int, err error)
	GetEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error)

	GetListDecode(ctx *fasthttp.RequestCtx) (userID, groupID int, err error)
	GetListEncode(response []models2.GroupPreview, ctx *fasthttp.RequestCtx) (err error)

	InternalGetListDecode(ctx *fasthttp.RequestCtx) (userID, groupID int, err error)
	InternalGetListEncode(response []models2.GroupPreview, ctx *fasthttp.RequestCtx) (err error)

	InternalGetPermissionDecode(ctx *fasthttp.RequestCtx) (userID, groupID int, err error)
	InternalGetPermissionEncode(response models2.UserRole, ctx *fasthttp.RequestCtx) (err error)

	InviteDecode(ctx *fasthttp.RequestCtx) (request models.InviteUserRequest, err error)

	ChangeRoleDecode(ctx *fasthttp.RequestCtx) (request models.ChangeRoleRequest, err error)
	ExpelDecode(ctx *fasthttp.RequestCtx) (request models.ExpelUserRequest, err error)

	ResolveDecode(ctx *fasthttp.RequestCtx) (request models.ResolveInviteLinkRequest, err error)
	AddLinkDecode(ctx *fasthttp.RequestCtx) (request models.AddInviteLinkRequest, userID int, err error)
	RemoveLinkDecode(ctx *fasthttp.RequestCtx) (request models.RemoveInviteLinkRequest, err error)
	ListLinkDecode(ctx *fasthttp.RequestCtx) (request models.ListInviteLinkRequest, err error)
}

type transport struct {
	validator *validator.Validate
}

func NewTransport() Transport {
	return &transport{
		validator: validator.New(),
	}
}

func (t transport) CreateDecode(ctx *fasthttp.RequestCtx) (request models2.Group, err error) {
	var group models2.Group
	err = json.Unmarshal(ctx.Request.Body(), &group)
	if err != nil {
		return
	}

	err = t.validator.Struct(group)
	if err != nil {
		return
	}

	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		group.CreateBy = userID
		request = group
		return
	}

	return request, errors.New("userID not found")
}

func (t transport) CreateEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) UpdateDecode(ctx *fasthttp.RequestCtx) (request models2.Group, userID int, err error) {
	var group models2.Group
	var ok bool
	err = json.Unmarshal(ctx.Request.Body(), &group)
	if err != nil {
		return
	}

	err = t.validator.Struct(group)
	if err != nil {
		return
	}

	groupIDStr := ctx.UserValue("groupID").(string)
	group.ID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}

	userID, ok = ctx.UserValue("userID").(int)
	if ok {
		request = group
		return
	}


	return request, userID, errors.New("userID not found")
}

func (t transport) UpdateEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error) {
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
	var ok bool
	groupIDStr := ctx.UserValue("groupID").(string)
	groupID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}

	userID, ok = ctx.UserValue("userID").(int)
	if ok {
		return
	}

	return groupID, userID, errors.New("userID not found")
}

func (t transport) DeleteEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error) {
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
	var ok bool
	groupIDStr := ctx.UserValue("groupID").(string)
	groupID, err = strconv.Atoi(groupIDStr)
	if err != nil {
		return
	}

	userID, ok = ctx.UserValue("userID").(int)
	if ok {
		return
	}

	return groupID, userID, errors.New("userID not found")
}

func (t transport) GetEncode(response models2.Group, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetListDecode(ctx *fasthttp.RequestCtx) (userID, groupID int, err error) {
	var ok bool
	_group := ctx.QueryArgs().Peek("group_id")
	if _group != nil {
		groupID, _ = strconv.Atoi(string(_group))
	}

	userID, ok = ctx.UserValue("userID").(int)
	if ok {
		return
	}

	return groupID, userID, errors.New("userID not found")
}

func (t transport) GetListEncode(response []models2.GroupPreview, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) InternalGetListDecode(ctx *fasthttp.RequestCtx) (userID, groupID int, err error) {
	_group := ctx.QueryArgs().Peek("group_id")
	if _group != nil {
		groupID, _ = strconv.Atoi(string(_group))
	}

	_userID := ctx.QueryArgs().Peek("user_id")
	if _userID != nil {
		userID, _ = strconv.Atoi(string(_userID))
	}

	return
}

func (t transport) InternalGetListEncode(response []models2.GroupPreview, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) InternalGetPermissionDecode(ctx *fasthttp.RequestCtx) (userID, groupID int, err error) {
	_group := ctx.QueryArgs().Peek("group_id")
	if _group != nil {
		groupID, _ = strconv.Atoi(string(_group))
	}

	_userID := ctx.QueryArgs().Peek("user_id")
	if _userID != nil {
		userID, _ = strconv.Atoi(string(_userID))
	}

	return
}

func (t transport) InternalGetPermissionEncode(response models2.UserRole, ctx *fasthttp.RequestCtx) (err error) {
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
	var ok bool
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}

	if request.Group == 0 {
		if urlId, err := http.GetUrlParamInt(ctx, "groupID"); err == nil {
			request.Group = urlId
		}
	}
	err = t.validator.Struct(request)

	request.CreatorID, ok = ctx.UserValue("userID").(int)
	if ok {
		return
	}

	return request, errors.New("userID not found")
}

func (t transport) ChangeRoleDecode(ctx *fasthttp.RequestCtx) (request models.ChangeRoleRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) ExpelDecode(ctx *fasthttp.RequestCtx) (request models.ExpelUserRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
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

func (t transport) AddLinkDecode(ctx *fasthttp.RequestCtx) (request models.AddInviteLinkRequest, userID int, err error) {
	body := ctx.Request.Body()
	if body != nil && len(body) > 0 {
		err = json.Unmarshal(ctx.Request.Body(), &request)
		if err != nil {
			return
		}
	}
	if request.Group == 0 {
		if urlId, err := http.GetUrlParamInt(ctx, "groupID"); err == nil {
			request.Group = urlId
		}
	}
	err = t.validator.Struct(request)
	var ok bool
	userID, ok = ctx.UserValue("userID").(int)
	if ok {
		return
	}
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
