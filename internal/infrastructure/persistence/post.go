package persistence

import (
	"database/sql"
	"forum-api/internal/domain/models"
	"time"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type Post struct {
	db *pgx.ConnPool
}

func newPost(db *pgx.ConnPool) *Post {
	return &Post{db: db}
}

func (p *Post) InsertInto(posts []*models.Post, thread *models.Thread) error {
	tx, err := p.db.Begin()
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
		posts[i].Thread = thread.ID
		posts[i].Forum = thread.Forum

		err = tx.QueryRow(
			"post_insert_into",
			posts[i].Author,
			posts[i].Message,
			posts[i].Parent,
			posts[i].Thread,
			posts[i].Forum).Scan(&posts[i].ID, &created)

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

func (p *Post) GetById(post *models.Post) error {
	tx, err := p.db.Begin()
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
	err = tx.QueryRow("post_get_by_id", post.ID).
		Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Path)
	if err != nil {
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (p *Post) Update(post *models.Post) error {
	tx, err := p.db.Begin()
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
	err = tx.QueryRow("post_update", post.Message, post.ID).
		Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	if err != nil {
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (p *Post) Prepare() error {
	_, err := p.db.Prepare("post_insert_into",
		"INSERT INTO post (usr, message,  parent, thread, forum, created) "+
			"VALUES ($1, $2, $3, $4, $5, current_timestamp) "+
			"RETURNING id, created",
	)
	if err != nil {
		return err
	}

	_, err = p.db.Prepare("post_get_by_id",
		"SELECT p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread, p.path "+
			"FROM post p "+
			"WHERE p.id = $1",
	)
	if err != nil {
		return err
	}

	_, err = p.db.Prepare("post_update",
		"UPDATE post SET message = $1, isEdited = true "+
			"WHERE id = $2 "+
			"RETURNING usr, created, forum, isEdited, message, parent, thread",
	)
	if err != nil {
		return err
	}

	_, err = p.db.Prepare("forum_posts_update",
		"UPDATE forum  SET posts = (posts + $1) "+
			"where slug = $2",
	)
	if err != nil {
		return err
	}

	return nil
}
