package group

import (
	"errors"
	"fmt"
	accountApi "github.com/Solar-2020/Account-Backend/pkg/api"
	"github.com/Solar-2020/Group-Backend/internal"
	"github.com/Solar-2020/Group-Backend/internal/clients/account"
	"github.com/Solar-2020/Group-Backend/internal/models"
	models2 "github.com/Solar-2020/Group-Backend/pkg/models"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Service interface {
	Create(request models2.Group) (response models2.Group, err error)
	Update(request models2.Group, userID int) (response models2.Group, err error)
	Delete(groupID, userID int) (response models2.Group, err error)
	Get(groupID, userID int) (response models2.Group, err error)
	GetList(groupID, userID int) (response []models2.GroupPreview, err error)

	InternalGetList(groupID, userID int) (response []models2.GroupPreview, err error)

	Invite(request models.InviteUserRequest) (response models.InviteUserResponse, err error)
	ChangeRole(request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error)
	ExpelUser(request models.ExpelUserRequest) (response models.ExpelUserResponse, err error)

	ResolveGroup(request models.ResolveInviteLinkRequest) (response models.ResolveInviteLinkResponse, err error)
	AddGroupInviteLink(request models.AddInviteLinkRequest, userID int) (response models.AddInviteLinkResponse, err error)
	RemoveGroupInviteLink(request models.RemoveInviteLinkRequest) (response models.RemoveInviteLinkRsponse, err error)
	ListGroupInviteLink(request models.ListInviteLinkRequest) (response models.ListInviteLinkResponse, err error)

	CheckPermission(group models2.Group, action models2.GroupAction, userID int) error

	GetUserRole(groupID, userID int) (role models.UserRole, err error)
}

var (
	inviteHashParse = regexp.MustCompile(`http(?:s)?:\/\/.*\/(\w+)(?:\/)?`)
)

type service struct {
	groupStorage  groupStorage
	accountClient account.Client
}

func NewService(groupStorage groupStorage) Service {
	return &service{
		groupStorage: groupStorage,
	}
}

func (s *service) Create(request models2.Group) (response models2.Group, err error) {
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

func (s *service) Update(request models2.Group, userID int) (response models2.Group, err error) {
	err = s.checkAdminPermission(request.ID, userID)
	if err != nil {
		return
	}

	err = s.validateGroup(request)
	if err != nil {
		return
	}

	err = s.checkUnique(request)
	if err != nil {
		return
	}

	response, err = s.groupStorage.UpdateGroup(request)

	return
}

func (s *service) Delete(groupID, userID int) (response models2.Group, err error) {
	err = s.checkAdminPermission(groupID, userID)
	if err != nil {
		return
	}

	response, err = s.groupStorage.UpdateGroupStatus(groupID, 2)

	return
}

func (s *service) Get(groupID, userID int) (response models2.Group, err error) {
	err = s.checkUserPermission(groupID, userID)
	if err != nil {
		return
	}

	response, err = s.groupStorage.SelectGroupByID(groupID)

	return
}

func (s *service) GetList(groupID, userID int) (response []models2.GroupPreview, err error) {
	response, err = s.groupStorage.SelectGroupsByUserID(userID, groupID)
	return
}

func (s *service) InternalGetList(groupID, userID int) (response []models2.GroupPreview, err error) {
	response, err = s.groupStorage.SelectGroupsByUserID(userID, groupID)
	return
}

func (s *service) Invite(request models.InviteUserRequest) (response models.InviteUserResponse, err error) {
	// Можно передавать смешанные списки по UserID и Email. Собираем единый.
	userIds := func() map[int]bool {
		m := make(map[int]bool)
		for _, id := range request.UserID {
			m[id] = true
		}
		return m
	}()
	for _, email := range request.User {
		uid, err := s.emailToUid(email)
		if err != nil {
			return response, errors.New(err.Error() + fmt.Sprintf("cant add user %s", email))
		}
		if _, ok := userIds[uid]; !ok {
			request.UserID = append(request.UserID, uid)
			userIds[uid] = true
		}
	}

	//addedUsers := make([]string, 0, len(request.UserID))
	addedUsersID := make([]int, 0, len(request.UserID))

	for i, userId := range request.UserID {
		err_ := s.groupStorage.InsertUser(request.Group, userId, int(request.Role))
		if err_ != nil {
			if err == nil {
				err = fmt.Errorf("")
			}
			err = fmt.Errorf("%s; %s", err, fmt.Sprintf("[%d]: %s", i, err_))
		} else {
			addedUsersID = append(addedUsersID, userId)
		}
	}
	response = models.InviteUserResponse{
		Group: request.Group, Role: request.Role, UserID: addedUsersID,
	}
	return
}

func (s *service) ChangeRole(request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error) {
	if request.UserID == 0 {
		request.UserID, err = s.emailToUid(request.User)
		if err != nil {
			err = fmt.Errorf("bad user: %s", err)
			return
		}
	}
	newRole, err := s.groupStorage.EditUserRole(request.Group, request.UserID, int(request.Role))
	response.Role = models2.MemberRole(newRole)
	return
}

func (s *service) ExpelUser(request models.ExpelUserRequest) (response models.ExpelUserResponse, err error) {
	if request.UserID == 0 {
		request.UserID, err = s.emailToUid(request.User)
		if err != nil {
			err = fmt.Errorf("bad user: %s", err)
			return
		}
	}
	err = s.groupStorage.RemoveUser(int(request.Group), request.UserID)
	response.User = request.User
	return
}

func (s *service) CheckPermission(group models2.Group, action models2.GroupAction, userID int) error {
	denied := fmt.Errorf("denied")
	if action == models2.ActionCreate {
		if group.CreateBy != 0 && userID != group.CreateBy {
			return denied
		}
	}
	switch action {
	case models2.ActionCreate:
		if group.CreateBy != 0 && userID != group.CreateBy {
			return denied
		}
	case models2.ActionGet:
		return s.checkUserPermission(group.ID, userID)
	case models2.ActionEdit, models2.ActionEditRole, models2.ActionInvite, models2.ActionExpel, models2.ActionRemove:
		return s.checkAdminPermission(group.ID, userID)
	}
	return denied
}

func (s *service) AddGroupInviteLink(request models.AddInviteLinkRequest, userID int) (response models.AddInviteLinkResponse, err error) {
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
	err = s.groupStorage.AddShortLinkToGroup(request.Group, line, userID)
	if err != nil {
		return
	}
	response.Group = request.Group
	response.Link = s.getLinkFromHash(line)
	return
}
func (s *service) RemoveGroupInviteLink(request models.RemoveInviteLinkRequest) (response models.RemoveInviteLinkRsponse, err error) {
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
func (s *service) ListGroupInviteLink(request models.ListInviteLinkRequest) (response models.ListInviteLinkResponse, err error) {
	response.Links, err = s.groupStorage.ListShortLinksToGroup(request.Group)
	for i, elem := range response.Links {
		response.Links[i].Link = s.getLinkFromHash(elem.Link)
	}
	response.Group = request.Group
	return
}

func (s *service) ResolveGroup(request models.ResolveInviteLinkRequest) (response models.ResolveInviteLinkResponse, err error) {
	linkHash, _ := s.getHashFromLink(request.Link)
	id, err := s.groupStorage.HashToGroupID(linkHash)
	temp, err := s.groupStorage.SelectGroupByID(id)
	response.Group = temp.ID
	return
}

func (s *service) GetUserRole(groupID, userID int) (role models.UserRole, err error) {
	role.GroupID = groupID
	role.UserID = userID
	role.RoleID, err = s.groupStorage.SelectGroupRole(groupID, userID)
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

func (s *service) emailToUid(email string) (uid int, err error) {
   client := accountApi.AccountClient{
		   Addr:    internal.Config.AccountServiceAddress,
   }
   user, err := client.GetUserByEmail(email)
   if err != nil {
		   err = fmt.Errorf("bad user: %s", err)
		   return
   }
   uid = user.ID
   return
}
