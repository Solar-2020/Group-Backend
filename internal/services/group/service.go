package group

import (
	"errors"
	"github.com/Solar-2020/Group-Backend/internal/models"
)

type Service interface {
	Create(request models.Group) (response models.Group, err error)
	Update(group models.Group, userID int) (response models.Group, err error)
	Delete(groupID, userID int) (response models.Group, err error)
	Get(groupID, userID int) (response models.Group, err error)
	GetList(userID int) (response []models.GroupPreview, err error)
	Invite(request models.InviteUserRequest) (response models.InviteUserResponse, err error)
	ChangeRole(request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error)
	ExpelUser(request models.ExpelUserRequest) (response models.ExpelUserResponse, err error)
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

func (s *service) Invite(request models.InviteUserRequest) (response models.InviteUserResponse, err error) {
	// TODO: userEmail -> userID
	userID := 5
	err = s.groupStorage.InsertUser(request.Group, userID, int(request.Role))
	response = models.InviteUserResponse{
		Group: request.Group, User: request.User, Role: request.Role,
	}
	return
}

func (s *service) ChangeRole(request models.ChangeRoleRequest) (response models.ChangeRoleResponse, err error) {
	// TODO: userEmail -> userID
	userID := 5
	newRole, err := s.groupStorage.EditUserRole(request.Group, userID, int(request.Role))
	response.Role = models.MemberRole(newRole)
	return
}

func (s *service) ExpelUser(request models.ExpelUserRequest) (response models.ExpelUserResponse, err error) {
	// TODO: userEmail -> userID
	userID := 5
	err = s.groupStorage.RemoveUser(int(request.Group), userID)
	response.User = request.User
	return
}