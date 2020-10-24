package group

import (
	"errors"
	"fmt"
	"github.com/Solar-2020/GoUtils/context"
	"github.com/Solar-2020/Group-Backend/internal"
	"github.com/Solar-2020/Group-Backend/internal/models"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Service interface {
	Create(ctx context.Context, request models.CreateRequest) (response models.CreateResponse, err error)
	Update(ctx context.Context, request models.UpdateRequest) (response models.UpdateResponse, err error)
	Delete(ctx context.Context, request models.DeleteRequest) (response models.DeleteResponse, err error)
	Get(ctx context.Context, request models.GetRequest) (response models.GetResponse, err error)
	GetList(ctx context.Context, request models.GetListRequest) (response models.GetListResponse, err error)

	Invite(ctx context.Context, request models.InviteUserRequest) (response models.InviteUserResponse, err error)
	ChangeRole(ctx context.Context, request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error)
	ExpelUser(ctx  context.Context, request models.ExpelUserRequest) (response models.ExpelUserResponse, err error)

	ResolveGroup(ctx context.Context, request models.ResolveInviteLinkRequest) (response models.ResolveInviteLinkResponse, err error)
	AddGroupInviteLink(ctx context.Context, request models.AddInviteLinkRequest) (response models.AddInviteLinkResponse, err error)
	RemoveGroupInviteLink(ctx context.Context, request models.RemoveInviteLinkRequest) (response models.RemoveInviteLinkRsponse, err error)
	ListGroupInviteLink(ctx context.Context, request models.ListInviteLinkRequest) (response models.ListInviteLinkResponse, err error)

	CheckPermission(ctx context.Context, group models.Group, action models.GroupAction) error
}

var (
	inviteHashParse = regexp.MustCompile(`http(?:s)?:\/\/.*\/(\w+)(?:\/)?`)
)

type service struct {
	groupStorage groupStorage
}

func NewService(groupStorage groupStorage) Service {
	return &service{
		groupStorage: groupStorage,
	}
}

func (s *service) Create(ctx context.Context, request models.CreateRequest) (response models.CreateResponse, err error) {
	if request.CreateBy == 0 {
		request.CreateBy = ctx.UserID
	}
	err = s.validateGroup(request.Group)
	if err != nil {
		return
	}

	err = s.checkUnique(request.Group)
	if err != nil {
		return
	}

	response.Group, err = s.groupStorage.InsertGroup(request.Group)
	if err != nil {
		return
	}

	err = s.groupStorage.InsertUser(response.ID, response.CreateBy, 1)
	return
}

func (s *service) validateGroup(group models.Group) (err error) {
	if len(group.URL) < 3 || len(group.URL) > 20 {
		return errors.New("Недопустимая длина ссылки")
	}

	if len(group.Description) > 500 {
		return errors.New("Слишком длинное описание группы")
	}

	if len(group.Title) > 100 {
		return errors.New("Слишком длинное название")
	}

	return
}

func (s *service) checkUnique(group models.Group) (err error) {

	return
}

func (s *service) checkAdminPermission(groupID, userID int) (err error) {
	roleID, err := s.groupStorage.SelectGroupRole(groupID, userID)
	if err != nil {
		return errors.New("У Вас недостаточно прав для совершения данной операции")
	}

	if !(roleID == 1 || roleID == 2) {
		return errors.New("У Вас недостаточно прав для совершения данной операции")
	}

	return
}

func (s *service) checkUserPermission(groupID, userID int) (err error) {
	roleID, err := s.groupStorage.SelectGroupRole(groupID, userID)
	if err != nil {
		return errors.New("У Вас недостаточно прав для совершения данной операции")
	}

	if !(roleID == 1 || roleID == 2 || roleID == 3) {
		return errors.New("У Вас недостаточно прав для совершения данной операции")
	}

	return
}

func (s *service) Update(ctx context.Context, request models.UpdateRequest) (response models.UpdateResponse, err error) {
	if request.UserID == 0 {
		request.UserID = ctx.UserID
	}
	err = s.checkAdminPermission(request.Group.ID, request.UserID)
	if err != nil {
		return
	}

	err = s.validateGroup(request.Group)
	if err != nil {
		return
	}

	err = s.checkUnique(request.Group)
	if err != nil {
		return
	}

	response.Group, err = s.groupStorage.UpdateGroup(request.Group)

	return
}

func (s *service) Delete(ctx context.Context, request models.DeleteRequest) (response models.DeleteResponse, err error) {
	if request.UserID == 0 {
		request.UserID = ctx.UserID
	}
	err = s.checkAdminPermission(request.GroupID, request.UserID)
	if err != nil {
		return
	}

	response.Group, err = s.groupStorage.UpdateGroupStatus(request.GroupID, 2)

	return
}

func (s *service) Get(ctx context.Context, request models.GetRequest) (response models.GetResponse, err error) {
	if request.UserID == 0 {
		request.UserID = ctx.UserID
	}
	err = s.checkUserPermission(request.GroupID, request.UserID)
	if err != nil {
		return
	}

	response.Group, err = s.groupStorage.SelectGroupByID(request.GroupID)

	return
}

func (s *service) GetList(ctx context.Context, request models.GetListRequest) (response models.GetListResponse, err error) {
	if request.UserID == 0 {
		request.UserID = ctx.UserID
	}
	response.Groups, err = s.groupStorage.SelectGroupsByUserID(request.UserID)
	response.UserID = request.UserID
	return
}

func (s *service) Invite(ctx context.Context, request models.InviteUserRequest) (response models.InviteUserResponse, err error) {
	// TODO: userEmail -> userID
	for i := range request.User {
		err_ := s.groupStorage.InsertUser(request.Group, request.UserID[i], int(request.Role))
		if err_ != nil {
			if err == nil {
				err = fmt.Errorf("")
			}
			err = fmt.Errorf("%s; %s", err, fmt.Sprintf("[%d]: %s", i, err_))
		}
	}
	response = models.InviteUserResponse{
		Group: request.Group, User: request.User, Role: request.Role,
	}
	return
}

func (s *service) ChangeRole(ctx context.Context, request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error) {
	// TODO: userEmail -> userID
	userID := request.UserID
	newRole, err := s.groupStorage.EditUserRole(request.Group, userID, int(request.Role))
	response.Role = models.MemberRole(newRole)
	return
}

func (s *service) ExpelUser(ctx context.Context, request models.ExpelUserRequest) (response models.ExpelUserResponse, err error) {
	// TODO: userEmail -> userID
	userID := request.UserID
	err = s.groupStorage.RemoveUser(int(request.Group), userID)
	response.User = request.User
	return
}


func (s *service) CheckPermission(ctx context.Context, groupRequest models.Group, action models.GroupAction) error {
	denied := fmt.Errorf("denied")
	if action == models.ActionCreate {
		if groupRequest.CreateBy != 0 && ctx.UserID != groupRequest.CreateBy {
			return denied
		}
	}
	switch action {
	case models.ActionGet:
		return s.checkUserPermission(groupRequest.ID, ctx.UserID)
	case models.ActionCreate:
		if groupRequest.CreateBy != 0 && ctx.UserID != groupRequest.CreateBy {
			return denied
		}
	case models.ActionEdit, models.ActionEditRole, models.ActionInvite, models.ActionExpel, models.ActionRemove:
		return s.checkAdminPermission(groupRequest.ID, ctx.UserID)
	}
	return denied
}

func (s *service) AddGroupInviteLink(ctx context.Context, request models.AddInviteLinkRequest) (response models.AddInviteLinkResponse, err error) {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"_")
	length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	line := b.String()
	err = s.groupStorage.AddShortLinkToGroup(request.Group, line, ctx.Session.UserID)
	if err != nil {
		return
	}
	response.Group = request.Group
	response.Link = s.getLinkFromHash(line)
	return
}
func (s *service)  RemoveGroupInviteLink(ctx context.Context, request models.RemoveInviteLinkRequest) (response models.RemoveInviteLinkRsponse, err error) {
	errs := make([]string, 0)
	removedLinks := make([]string, len(request.Links))
	for i, item := range request.Links {
		linkHash, _ := s.getHashFromLink(item)
		err = s.groupStorage.RemoveLinkToGroup(request.Group, linkHash)
		if err == nil {
			removedLinks = append(removedLinks, item)
		} else {
			errs = append(errs, fmt.Sprintf("[%d]: %s", i, item))
		}
	}
	if len(errs) != 0 {
		err = fmt.Errorf(strings.Join(errs, ",\n"))
	}
	response.Links = removedLinks
	response.Group = request.Group
	return
}
func (s *service)  ListGroupInviteLink(ctx context.Context, request models.ListInviteLinkRequest) (response models.ListInviteLinkResponse, err error) {
	response.Links, err = s.groupStorage.ListShortLinksToGroup(request.Group)
	for i, elem := range response.Links {
		response.Links[i].Link = s.getLinkFromHash(elem.Link)
	}
	response.Group = request.Group
	return
}

func (s *service) ResolveGroup(ctx  context.Context, request models.ResolveInviteLinkRequest) (response models.ResolveInviteLinkResponse, err error) {
	linkHash, _ := s.getHashFromLink(request.Link)
	id, err := s.groupStorage.HashToGroupID(linkHash)
	temp, err := s.groupStorage.SelectGroupByID(id)
	response.Group = temp.ID
	return
}

func (s *service) getHashFromLink(src string) (res string, err error) {
	parseRes := inviteHashParse.FindStringSubmatch(src)
	res = src
	if parseRes != nil {
		res = parseRes[1]
	} else {
		err = fmt.Errorf("bad link")
	}
	return
}

func (s *service) getLinkFromHash(src string) string {
	return fmt.Sprintf("%s/%s", internal.Config.InviteLinkPrefix, src)
}