package group

import (
	"github.com/Solar-2020/Group-Backend/internal/models"
)

type groupStorage interface {
	InsertGroup(group models.Group) (groupReturn models.Group, err error)

	InsertUser(groupID, userID, roleID int) (err error)

	UpdateGroup(group models.Group) (groupReturn models.Group, err error)

	UpdateGroupStatus(groupID, statusID int) (group models.Group, err error)

	SelectGroupByID(groupID int) (group models.Group, err error)

	SelectGroupRole(groupID, userID int) (roleID int, err error)
}
