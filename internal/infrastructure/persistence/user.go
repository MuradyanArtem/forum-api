package persistence

import (
	"database/sql"
	"forum-api/internal/domain/models"

	"github.com/jackc/pgx"
)

type User struct {
	db *pgx.ConnPool
}

func newUser(db *pgx.ConnPool) *User {
	return &User{db: db}
}

func (userDB *User) InsertInto(user *models.User) error {
	tx, err := userDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	var info string
	about := &sql.NullString{}

	if user.About != "" {
		about.String = user.About
		about.Valid = true
	}

	err = tx.QueryRow("user_insert", user.Email, user.Fullname, user.Nickname, about).Scan(&info)
	if err != nil {
		return err
	}

	return nil
}

func (userDB *User) GetByNickname(user *models.User) error {
	tx, err := userDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	row := tx.QueryRow("user_get_by_nickname", user.Nickname)

	fullname := &sql.NullString{}
	about := &sql.NullString{}

	if err := row.Scan(&user.Email, fullname, &user.Nickname, about); err != nil {
		return err
	}

	if about.Valid {
		user.About = about.String
	}

	if fullname.Valid {
		user.Fullname = fullname.String
	}

	return nil
}

func (userDB *User) GetByNicknameOrEmail(user *models.User) ([]models.User, error) {
	tx, err := userDB.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	rows, err := tx.Query("user_get_by_nickname_or_email", user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)
	for rows.Next() {
		fullname := &sql.NullString{}
		about := &sql.NullString{}

		if err := rows.Scan(&user.Email, fullname, &user.Nickname, about); err != nil {
			return nil, err
		}

		if about.Valid {
			user.About = about.String
		}

		if fullname.Valid {
			user.Fullname = fullname.String
		}

		users = append(users, *user)
	}
	rows.Close()

	return users, nil
}

func (userDB *User) Update(user *models.User) error {
	tx, err := userDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec("user_update", user.Email, user.Nickname, user.Fullname, user.About)
	if err != nil {
		return err
	}

	return nil
}

func (userDB *User) DeleteAll() error {
	tx, err := userDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	_, err = userDB.db.Exec("DELETE FROM usr")
	if err != nil {
		return err
	}

	return nil
}

func (userDB *User) GetStatus(s *models.Status) error {
	tx, err := userDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	rows, err := tx.Query(
		"SELECT count(*) FROM forum " +
			"UNION ALL " +
			"SELECT count(*) " +
			"FROM post " +
			"UNION ALL " +
			"SELECT count(*) FROM thread " +
			"UNION ALL " +
			"SELECT count(*) FROM usr",
	)
	if err != nil {
		return err
	}

	i := 0
	for rows.Next() {
		var err error
		switch i {
		case 0:
			err = rows.Scan(&s.Forum)
		case 1:
			err = rows.Scan(&s.Post)
		case 2:
			err = rows.Scan(&s.Thread)
		case 3:
			err = rows.Scan(&s.User)
		}
		if err != nil {
			rows.Close()
			return err
		}
		i++
	}
	rows.Close()

	return nil
}

func (userDB *User) Prepare() error {
	_, err := userDB.db.Prepare("user_insert",
		"INSERT INTO usr (email, fullname, nickname, about) "+
			"VALUES ($1, $2, $3, $4) "+
			"ON CONFLICT DO NOTHING "+
			"RETURNING email",
	)
	if err != nil {
		return err
	}

	_, err = userDB.db.Prepare("user_get_by_nickname",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM usr u "+
			"WHERE nickname = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = userDB.db.Prepare("user_get_by_nickname_or_email",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM usr u "+
			"WHERE nickname = $1 OR email = $2",
	)
	if err != nil {
		return err
	}

	_, err = userDB.db.Prepare("user_update",
		"UPDATE usr SET "+
			"email = $1, "+
			"nickname = $2, "+
			"fullname = $3, "+
			"about = $4 "+
			"WHERE nickname = $2 RETURNING email",
	)
	if err != nil {
		return err
	}

	return nil
}