package api

import (
	service "github.com/Solar-2020/GoUtils/http"
	models2 "github.com/Solar-2020/Group-Backend/pkg/models"
	urllib "net/url"
	"strconv"
)

type GroupServiceInterface interface {
	//UsersGroups(userID int) ([]models.GroupPreview, error)
	UserGroupPreview(userID int, groupID int) (models2.GroupPreview, error)
}

type GroupClient struct {
	service.Service
	Addr string
}
func (c *GroupClient) Address () string { return c.Addr }
type UsersGroupsRequest struct {
	Uid int `json:"uid"`
	GroupID int `json:"group_id"`
}
func (r *UsersGroupsRequest) QueryParams() (values urllib.Values, err error) {
	values = urllib.Values{}
	values.Set("uid", strconv.Itoa(r.Uid))
	values.Set("group_id", strconv.Itoa(r.GroupID))
	return
}
func (c *GroupClient) UsersGroupsPreview(userID, groupID int) (res []models2.GroupPreview, err error) {
	const endpoind = "/internal/group/list"
	req := UsersGroupsRequest{
		Uid:     userID,
		GroupID: groupID,
	}

	endpoint := service.ServiceEndpoint{
		Service:  c,
		Endpoint: endpoind,
		Method:   "GET",
		//ContentType: "application/json",
	}
	err = endpoint.Send(&req, &res)
	return
}
