package persistence

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"

	"github.com/jackc/pgx"
)

type User struct {
	db *pgx.ConnPool
}

func newUser(db *pgx.ConnPool) *User {
	return &User{db: db}
}

func (u *User) Insert(user *models.User) error {
	if _, err := u.db.Exec("INSERT INTO users (nickname, "+
		"email, about, fullname) "+
		"VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Email,
		user.About,
		user.Fullname); err != nil {
		return infrastructure.ErrConflict
	}

	return nil
}

func (u *User) Update(user *models.User) error {
	query := "UPDATE users SET   "
	if user.About != "" {
		query += ` about = '` + user.About + "' , "
	}
	if user.Fullname != "" {
		query += " fullname = '" + user.Fullname + "' , "
	}
	if user.Email != "" {
		query += " email = '" + user.Email + "' , "
	}
	query = query[:len(query)-2]
	query += "WHERE nickname = '" + user.Nickname + "' RETURNING about, email, fullname, nickname "
	if user.About == "" && user.Email == "" && user.Fullname == "" {
		query = "SELECT about, email, fullname, nickname FROM users WHERE nickname = '" + user.Nickname + "' "
	}
	err := u.db.QueryRow(query).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return infrastructure.ErrNotExists
		default:
			return infrastructure.ErrConflict
		}
	}
	return nil
}

func (u *User) SelectByNickname(nickname string) (models.User, error) {
	user := models.User{}
	err := u.db.QueryRow(
		"SELECT users.email, users.nickname, users.about, users.fullname "+
			"FROM users "+
			"WHERE nickname = $1",
		nickname).Scan(&user.Email, &user.Nickname, &user.About, &user.Fullname)
	switch err {
	case pgx.ErrNoRows:
		return models.User{}, infrastructure.ErrNotExists
	default:
		return user, err
	}
}

func (u *User) SelectByEmailOrNickname(nickname string, email string) (models.UserSlice, error) {
	rows, err := u.db.Query(
		"SELECT users.email, users.nickname, users.about, users.fullname "+
			"FROM users "+
			"WHERE nickname = $1 OR email = $2",
		nickname, email)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.Email, &user.Nickname, &user.About, &user.Fullname)
		users = append(users, user)
	}
	return users, nil
}

func (u *User) SelectNicknameWithCase(nickname string) (string, error) {
	var result string
	err := u.db.QueryRow("SELECT nickname FROM users "+
		"WHERE nickname = $1", nickname).
		Scan(&result)

	switch err {
	case pgx.ErrNoRows:
		return "", infrastructure.ErrNotExists
	case nil:
		return result, nil
	default:
		return "", err
	}
}

func (u *User) DeleteAll() error {
	_, err := u.db.Exec("TRUNCATE  votes, posts, forum_users, threads, forums, users CASCADE")
	return err
}

func (u *User) GetStatus(status *models.Status) error {
	return u.db.QueryRow("SELECT "+
		"(SELECT COUNT(*) FROM forums) as forums_status, "+
		"(SELECT COUNT(*) FROM threads) as threads_status, "+
		"(SELECT COUNT(*) FROM posts) as posts_status, "+
		"(SELECT COUNT(*) FROM users) as users_status").
		Scan(
			&status.Forum,
			&status.Thread,
			&status.Post,
			&status.User)
}
