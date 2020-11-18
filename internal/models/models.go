package models

import "github.com/pkg/errors"

var (
	ErrorNoMembership = errors.New("Вы не состоите в данной группе")
	ErrorNoPermission = errors.New("У Вас не достаточно прав")
)

type Permission struct {
	RoleID   int
	ActionID int
}
