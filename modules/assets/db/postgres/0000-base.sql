-- +migrate Down
DROP INDEX IF EXISTS users_nickname_index;
DROP INDEX IF EXISTS users_email_index;
DROP INDEX IF EXISTS forums_user_index;
DROP INDEX IF EXISTS forums_slug_index;
DROP INDEX IF EXISTS threads_slug_index;
DROP INDEX IF EXISTS threads_forum_index;
DROP INDEX IF EXISTS threads_owner_index;
DROP INDEX IF EXISTS posts_thread_index;
DROP INDEX IF EXISTS posts_parent_index;
DROP INDEX IF EXISTS posts_author_index;
DROP INDEX IF EXISTS posts_parent_path_index;
DROP INDEX IF EXISTS votes_user_thread_index;
DROP TABLE votes;
DROP TABLE posts;
DROP TABLE threads;
DROP TABLE forums;
DROP TABLE users;

-- +migrate Up
CREATE EXTENSION LTREE;

-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
  id       BIGSERIAL NOT NULL PRIMARY KEY,
  about    TEXT,
  email    TEXT UNIQUE NOT NULL,
  fullname VARCHAR(64) NOT NULL,
  nickname TEXT UNIQUE NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_low_nickname_index
  ON users (lower(nickname));
CREATE UNIQUE INDEX IF NOT EXISTS users_email_index
  ON users (lower(email));

-- +migrate Up
CREATE TABLE IF NOT EXISTS forums (
  id      BIGSERIAL NOT NULL PRIMARY KEY,
  slug    TEXT UNIQUE NOT NULL,
  author  TEXT REFERENCES users (nickname),
  title   VARCHAR(255) NOT NULL,
  posts   BIGINT  DEFAULT 0,
  threads INTEGER DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS forums_low_slug_index
  ON forums (lower(slug));
CREATE UNIQUE INDEX IF NOT EXISTS forums_slug_index
  ON forums (slug);

-- +migrate Up
CREATE TABLE IF NOT EXISTS threads (
  id        SERIAL PRIMARY KEY,
  forum     TEXT REFERENCES forums (slug),
  author    TEXT REFERENCES users (nickname),
  created   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  message   TEXT NOT NULL,
  slug      TEXT,
  title     VARCHAR(255) NOT NULL,
  votes     INTEGER DEFAULT 0 NOT NULL
);
CREATE INDEX IF NOT EXISTS threads_slug_index
  ON threads (lower(slug));
CREATE INDEX IF NOT EXISTS threads_id_index
  ON threads (id);
CREATE INDEX IF NOT EXISTS threads_forum_index
  ON threads (forum);
CREATE INDEX IF NOT EXISTS threads_forum_created_index
  ON threads (forum, created);
CREATE INDEX IF NOT EXISTS threads_owner_index
  ON threads (author);

-- +migrate Up
CREATE TABLE IF NOT EXISTS posts (
  id        BIGSERIAL PRIMARY KEY,
  forum     TEXT REFERENCES forums (slug),
  thread    BIGINT REFERENCES threads (id),
  author    TEXT REFERENCES users (nickname),
  created   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  is_edited BOOLEAN NOT NULL DEFAULT FALSE,
  message   TEXT NOT NULL,
  parent    BIGINT DEFAULT 0,
  path      LTREE
);
CREATE INDEX IF NOT EXISTS posts_thread_id_index
  ON posts (thread, id);
CREATE INDEX IF NOT EXISTS posts_thread_index
  ON posts (thread);
CREATE INDEX IF NOT EXISTS posts_id_index
  ON posts (id);
CREATE INDEX IF NOT EXISTS posts_author_index
  ON posts (author);
CREATE INDEX IF NOT EXISTS posts_forum_index
  ON posts (forum);
CREATE INDEX IF NOT EXISTS posts_path_index
  ON posts (thread, path);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_parent_path() RETURNS TRIGGER AS
$update_parent_path$
  DECLARE p LTREE;
  BEGIN
    IF (NEW.parent = 0)
      THEN
        NEW.path = new.id :: TEXT :: LTREE;
    ELSEIF TG_OP = 'INSERT'
      THEN
        SELECT posts.path || NEW.id :: TEXT FROM posts WHERE id = NEW.parent INTO p;
        IF p IS NULL
          THEN
            RAISE EXCEPTION 'Invalid parent_id %', NEW.parent;
        END IF;
        NEW.path = p;
    END IF;
    RETURN NEW;
  END;
$update_parent_path$
LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER parent_path_tgr
BEFORE INSERT
ON posts
FOR EACH ROW
EXECUTE PROCEDURE update_parent_path();

-- +migrate Up
CREATE TABLE IF NOT EXISTS votes (
  id        SERIAL NOT NULL PRIMARY KEY,
  author    TEXT REFERENCES users (nickname),
  thread    BIGINT REFERENCES threads (id),
  voice     INTEGER NOT NULL
);
CREATE UNIQUE INDEX votes_user_thread_index
  ON votes (lower(author), thread);
CREATE INDEX votes_thread_index
  ON votes (thread);
