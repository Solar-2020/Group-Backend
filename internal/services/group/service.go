package group

import (
	"errors"
	"fmt"
	"github.com/Solar-2020/GoUtils/context"
	"github.com/Solar-2020/Group-Backend/internal"
	"github.com/Solar-2020/Group-Backend/internal/models"
	models2 "github.com/Solar-2020/Group-Backend/pkg/models"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Service interface {
	Create(ctx context.Context, request models2.Group) (response models2.Group, err error)
	Update(ctx context.Context, group models2.Group) (response models2.Group, err error)
	Delete(ctx context.Context, groupID int) (response models2.Group, err error)
	Get(ctx context.Context, groupID int) (response models2.Group, err error)
	GetList(ctx context.Context, groupID int) (response []models2.GroupPreview, err error)

	Invite(ctx context.Context, request models.InviteUserRequest) (response models.InviteUserResponse, err error)
	ChangeRole(ctx context.Context, request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error)
	ExpelUser(ctx  context.Context, request models.ExpelUserRequest) (response models.ExpelUserResponse, err error)

	ResolveGroup(ctx context.Context, request models.ResolveInviteLinkRequest) (response models.ResolveInviteLinkResponse, err error)
	AddGroupInviteLink(ctx context.Context, request models.AddInviteLinkRequest) (response models.AddInviteLinkResponse, err error)
	RemoveGroupInviteLink(ctx context.Context, request models.RemoveInviteLinkRequest) (response models.RemoveInviteLinkRsponse, err error)
	ListGroupInviteLink(ctx context.Context, request models.ListInviteLinkRequest) (response models.ListInviteLinkResponse, err error)

	CheckPermission(ctx context.Context, group models2.Group, action models2.GroupAction) error
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

func (s *service) Create(ctx context.Context, request models2.Group) (response models2.Group, err error) {
	request.CreateBy = ctx.Session.Uid
	err = s.validateGroup(request)
	if err != nil {
		return
	}

	err = s.checkUnique(request)
	if err != nil {
		return
	}

	response, err = s.groupStorage.InsertGroup(request)
	if err != nil {
		return
	}

	err = s.groupStorage.InsertUser(response.ID, response.CreateBy, 1)
	return
}

func (s *service) validateGroup(group models2.Group) (err error) {
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

func (s *service) checkUnique(group models2.Group) (err error) {

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

func (s *service) Update(ctx context.Context, group models2.Group) (response models2.Group, err error) {
	err = s.checkAdminPermission(group.ID, ctx.Session.Uid)
	if err != nil {
		return
	}

	err = s.validateGroup(group)
	if err != nil {
		return
	}

	err = s.checkUnique(group)
	if err != nil {
		return
	}

	response, err = s.groupStorage.UpdateGroup(group)

	return
}

func (s *service) Delete(ctx context.Context, groupID int) (response models2.Group, err error) {
	err = s.checkAdminPermission(groupID, ctx.Session.Uid)
	if err != nil {
		return
	}

	response, err = s.groupStorage.UpdateGroupStatus(groupID, 2)

	return
}

func (s *service) Get(ctx context.Context, groupID int) (response models2.Group, err error) {
	err = s.checkUserPermission(groupID, ctx.Session.Uid)
	if err != nil {
		return
	}

	response, err = s.groupStorage.SelectGroupByID(groupID)

	return
}

func (s *service) GetList(ctx context.Context, groupID int) (response []models2.GroupPreview, err error) {
	response, err = s.groupStorage.SelectGroupsByUserID(ctx.Session.Uid, groupID)
	return
}

func (s *service) Invite(ctx context.Context, request models.InviteUserRequest) (response models.InviteUserResponse, err error) {
	// TODO: userEmail -> userID
	addedUsers := make([]string, 0, len(request.UserID))
	addedUsersID := make([]int, 0, len(request.UserID))
	for i, userId := range request.UserID {
		err_ := s.groupStorage.InsertUser(request.Group, userId, int(request.Role))
		if err_ != nil {
			if err == nil {
				err = fmt.Errorf("")
			}
			err = fmt.Errorf("%s; %s", err, fmt.Sprintf("[%d]: %s", i, err_))
		} else {
			addedUsers = append(addedUsers, request.User[i])
			addedUsersID = append(addedUsersID, userId)
		}
	}
	response = models.InviteUserResponse{
		Group: request.Group, User: addedUsers, Role: request.Role, UserID: addedUsersID,
	}
	return
}

func (s *service) ChangeRole(ctx context.Context, request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error) {
	// TODO: userEmail -> userID
	userID := request.UserID
	newRole, err := s.groupStorage.EditUserRole(request.Group, userID, int(request.Role))
	response.Role = models2.MemberRole(newRole)
	return
}

func (s *service) ExpelUser(ctx context.Context, request models.ExpelUserRequest) (response models.ExpelUserResponse, err error) {
	// TODO: userEmail -> userID
	userID := request.UserID
	err = s.groupStorage.RemoveUser(int(request.Group), userID)
	response.User = request.User
	return
}


func (s *service) CheckPermission(ctx context.Context, groupRequest models2.Group, action models2.GroupAction) error {
	denied := fmt.Errorf("denied")
	if action == models2.ActionCreate {
		if groupRequest.CreateBy != 0 && ctx.Session.Uid != groupRequest.CreateBy {
			return denied
		}
	}
	switch action {
	case models2.ActionCreate:
		if groupRequest.CreateBy != 0 && ctx.Session.Uid != groupRequest.CreateBy {
			return denied
		}
	case models2.ActionGet:
		return s.checkUserPermission(groupRequest.ID, ctx.Session.Uid)
	case models2.ActionEdit, models2.ActionEditRole, models2.ActionInvite, models2.ActionExpel, models2.ActionRemove:
		return s.checkAdminPermission(groupRequest.ID, ctx.Session.Uid)
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
	err = s.groupStorage.AddShortLinkToGroup(request.Group, line, ctx.Session.Uid)
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