package group

import (
	"errors"
	"fmt"
	"github.com/Solar-2020/GoUtils/context"
	"github.com/Solar-2020/Group-Backend/internal/models"
)

type Service interface {
	Create(request models.Group) (response models.Group, err error)
	Update(group models.Group, userID int) (response models.Group, err error)
	Delete(groupID, userID int) (response models.Group, err error)
	Get(groupID, userID int) (response models.Group, err error)
	GetList(userID int) (response []models.GroupPreview, err error)
	Invite(ctx context.Context, request models.InviteUserRequest) (response models.InviteUserResponse, err error)
	ChangeRole(ctx context.Context, request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error)
	ExpelUser(ctx  context.Context, request models.ExpelUserRequest) (response models.ExpelUserResponse, err error)

	CheckPermission(ctx context.Context, group models.Group, action models.GroupAction) error
}

type service struct {
	groupStorage groupStorage
}

func NewService(groupStorage groupStorage) Service {
	return &service{
		groupStorage: groupStorage,
	}
}

func (s *service) Create(request models.Group) (response models.Group, err error) {
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

func (s *service) Update(group models.Group, userID int) (response models.Group, err error) {
	err = s.checkAdminPermission(group.ID, userID)
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

func (s *service) Delete(groupID, userID int) (response models.Group, err error) {
	err = s.checkAdminPermission(groupID, userID)
	if err != nil {
		return
	}

	response, err = s.groupStorage.UpdateGroupStatus(groupID, 2)

	return
}

func (s *service) Get(groupID, userID int) (response models.Group, err error) {
	err = s.checkUserPermission(groupID, userID)
	if err != nil {
		return
	}

	response, err = s.groupStorage.SelectGroupByID(groupID)

	return
}

func (s *service) GetList(userID int) (response []models.GroupPreview, err error) {
	response, err = s.groupStorage.SelectGroupsByUserID(userID)

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
			err = fmt.Errorf("%s\n%s", err, fmt.Sprintf("[%d]: %s", i, err_))
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
