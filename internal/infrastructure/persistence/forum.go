package persistence

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"

	"github.com/jackc/pgx"
)

type Forum struct {
	db *pgx.ConnPool
}

func newForum(db *pgx.ConnPool) *Forum {
	return &Forum{db: db}
}

func (f *Forum) Insert(forum *models.Forum) error {
	if err := f.db.QueryRow("INSERT INTO forums (slug, title, nickname) " +
		"VALUES ($1, $2, $3) "+
		"RETURNING nickname", &forum.Slug, &forum.Title, &forum.User).Scan(&forum.User); err != nil {
		switch infrastructure.ErrCode(err) {
		case infrastructure.PgErrUniqueViolation:
			return infrastructure.ErrConflict
		default:
			return infrastructure.ErrNotExists
		}
	}
	return nil
}

func (f *Forum) SelectBySlug(slug string) (*models.Forum, error) {
	forum := &models.Forum{}
	if err := f.db.QueryRow("SELECT forums.title, forums.nickname, forums.slug, forums.posts, forums.threads "+
	"FROM forums "+
	"WHERE slug = $1 ", slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
		return nil, err
	}
	return forum, nil
}

func (f *Forum) GetUsersByForum(slug string, desc bool, since string, limit int) (models.UserSlice, error) {
	users := models.UserSlice{}
	query := "SELECT users.about, users.email, users.fullname, users.nickname " +
		"FROM forum_users " +
		"JOIN users on users.nickname = forum_users.author " +
		"WHERE slug = $1 "
	rows := &pgx.Rows{}
	var err error
	if limit > 0 && since != "" {
		if desc {
			query += "AND lower(users.nickname) < lower($2::text) ORDER BY users.nickname  DESC LIMIT $3"
		} else {
			query += "AND lower(users.nickname)  > lower($2::text) ORDER BY users.nickname  ASC LIMIT $3"
		}
		rows, err = f.db.Query(query, &slug, &since, &limit)
	} else {
		if limit > 0 {
			if desc {
				query += "ORDER BY users.nickname DESC LIMIT $2"
			} else {
				query += "ORDER BY users.nickname ASC LIMIT $2"
			}
			rows, err = f.db.Query(query, &slug, &limit)
		} else if since != "" {
			if desc {
				query += "AND lower(users.nickname) < lower($2::text) ORDER BY users.nickname DESC "
			} else {
				query += "AND lower(users.nickname) > lower($2::text) ORDER BY users.nickname ASC "
			}
			rows, err = f.db.Query(query, &slug, &since)
		} else {
			rows, err = f.db.Query(query, &slug)
		}
	}
	if err != nil {
		return users, nil
	}
	for rows.Next() {
		user := models.User{}

		err := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (f *Forum) SelectForumWithCase(slug string) (string, error) {
	var res string
	err := f.db.QueryRow("SELECT slug FROM forums WHERE slug = $1", &slug).Scan(&res)
	return res, err
}
