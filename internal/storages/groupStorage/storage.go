package groupStorage

import (
	"database/sql"
	"fmt"
	"github.com/Solar-2020/Group-Backend/internal/models"
	"github.com/lib/pq"
)

const (
	queryReturningID = "RETURNING id;"
	userGroupsTable = "users_groups"
	groupLinksTable = "group_links"
	pgErrorUniqueConstraint = "23505"
)

type Storage interface {
	InsertGroup(group models.Group) (groupReturn models.Group, err error)

	InsertUser(groupID, userID, roleID int) (err error)
	EditUserRole(groupID, userID, roleID int) (resultRole int, err error)
	RemoveUser(groupID, userID int) (err error)

	UpdateGroup(group models.Group) (groupReturn models.Group, err error)

	UpdateGroupStatus(groupID, statusID int) (group models.Group, err error)

	SelectGroupByID(groupID int) (group models.Group, err error)

	SelectGroupRole(groupID, userID int) (roleID int, err error)

	SelectGroupsByUserID(userID int) (group []models.GroupPreview, err error)
	HashToGroupID(line string) (groupID  int, err error)
	RemoveLinkToGroup(groupID int, link string) (err error)
	ListShortLinksToGroup(groupID int) (res []models.GroupInviteLink, err error)
	AddShortLinkToGroup(groupID int, link string, author int) (err error)
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
	if pgErr, ok := err.(*pq.Error); ok {
		if pgErr.Code == pgErrorUniqueConstraint {
			err = fmt.Errorf("exists")
		}
	}

	return
}

func (s *storage) EditUserRole(groupID, userID, roleID int) (resultRole int, err error) {
	const sqlQuery = `
	UPDATE %s SET role_id=$1 WHERE group_id=$2 AND user_id=$3
	RETURNING role_id`

	row := s.db.QueryRow(
		fmt.Sprintf(sqlQuery, userGroupsTable),
		roleID, groupID, userID)
	if row == nil {
		err = fmt.Errorf("nil row")
		return
	}
	err = row.Scan(&resultRole)
	return
}

func (s *storage) RemoveUser(groupID, userID int) (err error) {
	const sqlQuery = `
	DELETE FROM %s WHERE group_id=$1 AND user_id=$2`
	res, err := s.db.Exec(fmt.Sprintf(sqlQuery, userGroupsTable), groupID, userID)
	if err != nil {
		return
	}

	if c, err2 := res.RowsAffected(); err2 == nil && c < 1 {
		err = fmt.Errorf("removed nothing")
	}
	return
}


func (s *storage) SelectGroupsByUserID(userID int) (groups []models.GroupPreview, err error) {
	const sqlQuery = `
	SELECT g.id,
		   g.title,
		   g.description,
		   g.url,
		   g.avatar_url,
		   r.id,
		   r.title,
		   g.status_id
	FROM groups AS g
			 JOIN users_groups AS ug ON g.id = ug.group_id
			 JOIN roles AS r ON ug.role_id = r.id
	WHERE ug.user_id = $1`

	groups = make([]models.GroupPreview, 0)
	rows, err := s.db.Query(sqlQuery, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tempGroup models.GroupPreview
		err = rows.Scan(&tempGroup.ID, &tempGroup.Title, &tempGroup.Description, &tempGroup.URL,
			&tempGroup.AvatarURL, &tempGroup.UserRoleID, &tempGroup.UserRole, &tempGroup.Status)
		if err != nil {
			return
		}
		tempGroup.UserID = userID
		groups = append(groups, tempGroup)
	}
	return
}

func (s *storage) HashToGroupID(line string) (groupID  int, err error) {
	const sqlTemplate = `SELECT group_id from %s WHERE link=$1`
	query := fmt.Sprintf(sqlTemplate, groupLinksTable)

	row := s.db.QueryRow(query, line)
	if row == nil {
		err = fmt.Errorf("nil row")
		return
	}
	err = row.Scan(&groupID)
	return
}

func (s *storage) AddShortLinkToGroup(groupID int, link string, author int) (err error) {
	const sqlTemplate = `INSERT INTO %s (group_id, link, author) VALUES ($1, $2, $3)`
	query := fmt.Sprintf(sqlTemplate, groupLinksTable)

	res, err := s.db.Exec(query, groupID, link, author)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return
	}
	if c < 1 {
		err = fmt.Errorf("not added")
	}
	return
}

func (s *storage) ListShortLinksToGroup(groupID int) (res []models.GroupInviteLink, err error) {
	const sqlTemplate = `SELECT link, created, author from %s WHERE group_id=$1`
	query := fmt.Sprintf(sqlTemplate, groupLinksTable)

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		return
	}

	for rows.Next(){
		link := models.GroupInviteLink{}
		err = rows.Scan(&link.Link, &link.Added, &link.Author)
		if err != nil {
			return
		}
		res = append(res, link)
	}
	return
}

func (s *storage) RemoveLinkToGroup(groupID int, link string) (err error) {
	const sqlTemplate = `DELETE FROM %s WHERE group_id=$1 AND link=$2`
	query := fmt.Sprintf(sqlTemplate, groupLinksTable)

	res, err := s.db.Exec(query, groupID, link)
	if err != nil {
		return
	}
	c, err := res.RowsAffected()
	if err != nil {
		return
	}
	if c < 1 {
		err = fmt.Errorf("not removed")
	}
	return
}