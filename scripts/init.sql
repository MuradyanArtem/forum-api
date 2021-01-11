
ALTER SYSTEM
SET checkpoint_completion_target
= '0.9';

ALTER SYSTEM
SET wal_buffers
= '6912kB';

ALTER SYSTEM
SET default_statistics_target
= '100';

ALTER SYSTEM
SET effective_io_concurrency
= '200';

ALTER SYSTEM
SET seq_page_cost
= '0.1';

ALTER SYSTEM
SET random_page_cost
= '0.1';

ALTER SYSTEM
SET max_worker_processes
= '4';

ALTER SYSTEM
SET max_parallel_workers_per_gather
= '2';

ALTER SYSTEM
SET max_parallel_workers
= '4';

ALTER SYSTEM
SET max_parallel_maintenance_workers
= '2';

CREATE EXTENSION
IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C"             NOT NULL UNIQUE,
    email    CITEXT COLLATE "C"             NOT NULL UNIQUE,
    about    TEXT                           NOT NULL,
    fullname TEXT                           NOT NULL
);

CREATE INDEX index_users_nickname_hash ON users using hash (nickname);
CREATE INDEX index_users_email_hash ON users using hash (email);

CREATE UNLOGGED TABLE forums
(
    slug     CITEXT COLLATE "C" PRIMARY KEY                                   NOT NULL,
    title    TEXT                                                             NOT NULL,
    nickname CITEXT COLLATE "C" REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    posts    INTEGER DEFAULT 0                                                NOT NULL,
    threads  INTEGER DEFAULT 0                                                NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users (nickname)
);

CREATE INDEX ON forums (slug, title, nickname, posts, threads);
CREATE INDEX ON forums USING hash (slug);
CREATE INDEX ON forums (nickname);


CREATE UNLOGGED TABLE forum_users
(
    author CITEXT COLLATE "C" REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    slug   CITEXT COLLATE "C" REFERENCES forums (slug)    ON DELETE CASCADE NOT NULL,
    FOREIGN KEY (slug)        REFERENCES forums (slug)    ON DELETE CASCADE,
    FOREIGN KEY (author)      REFERENCES users (nickname) ON DELETE CASCADE,
    UNIQUE      (slug, author)
);

CREATE INDEX index_forum_users ON forum_users (slug, author);
CREATE INDEX ON forum_users (author);
CLUSTER forum_users USING index_forum_users;

CREATE UNLOGGED TABLE threads
(
    author     CITEXT COLLATE "C" REFERENCES users(nickname) ON DELETE CASCADE    NOT NULL,
    created    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP                 NOT NULL,
    forum_slug CITEXT COLLATE "C" REFERENCES forums (slug)   ON DELETE CASCADE    NOT NULL,
    id         SERIAL PRIMARY KEY                                                 NOT NULL,
    message    TEXT                                                               NOT NULL,
    slug       CITEXT COLLATE "C",
    title      TEXT                                                               NOT NULL,
    votes      INTEGER DEFAULT 0                                                  NOT NULL,
    FOREIGN KEY (forum_slug) REFERENCES forums (slug)        ON DELETE CASCADE,
    FOREIGN KEY (author)     REFERENCES users (nickname)     ON DELETE CASCADE

);

CREATE INDEX ON threads (slug, created);
CREATE INDEX ON index_threads_created on threads (created);

CREATE INDEX ON threads using hash (slug);
CREATE INDEX ON threads using hash (id);

CREATE UNLOGGED TABLE posts
(
    author     CITEXT COLLATE "C" REFERENCES users (nickname) ON DELETE CASCADE  NOT NULL,
    created    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP                NOT NULL,
    forum_slug CITEXT COLLATE "C" REFERENCES forums (slug)    ON DELETE CASCADE  NOT NULL,
    id         SERIAL PRIMARY KEY                                                NOT NULL,
    edited     BOOL DEFAULT 'false'                                              NOT NULL,
    message    TEXT                                                              NOT NULL,
    parent     INTEGER                                                           NOT NULL,
    thread     INTEGER REFERENCES threads (id)                ON DELETE CASCADE  NOT NULL,
    path       INTEGER ARRAY DEFAULT '{}'                                        NOT NULL,
    FOREIGN KEY (forum_slug) REFERENCES forums (slug)         ON DELETE CASCADE,
    FOREIGN KEY (author)     REFERENCES users (nickname)      ON DELETE CASCADE,
    FOREIGN KEY (thread)     REFERENCES threads (id)          ON DELETE CASCADE
);

CREATE INDEX ON posts (id);
CREATE INDEX ON posts (thread, created, id);
CREATE INDEX ON posts (thread, id);
CREATE INDEX ON posts (thread, path);
CREATE INDEX ON posts (thread, parent, path);
CREATE INDEX ON posts ((path[1]), path);

CREATE UNLOGGED TABLE votes
(
    nickname  CITEXT COLLATE "C" REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    thread_id INTEGER            REFERENCES threads (id)     ON DELETE CASCADE NOT NULL,
    vote      SMALLINT                                                         NOT NULL,
    FOREIGN KEY (thread_id)      REFERENCES threads (id),
    FOREIGN KEY (nickname)       REFERENCES users (nickname),
    UNIQUE      (thread_id, nickname)
);

CREATE UNIQUE INDEX ON votes (thread_id, nickname);


-- PATH TO POST UPDATE
CREATE FUNCTION update_path() RETURNS TRIGGER AS
$$
DECLARE
    temp INT ARRAY;
    t    INTEGER;
BEGIN
    IF new.parent ISNULL OR new.parent = 0 THEN
        new.path = ARRAY [new.id];
    ELSE
        SELECT thread
        INTO t
        FROM posts
        WHERE id = new.parent;
        IF t ISNULL OR t <> new.thread THEN
            RAISE EXCEPTION 'Not in this thread ID ' USING HINT = 'Please check your parent ID';
        END IF;

        SELECT path
        INTO temp
        FROM posts
        WHERE id = new.parent;
        new.path = array_append(temp, new.id);

    END IF;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();

-- VOTE VALUE UPDATE
CREATE FUNCTION vote_count_upd() RETURNS TRIGGER AS
$$
BEGIN
    IF (old.vote != new.vote) THEN
        UPDATE threads
        SET votes = (votes - old.vote + new.vote)
        WHERE id = new.thread_id;
    END IF;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER vote_count_upd
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE vote_count_upd();

CREATE FUNCTION vote_count_insert() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE threads
    SET votes = (votes + new.vote)
    WHERE id = new.thread_id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER vote_count_insert
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE vote_count_insert();


-- UPDATE FORUM_USER TABLE AFTER INSERTS
CREATE FUNCTION insert_forum_user_from_threads_or_psoots() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO forum_users
    VALUES (new.author, new.forum_slug)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_user_from_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user_from_threads_or_psoots();

CREATE TRIGGER update_forum_user_from_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user_from_threads_or_psoots();

-- UPDATE POSTS AND THREADS COUNTERS IN FORUMS
CREATE FUNCTION update_forum_counter_posts() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = new.forum_slug;

    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_counters_after_post_insert
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_counter_posts();

CREATE FUNCTION update_forum_counter_threads() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = new.forum_slug;

    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_counters_after_thread_insert
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_counter_threads();

VACUUM ANALYSE;
