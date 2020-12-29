package persistence

import (
	"database/sql"
	"forum-api/internal/domain/models"
	"strconv"
	"time"

	"github.com/jackc/pgx"
)

type Thread struct {
	db *pgx.ConnPool
}

func newThread(db *pgx.ConnPool) *Thread {
	return &Thread{db: db}
}

func (t *Thread) InsertIntoForumUsers(forum, nickname string) error {
	tx, err := t.db.Begin()
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

	var buffer string
	err = tx.QueryRow("get_forum_user", forum, nickname).Scan(&buffer)
	if err != nil {
		_, err = t.db.Exec("forum_users_insert_into", forum, nickname)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Thread) InsertInto(thread *models.Thread) error {
	tx, err := t.db.Begin()
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

	slug := &sql.NullString{}
	if thread.Slug != "" {
		slug.String = thread.Slug
		slug.Valid = true
	}

	created := &sql.NullString{}
	if thread.Created != "" {
		created.String = thread.Created
		created.Valid = true
	}

	row := tx.QueryRow("thread_insert_into", thread.Author, created, thread.Forum, thread.Message, thread.Title, slug)
	if err := row.Scan(&thread.ID); err != nil {
		return err
	}

	if err = t.InsertIntoForumUsers(thread.Forum, thread.Author); err != nil {
		return err
	}

	return nil
}

func (t *Thread) GetBySlug(thread *models.Thread) error {
	tx, err := t.db.Begin()
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

	row := tx.QueryRow("thread_get_by_slug", thread.Slug)

	created := sql.NullTime{}
	slug := sql.NullString{}

	if err = row.Scan(
		&thread.ID,
		&thread.Title,
		&thread.Message,
		&created,
		&slug,
		&thread.Author,
		&thread.Forum,
		&thread.Votes,
	); err != nil {
		return err
	}

	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}

	thread.Slug = ""
	if slug.Valid {
		thread.Slug = slug.String
	}

	return nil
}

func (t *Thread) GetById(thread *models.Thread) error {
	tx, err := t.db.Begin()
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

	row := tx.QueryRow("thread_get_by_id", thread.ID)

	created := sql.NullTime{}
	slug := sql.NullString{}

	if err = row.Scan(
		&thread.ID,
		&thread.Title,
		&thread.Message,
		&created,
		&slug,
		&thread.Author,
		&thread.Forum,
		&thread.Votes,
	); err != nil {
		return err
	}

	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}

	thread.Slug = ""
	if slug.Valid {
		thread.Slug = slug.String
	}

	return nil
}

func (t *Thread) GetBySlugOrId(thread *models.Thread) error {
	tx, err := t.db.Begin()
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

	Id, err := strconv.ParseInt(thread.Slug, 10, 64)
	if err == nil {
		thread.ID = Id
		thread.Slug = ""
	}

	if thread.Slug != "" {
		err = t.GetBySlug(thread)
	} else {
		err = t.GetById(thread)
	}

	if err != nil {
		return err
	}

	return nil
}

func (t *Thread) InsertIntoVotes(thread *models.Thread, vote *models.Vote) error {
	tx, err := t.db.Begin()
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

	voteNum := int32(0)
	err = tx.QueryRow("votes_get_info", vote.Nickname, vote.Thread).Scan(&voteNum)

	if voteNum == 0 {
		err = tx.QueryRow("votes_insert_into", vote.Vote, vote.Nickname, vote.Thread).Scan(&vote.Thread)
		thread.Votes += int64(vote.Vote)
	} else {
		if voteNum != vote.Vote {
			err = tx.QueryRow("votes_update", vote.Vote, vote.Nickname, vote.Thread).Scan(&vote.Thread)
			thread.Votes += 2 * int64(vote.Vote)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (t *Thread) Update(thread *models.Thread) error {
	tx, err := t.db.Begin()
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

	slug := sql.NullString{}
	created := sql.NullTime{}
	votes := sql.NullInt64{}

	switch true {
	case thread.Message == "" && thread.Title == "":
		err = t.GetBySlugOrId(thread)
	case thread.Message != "" && thread.Title == "":
		err = tx.QueryRow("thread_update_message",
			thread.Message,
			thread.Slug,
		).Scan(
			&thread.ID,
			&thread.Title,
			&thread.Message,
			&created,
			&slug,
			&thread.Author,
			&thread.Forum,
			&votes,
		)
	case thread.Message == "" && thread.Title != "":
		err = tx.QueryRow("thread_update_title",
			thread.Title,
			thread.Slug,
		).Scan(
			&thread.ID,
			&thread.Title,
			&thread.Message,
			&created,
			&slug,
			&thread.Author,
			&thread.Forum,
			&votes,
		)
	case thread.Message != "" && thread.Title != "":
		err = tx.QueryRow("thread_update_all",
			thread.Message,
			thread.Title,
			thread.Slug,
		).Scan(
			&thread.ID,
			&thread.Title,
			&thread.Message,
			&created,
			&slug,
			&thread.Author,
			&thread.Forum,
			&votes,
		)
	}
	if err != nil {
		return err
	}

	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}

	if slug.Valid {
		thread.Slug = slug.String
	}

	if votes.Valid {
		thread.Votes = votes.Int64
	}

	return nil
}

func (t *Thread) GetPosts(thread *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	tx, err := t.db.Begin()
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

	posts := make([]models.Post, 0)
	var rows *pgx.Rows

	if sort == "tree" {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_asc", thread.ID)

		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_desc", thread.ID)

		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_asc_with_since", thread.ID, since)

		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_desc_with_since", thread.ID, since)

		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_asc_with_limit", thread.ID, limit)

		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_desc_with_limit", thread.ID, limit)

		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_asc_with_since_with_limit", thread.ID, since, limit)

		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_desc_with_since_with_limit", thread.ID, since, limit)
		}
	} else if sort == "parent_tree" {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_asc", thread.ID)

		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_desc", thread.ID)

		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_asc_with_since", thread.ID, since)

		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_desc_with_since", thread.ID, since)

		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_asc_with_limit", thread.ID, limit)

		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_desc_with_limit", thread.ID, limit)

		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_asc_with_since_with_limit", thread.ID, since, limit)

		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_desc_with_since_with_limit", thread.ID, since, limit)
		}
	} else {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_post_flat_asc", thread.ID)

		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_post_flat_desc", thread.ID)

		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_post_flat_asc_with_since", thread.ID, since)

		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_post_flat_desc_with_since", thread.ID, since)

		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_post_flat_asc_with_limit", thread.ID, limit)

		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_post_flat_desc_with_limit", thread.ID, limit)

		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_post_flat_asc_with_since_with_limit", thread.ID, since, limit)

		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_post_flat_desc_with_since_with_limit", thread.ID, since, limit)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		created := sql.NullTime{}
		p := models.Post{}

		err := rows.Scan(&p.ID, &p.Author, &created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
		if err != nil {
			return nil, err
		}

		if created.Valid {
			p.Created = created.Time.Format(time.RFC3339Nano)
		}

		posts = append(posts, p)
	}
	rows.Close()

	return posts, nil
}

func (t *Thread) Prepare() error {
	_, err := t.db.Prepare("thread_insert_into",
		"INSERT INTO thread (usr, created, forum, message, title, slug) VALUES ($1, $2, $3, $4, $5, $6)"+
			"ON CONFLICT DO NOTHING "+
			"RETURNING id",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_forum_user",
		"SELECT nickname FROM forum_users "+
			"WHERE forum = $1 AND nickname = $2 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("forum_users_insert_into",
		"INSERT INTO forum_users (forum, nickname) "+
			"VALUES ($1,$2) ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_get_by_slug",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes  "+
			"FROM thread t "+
			"WHERE t.slug = $1",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_get_by_id",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes "+
			"FROM thread t "+
			"WHERE t.id = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("votes_insert_into",
		"INSERT INTO vote (vote, usr, thread) VALUES ($1 , $2, $3) "+
			"RETURNING thread",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("votes_update",
		"UPDATE vote SET vote = $1 WHERE usr = $2 and thread = $3 "+
			"RETURNING thread",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("votes_get_info",
		"SELECT vote FROM vote "+
			"WHERE usr = $1 and thread = $2 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_update_all",
		"UPDATE thread SET "+
			"message = $1, "+
			"title = $2 "+
			"WHERE id::citext = $3 or slug = $3 "+
			"RETURNING id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_update_message",
		"UPDATE thread SET "+
			"message = $1 "+
			"WHERE id::citext = $2 or slug = $2 "+
			"RETURNING id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_update_title",
		"UPDATE thread SET "+
			"title = $1 "+
			"WHERE id::citext = $2 or slug = $2 "+
			"RETURNING id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_asc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_desc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path DESC ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_asc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_desc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path DESC "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_asc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] > (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_desc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] < (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path DESC ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_asc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] > (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_tree_desc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] < (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path DESC "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_asc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_desc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path DESC "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_asc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path "+
			"	LIMIT $2"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_desc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path DESC "+
			"	LIMIT $2"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_asc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] > (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_desc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] < (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path DESC "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_asc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] > (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path "+
			"	LIMIT $3"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_posts_parent_desc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] < (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path DESC "+
			"	LIMIT $3"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_asc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created, p.id",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_desc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created DESC , p.id DESC ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_asc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created, p.id "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_desc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created DESC , p.id DESC "+
			"LIMIT $2",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_asc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id > $2 "+
			"ORDER BY p.created, p.id",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_desc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id < $2 "+
			"ORDER BY p.created DESC , p.id DESC ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_asc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id > $2 "+
			"ORDER BY p.created, p.id "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("thread_post_flat_desc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id < $2 "+
			"ORDER BY p.created DESC , p.id DESC "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	return nil
}
