package persistence

import (
	"errors"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"strconv"

	"github.com/jackc/pgx"
)

type Thread struct {
	db *pgx.ConnPool
}

func newThread(db *pgx.ConnPool) *Thread {
	return &Thread{db: db}
}

func (t *Thread) InsertThread(thread *models.Thread) error {
	err := t.db.QueryRow("INSERT INTO threads (author, created, "+
		"forum_slug, message, slug, title) "+
		"VALUES ($1, $2, $3, $4, $5, $6) "+
		"RETURNING id", &thread.Author, &thread.Created, &thread.Forum, &thread.Message, &thread.Slug,
		&thread.Title).Scan(&thread.ID)
	if err != nil {
		switch infrastructure.ErrCode(err) {
		case infrastructure.PgErrUniqueViolation:
			return infrastructure.ErrConflict
		default:
			return errors.New(err.Error() + " " + thread.Forum + " " + thread.Author)
		}
	}
	return nil
}

func (t *Thread) SelectThreadByID(id int) (*models.Thread, error) {
	thread := &models.Thread{}
	err := t.db.QueryRow("SELECT * FROM threads WHERE id = $1", id).
		Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.ID,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (t *Thread) SelectThreadsByForum(slug string, limit int, since string, desc bool) (models.Threads, error) {
	threads := []models.Thread{}

	query := "SELECT author, created, forum_slug, id, message, slug, title, votes FROM threads WHERE forum_slug = $1 "
	rows := &pgx.Rows{}
	var err error
	if limit > 0 && since != "" {
		if desc {
			query += "AND created <= $2 ORDER BY created DESC LIMIT $3"
		} else {
			query += "AND created >= $2 ORDER BY created ASC LIMIT $3"
		}
		rows, err = t.db.Query(query, &slug, &since, &limit)
	} else {
		if limit > 0 {
			if desc {
				query += "ORDER BY created DESC LIMIT $2"
			} else {
				query += "ORDER BY created ASC LIMIT $2"
			}
			rows, err = t.db.Query(query, &slug, &limit)
		}
		if since != "" {
			if desc {
				query += "AND created <= $2 ORDER BY created DESC"
			} else {
				query += "AND created >= $2 ORDER BY created ASC"
			}
			rows, err = t.db.Query(query, &slug, &since)
		}
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		thread := models.Thread{}

		err := rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug,
			&thread.Title, &thread.Votes)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}

func (t *Thread) SelectThreadBySlug(slug string) (*models.Thread, error) {
	thread := models.Thread{}
	err := t.db.QueryRow("SELECT author, created, forum_slug, id, message, slug, title, votes FROM threads WHERE slug = $1", slug).
		Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.ID,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (t *Thread) Update(thread *models.Thread) error {
	err := t.db.QueryRow("UPDATE threads "+
		"SET author = $1, "+
		"forum_slug = $2, "+
		"message = $3, "+
		"slug = $4, "+
		"title = $5  "+
		"WHERE id = $6 "+
		"RETURNING threads.* ", &thread.Author, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.ID).
		Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.ID,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
	return err
}

func (t *Thread) UpdateBySlugOrID(s string, thread *models.Thread) error {
	query := "UPDATE THREADS SET   "
	if thread.Author != "" {
		query += "author = '" + thread.Author + "', "
	}
	if thread.Message != "" {
		query += "message = '" + thread.Message + "', "
	}
	if thread.Title != "" {
		query += "title = '" + thread.Title + "', "
	}
	query = query[:len(query)-2]
	value, err := strconv.Atoi(s)
	if err != nil {
		value = -1
	}
	query += " WHERE id = $1 OR slug = $2 RETURNING threads.*"
	if thread.Author == "" && thread.Message == "" && thread.Title == "" {
		query = "SELECT * FROM threads WHERE id = $1 OR slug = $2"
	}
	err = t.db.QueryRow(query, &value, &s).
		Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.ID,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
	switch err {
	case nil:
		return nil
	default:
		return err
	}
}

func (t *Thread) VoteBySlug(vote models.Vote, slug string) (models.Thread, error) {
	thread, err := t.SelectThreadBySlug(slug)
	if err != nil {
		return models.Thread{}, err
	}
	_, err = t.db.Exec(`INSERT INTO votes (thread_id, nickname, vote)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, nickname) DO UPDATE SET vote = $3`,
		thread.ID,
		vote.Nickname,
		vote.Voice,
	)
	if err != nil {
		return models.Thread{}, err
	}
	err = t.db.QueryRow(
		`SELECT votes FROM threads WHERE id = $1`,
		thread.ID).Scan(&thread.Votes)
	if err != nil {
		return models.Thread{}, err
	}
	return *thread, nil
}

func (t *Thread) VoteByID(vote models.Vote, id int) (models.Thread, error) {
	_, err := t.db.Exec(`
			INSERT INTO votes (thread_id, nickname, vote)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, nickname) DO UPDATE SET vote = $3`,
		id,
		vote.Nickname,
		vote.Voice,
	)
	if err != nil {
		return models.Thread{}, err
	}
	thread := &models.Thread{}
	if thread, err = t.SelectThreadByID(id); err != nil {
		return models.Thread{}, err
	}
	return *thread, nil
}

func (t *Thread) GetForumIDBySlug(s string) (int, string, error) {
	forum := ""
	res := 0
	err := t.db.QueryRow("SELECT id, forum_slug FROM threads WHERE slug = $1", s).Scan(&res, &forum)
	return res, forum, err
}

func (r *Thread) SelectForumByThreadID(id int) (string, error) {
	forum := ""
	err := t.db.QueryRow("SELECT forum_slug FROM threads WHERE id = $1", id).Scan(&forum)
	return forum, err
}
