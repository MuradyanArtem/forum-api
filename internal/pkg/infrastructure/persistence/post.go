package persistence

import (
	"database/sql"
	"forum-api/internal/pkg/domain/models"
	"time"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type PostDB struct {
	db *pgx.ConnPool
}

func NewPost(db *pgx.ConnPool) *PostDB {
	return &PostDB{db: db}
}

func (postDB *PostDB) InsertInto(posts []*models.Post, thread *models.Thread) error {
	tx, err := postDB.db.Begin()
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

	created := sql.NullTime{}
	for i, _ := range posts {
		posts[i].Thread = thread.Id
		posts[i].Forum = thread.Forum

		err = tx.QueryRow(
			"post_insert_into",
			posts[i].Author,
			posts[i].Message,
			posts[i].Parent,
			posts[i].Thread,
			posts[i].Forum).Scan(&posts[i].Id, &created)

		if err != nil {
			return err
		}

		if created.Valid {
			posts[i].Created = created.Time.Format(time.RFC3339Nano)
		}
	}

	if len(posts) > 0 {
		_, err := tx.Exec("forum_posts_update", len(posts), posts[0].Forum)
		if err != nil {
			logrus.Error("Error while update post count: " + err.Error())
		}
	}

	return nil
}

func (postDB *PostDB) GetById(post *models.Post) error {
	tx, err := postDB.db.Begin()
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

	created := sql.NullTime{}
	err = tx.QueryRow("post_get_by_id", post.Id).
		Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Path)
	if err != nil {
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (postDB *PostDB) Update(post *models.Post) error {
	tx, err := postDB.db.Begin()
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

	created := sql.NullTime{}
	err = tx.QueryRow("post_update", post.Message, post.Id).
		Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	if err != nil {
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (postDB *PostDB) Prepare() error {
	_, err := postDB.db.Prepare("post_insert_into",
		"INSERT INTO post (usr, message,  parent, thread, forum, created) "+
			"VALUES ($1, $2, $3, $4, $5, current_timestamp) "+
			"RETURNING id, created",
	)
	if err != nil {
		return err
	}

	_, err = postDB.db.Prepare("post_get_by_id",
		"SELECT p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread, p.path "+
			"FROM post p "+
			"WHERE p.id = $1",
	)
	if err != nil {
		return err
	}

	_, err = postDB.db.Prepare("post_update",
		"UPDATE post SET message = $1, isEdited = true "+
			"WHERE id = $2 "+
			"RETURNING usr, created, forum, isEdited, message, parent, thread",
	)
	if err != nil {
		return err
	}

	_, err = postDB.db.Prepare("forum_posts_update",
		"UPDATE forum  SET posts = (posts + $1) "+
			"where slug = $2",
	)
	if err != nil {
		return err
	}

	return nil
}
