package groupStorage

import (
	"database/sql"
	"github.com/Solar-2020/Group-Backend/internal/models"
)

const (
	queryReturningID = "RETURNING id;"
)

type Storage interface {
	InsertGroup(group models.Group) (groupReturn models.Group, err error)

	InsertUser(groupID, userID, roleID int) (err error)

	UpdateGroup(group models.Group) (groupReturn models.Group, err error)

	UpdateGroupStatus(groupID, statusID int) (group models.Group, err error)

	SelectGroupByID(groupID int) (group models.Group, err error)

	SelectGroupRole(groupID, userID int) (roleID int, err error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) InsertGroup(group models.Group) (groupReturn models.Group, err error) {
	const sqlQuery = `
	INSERT INTO groups(title, description, url, create_by, avatar_url)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, create_at, status_id;`

	err = s.db.QueryRow(sqlQuery, group.Title, group.Description, group.URL, group.CreateBy, group.AvatarURL).Scan(&group.ID, &group.CreatAt, &group.StatusID)
	return group, err
}

func (s *storage) UpdateGroup(group models.Group) (groupReturn models.Group, err error) {
	const sqlQuery = `
	UPDATE groups.groups
	SET title=$1,
		description=$2,
		url=$3,
		avatar_url=$4
	WHERE id = $5
	RETURNING id, title, description, url, create_by, create_at, status_id, avatar_url`

	err = s.db.QueryRow(sqlQuery, group.Title, group.Description, group.URL, group.AvatarURL, group.ID).Scan(&group.ID, &group.Title, &group.Description, &group.URL, &group.CreateBy, &group.CreatAt, &group.StatusID, &group.AvatarURL)
	return group, err
}

func (s *storage) UpdateGroupStatus(groupID, statusID int) (group models.Group, err error) {
	const sqlQuery = `
	UPDATE groups
	SET status_id = $1
	WHERE id = $2
	RETURNING id, title, description, url, create_by, create_at, status_id, avatar_url;`

	err = s.db.QueryRow(sqlQuery, statusID, groupID).Scan(&group.ID, &group.Title, &group.Description, &group.URL, &group.CreateBy, &group.CreatAt, &group.StatusID, &group.AvatarURL)
	return
}

func (s *storage) SelectGroupByID(groupID int) (group models.Group, err error) {
	const sqlQuery = `
	SELECT g.id,
		   g.title,
		   g.description,
		   g.url,
		   g.create_by,
		   g.create_at,
		   g.status_id,
		   g.avatar_url
	FROM groups as g
	WHERE g.id = $1;`

	err = s.db.QueryRow(sqlQuery, groupID).Scan(&group.ID, &group.Title, &group.Description, &group.URL, &group.CreateBy, &group.CreatAt, &group.StatusID, &group.AvatarURL)
	return
}

func (s *storage) SelectGroupRole(groupID, userID int) (roleID int, err error) {
	const sqlQuery = `
	SELECT ug.role_id
	FROM users_groups as ug
	WHERE ug.group_id = $1 AND ug.user_id = $2;`

	err = s.db.QueryRow(sqlQuery, groupID, userID).Scan(&roleID)
	return
}

func (s *storage) InsertUser(groupID, userID, roleID int) (err error) {
	const sqlQuery = `
	INSERT INTO users_groups(group_id, user_id, role_id)
	VALUES ($1, $2, $3);`

	_, err = s.db.Exec(sqlQuery, groupID, userID, roleID)

	return
}
