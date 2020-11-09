package group

import (
	models2 "github.com/Solar-2020/Group-Backend/pkg/models"
)

type groupStorage interface {
	InsertGroup(group models2.Group) (groupReturn models2.Group, err error)
	UpdateGroup(group models2.Group) (groupReturn models2.Group, err error)
	UpdateGroupStatus(groupID, statusID int) (group models2.Group, err error)
	SelectGroupByID(groupID int) (group models2.Group, err error)
	SelectGroupRole(groupID, userID int) (role models2.UserRole, err error)
	SelectGroupsByUserID(userID int, groupID int) (group []models2.GroupPreview, err error)

	InsertUser(groupID, userID, roleID int) (err error)
	EditUserRole(groupID, userID, roleID int) (resultRole int, err error)
	RemoveUser(groupID, userID int) (err error)

	HashToGroupID(line string) (groupID  int, err error)
	RemoveLinkToGroup(groupID int, link string) (err error)
	ListShortLinksToGroup(groupID int) (res []models2.GroupInviteLink, err error)
	AddShortLinkToGroup(groupID int, link string, author int) (err error)
}
