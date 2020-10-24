package group

import (
	"github.com/Solar-2020/Group-Backend/internal/models"
)

type groupStorage interface {
	InsertGroup(group models.Group) (groupReturn models.Group, err error)
	UpdateGroup(group models.Group) (groupReturn models.Group, err error)
	UpdateGroupStatus(groupID, statusID int) (group models.Group, err error)
	SelectGroupByID(groupID int) (group models.Group, err error)
	SelectGroupRole(groupID, userID int) (roleID int, err error)
	SelectGroupsByUserID(userID int) (group []models.GroupPreview, err error)

	InsertUser(groupID, userID, roleID int) (err error)
	EditUserRole(groupID, userID, roleID int) (resultRole int, err error)
	RemoveUser(groupID, userID int) (err error)

	HashToGroupID(line string) (groupID  int, err error)
	RemoveLinkToGroup(groupID int, link string) (err error)
	ListShortLinksToGroup(groupID int) (res []models.GroupInviteLink, err error)
	AddShortLinkToGroup(groupID int, link string, author int) (err error)
}
