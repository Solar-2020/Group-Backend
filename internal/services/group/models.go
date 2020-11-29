package group

import (
	account "github.com/Solar-2020/Account-Backend/pkg/models"
	"github.com/Solar-2020/Group-Backend/internal/models"
	group "github.com/Solar-2020/Group-Backend/pkg/models"
	"github.com/pkg/errors"
)

var (
	ErrorInternalServer = errors.New("Внутрення ошибка сервиса, повторите попытку позже")
)

type groupStorage interface {
	InsertGroup(group group.Group) (groupReturn group.Group, err error)
	UpdateGroup(group group.Group) (groupReturn group.Group, err error)
	UpdateGroupStatus(groupID, statusID int) (group group.Group, err error)
	SelectGroupByID(groupID int) (group group.Group, err error)
	SelectGroupRole(groupID, userID int) (role group.UserRole, err error)
	SelectPermission(actionID, roleID int) (permission models.Permission, err error)
	SelectGroupsByUserID(userID int, groupID int) (group []group.GroupPreview, err error)

	SelectUsersByGroupID(groupID int) (users []group.UserRole, err error)
	InsertUser(groupID, userID, roleID int) (err error)
	EditUserRole(groupID, userID, roleID int) (resultRole int, err error)
	RemoveUser(groupID, userID int) (err error)

	HashToGroupID(line string) (groupID int, err error)
	RemoveLinkToGroup(groupID int, link string) (err error)
	ListShortLinksToGroup(groupID int) (res []group.GroupInviteLink, err error)
	AddShortLinkToGroup(groupID int, link string, author int) (err error)
}

type accountClient interface {
	GetUserByUid(userID int) (user account.User, err error)
	GetUserByEmail(email string) (user account.User, err error)
	CreateUserAdvance(request account.UserAdvance) (userID int, err error)
}

type errorWorker interface {
	NewError(httpCode int, responseError error, fullError error) (err error)
}
