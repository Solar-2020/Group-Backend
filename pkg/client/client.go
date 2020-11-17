package client

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Client interface {
	CheckPermission(userID, groupId, actionID int) (err error)
}

type client struct {
	host   string
	secret string
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret}
}

func (c *client) CheckPermission(userID, groupId, actionID int) (err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/group/check-permission")

	req.URI().QueryArgs().Set("user_id", strconv.Itoa(userID))
	req.URI().QueryArgs().Set("group_id", strconv.Itoa(groupId))
	req.URI().QueryArgs().Set("action_id", strconv.Itoa(actionID))

	req.Header.Set("Authorization", c.secret)
	req.Header.SetMethod(fasthttp.MethodGet)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		return
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return
		}
		return errors.New(httpErr.Error)
	case fasthttp.StatusForbidden:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return
		}
		return ResponseError{
			StatusCode: resp.StatusCode(),
			Message:    ForbiddenStatus,
			Err:        errors.Errorf(ErrorUnknownStatusCode, resp.StatusCode()),
		}
	default:
		return ResponseError{
			StatusCode: resp.StatusCode(),
			Message:    InternalServerStatus,
			Err:        errors.Errorf(ErrorUnknownStatusCode, resp.StatusCode()),
		}
	}
}
