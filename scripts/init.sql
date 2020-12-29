
ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET wal_buffers = '6912kB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET seq_page_cost = '0.1';
ALTER SYSTEM SET random_page_cost = '0.1';
ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';

DROP trigger IF EXISTS path_updater                   ON post;
DROP trigger IF EXISTS thread_updater                 ON thread;
DROP trigger IF EXISTS forum_users_clear              ON forum_users;
DROP trigger IF EXISTS forum_user_insert_after_thread ON thread;
DROP trigger IF EXISTS forum_user_insert_after_post   ON post;
DROP trigger IF EXISTS insert_into_thread_votes       ON thread;
DROP trigger IF EXISTS update_thread_votes            ON thread;
DROP trigger IF EXISTS upd_forum_threads              ON thread;

DROP FUNCTION IF EXISTS update_path;
DROP FUNCTION IF EXISTS update_forum_users;
DROP FUNCTION IF EXISTS insert_thread_votes;
DROP FUNCTION IF EXISTS insert_into_forum_users;
DROP FUNCTION IF EXISTS update_forum_threads;

DROP TABLE IF EXISTS usr         CASCADE;
DROP TABLE IF EXISTS forum       CASCADE;
DROP TABLE IF EXISTS thread      CASCADE;
DROP TABLE IF EXISTS post        CASCADE;
DROP TABLE IF EXISTS vote        CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE EXTENSION IF NOT EXISTS CITEXT;

------------------------ USER ------------------------------------------
CREATE UNLOGGED TABLE usr
(
    id       SERIAL             PRIMARY KEY,
    email    CITEXT COLLATE "C" NOT NULL UNIQUE,
    nickname CITEXT COLLATE "C" NOT NULL UNIQUE,
    fullname TEXT               NOT NULL,
    about    TEXT
);

CREATE INDEX index_user_all ON usr (nickname, fullname, email, about);
CLUSTER usr USING index_user_all;

------------------------- FORUM ------------------------------------
CREATE UNLOGGED TABLE forum
(
    id      SERIAL             PRIMARY KEY,
    slug    CITEXT COLLATE "C" NOT NULL UNIQUE,
    usr     CITEXT COLLATE "C" NOT NULL REFERENCES usr (nickname) ON DELETE CASCADE,
    title   TEXT               NOT NULL,
    threads BIGINT DEFAULT 0,
    posts   BIGINT DEFAULT 0
);

CREATE INDEX index_forum_slug_hash ON forum USING HASH (slug);
CREATE INDEX index_usr_fk          ON forum (usr);
CREATE INDEX index_forum_all       ON forum (slug, title, usr, posts, threads);

------------------------- THREAD ---------------------------------------
CREATE UNLOGGED TABLE thread
(
    id      SERIAL             PRIMARY KEY,
    title   TEXT               NOT NULL,
    message TEXT               NOT NULL,
    created TIMESTAMP WITH TIME ZONE,
    slug    CITEXT COLLATE "C" UNIQUE,
    votes   INT DEFAULT 0,
    usr     CITEXT COLLATE "C" NOT NULL REFERENCES usr (nickname) ON DELETE CASCADE,
    forum   CITEXT COLLATE "C" NOT NULL REFERENCES forum (slug)   ON DELETE CASCADE
);

CREATE INDEX index_thread_forum_created ON thread (forum, created);
CREATE INDEX index_thread_slug          ON thread (slug);
CREATE INDEX index_thread_slug_hash     ON thread USING HASH (slug);
CREATE INDEX index_thread_all           ON thread (title, message, CREATEd, slug, usr, forum, votes);
CREATE INDEX index_thread_usr_fk        ON thread (usr);
CREATE INDEX index_thread_forum_fk      ON thread (forum);

------------------------- POST --------------------------------------------------------------
CREATE UNLOGGED TABLE post
(
    id       BIGSERIAL             PRIMARY KEY,
    message  TEXT                  NOT NULL,
    isedited BOOLEAN DEFAULT false NOT NULL,
    parent   INTEGER DEFAULT 0,
    created  TIMESTAMP,
    usr      CITEXT COLLATE "C"    NOT NULL REFERENCES usr (nickname) ON DELETE CASCADE,
    thread   INTEGER               NOT NULL REFERENCES thread         ON DELETE CASCADE,
    forum    CITEXT COLLATE "C"    NOT NULL REFERENCES forum (slug)   ON DELETE CASCADE,
    path     BIGINT[]
);

CREATE INDEX index_post_thread_id          ON post (thread, id);
CREATE INDEX index_post_thread_path        ON post (thread, path);
CREATE INDEX index_post_thread_parent_path ON post (thread, parent, path);
CREATE INDEX index_post_path1_path         ON post ((path[1]), path);
CREATE INDEX index_post_thread_created_id  ON post (thread, CREATEd, id);
CREATE INDEX index_post_usr_fk             ON post (usr);
CREATE INDEX index_post_forum_fk           ON post (forum);

--------------------------- VOTE ------------------------------------------
CREATE UNLOGGED TABLE vote
(
    id     SERIAL             PRIMARY KEY,
    vote   INTEGER            NOT NULL,
    usr    CITEXT COLLATE "C" NOT NULL references usr (nickname) ON DELETE CASCADE,
    thread INTEGER            NOT NULL references thread         ON DELETE CASCADE
);

CREATE UNIQUE INDEX vote_user_thread_unique ON vote (usr, thread);
CREATE INDEX index_vote_thread              ON vote (thread);

------------------------------ FORUM USERS -----------------------------------
CREATE UNLOGGED TABLE forum_users
(
    forum    CITEXT COLLATE "C" NOT NULL REFERENCES forum (slug)   ON DELETE CASCADE,
    nickname CITEXT COLLATE "C" NOT NULL REFERENCES usr (nickname) ON DELETE CASCADE
);

CREATE UNIQUE INDEX INDEX_forum_nickname ON forum_users (forum, nickname);
CLUSTER forum_users USING INDEX_forum_nickname;

---------------------- UPDATE PATH ---------------------------
CREATE OR REPLACE FUNCTION update_path()
RETURNS trigger AS
$BODY$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT thread, path
        FROM post
        WHERE thread = NEW.thread AND id = NEW.parent
        INTO first_parent_thread, parent_path;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Parent post not found in current thread' USING ERRCODE = '00404';
        END IF;
        NEW.path := parent_path || NEW.id;
    END IF;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER path_updater
    BEFORE INSERT
    ON post
    FOR EACH ROW
    EXECUTE PROCEDURE update_path();

-------------------------------- INSERT THREAD VOTES -----------------------
CREATE OR REPLACE FUNCTION insert_thread_votes()
RETURNS trigger AS
$BODY$
DECLARE
BEGIN
UPDATE thread SET votes = (votes + NEW.vote) WHERE id = NEW.thread;
RETURN new;
end;
$BODY$ LANGUAGE plpgsql;

CREATE trigger insert_thread_votes
    BEFORE INSERT
    ON vote
    FOR EACH ROW
    EXECUTE PROCEDURE insert_thread_votes();

------------------------------- UPDATE THREAD VOTES -------------------------
CREATE OR REPLACE FUNCTION update_thread_votes()
RETURNS trigger AS
$BODY$
BEGIN
    IF NEW.vote > 0 THEN
        UPDATE thread SET votes = (votes + 2) WHERE id = NEW.thread;
    ELSE
        UPDATE thread SET votes = (votes - 2) WHERE id = NEW.thread;
    END IF;
    RETURN new;
END;
$BODY$ LANGUAGE plpgsql;

CREATE trigger update_thread_votes
    BEFORE UPDATE
    ON vote
    FOR EACH ROW
    EXECUTE PROCEDURE update_thread_votes();

------------------------------- UPDATE FORUM USERS -------------------
CREATE OR REPLACE FUNCTION update_forum_users()
RETURNS TRIGGER AS
$BODY$
BEGIN
    INSERT INTO forum_user (forum, nickname) VALUES (NEW.forum, NEW.author)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER forum_usr_updater
    AFTER INSERT
    ON thread
    FOR EACH ROW
    EXECUTE PROCEDURE update_forum_users();

------------------------------- UPDATE FORUM THREADS -------------------
CREATE OR REPLACE FUNCTION update_forum_threads()
RETURNS trigger AS
$BODY$
BEGIN
    UPDATE forum SET threads = (threads + 1) WHERE slug = NEW.forum;
    RETURN NEW;
END;
$BODY$ language plpgsql;

CREATE TRIGGER upd_forum_threads
    AFTER INSERT
    ON thread
    FOR EACH row
    EXECUTE PROCEDURE update_forum_threads();
