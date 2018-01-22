-- +migrate Up
DROP INDEX IF EXISTS users_nickname_index;
DROP INDEX IF EXISTS users_email_index;
DROP INDEX IF EXISTS users_email_nickname_index;
DROP INDEX IF EXISTS users_nickname_id_index;

-- +migrate Up
DROP INDEX IF EXISTS forums_low_slug_index;
DROP INDEX IF EXISTS forums_slug_index;

-- +migrate Up
DROP INDEX IF EXISTS threads_slug_index;
DROP INDEX IF EXISTS threads_forum_id_index;
DROP INDEX IF EXISTS threads_forum_index;
DROP INDEX IF EXISTS threads_created_index;
DROP INDEX IF EXISTS threads_forum_created_index;
DROP INDEX IF EXISTS threads_author_index;

-- +migrate Up
DROP INDEX IF EXISTS posts_thread_index;
DROP INDEX IF EXISTS posts_thread_path_index;
DROP INDEX IF EXISTS posts_root_id_index;
DROP INDEX IF EXISTS posts_thread_id_index;
DROP INDEX IF EXISTS posts_thread_parent_index;
DROP INDEX IF EXISTS posts_thread_parent_path_index;

-- +migrate Up
DROP INDEX IF EXISTS votes_user_thread_index;
DROP INDEX IF EXISTS votes_thread_index;
DROP INDEX IF EXISTS forum_users_slug_index;

-- +migrate Up
DROP TRIGGER IF EXISTS parent_path_tgr ON posts;

-- +migrate Up
DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS users;

-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
  id       BIGSERIAL NOT NULL PRIMARY KEY,
  about    TEXT,
  email    TEXT UNIQUE NOT NULL,
  fullname VARCHAR(64) NOT NULL,
  nickname TEXT UNIQUE NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_nickname_index
  ON users (lower(nickname));
CREATE UNIQUE INDEX IF NOT EXISTS users_email_index
  ON users (lower(email));
CREATE INDEX IF NOT EXISTS users_email_nickname_index
  ON users (lower(email), lower(nickname));
CREATE INDEX IF NOT EXISTS users_nickname_id_index
  ON users (id, lower(nickname) DESC);

-- +migrate Up
CREATE TABLE IF NOT EXISTS forums (
  id      BIGSERIAL NOT NULL PRIMARY KEY,
  slug    TEXT UNIQUE NOT NULL,
  author  TEXT,
  title   VARCHAR(255) NOT NULL,
  posts   BIGINT  DEFAULT 0,
  threads INTEGER DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS forums_low_slug_index
  ON forums (lower(slug));

-- +migrate Up
CREATE TABLE IF NOT EXISTS threads (
  id        SERIAL PRIMARY KEY,
  forum     TEXT,
  forum_id  BIGINT,
  author_id BIGINT,
  author    TEXT,
  created   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  message   TEXT NOT NULL,
  slug      TEXT,
  title     VARCHAR(255) NOT NULL,
  votes     INTEGER DEFAULT 0 NOT NULL
);
CREATE INDEX IF NOT EXISTS threads_slug_index
  ON threads (lower(slug));
CREATE INDEX IF NOT EXISTS threads_forum_id_index
  ON threads (forum_id);
CREATE INDEX IF NOT EXISTS threads_forum_index
  ON threads (forum);
CREATE INDEX IF NOT EXISTS threads_forum_created_index
  ON threads (forum, created DESC);
CREATE INDEX IF NOT EXISTS threads_created_index
  ON threads (created);
CREATE INDEX IF NOT EXISTS threads_author_index
  ON threads (author);

-- +migrate Up
CREATE TABLE IF NOT EXISTS posts (
  id        BIGSERIAL PRIMARY KEY,
  forum     TEXT,
  thread    BIGINT,
  author    TEXT,
  created   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  is_edited BOOLEAN NOT NULL DEFAULT FALSE,
  message   TEXT NOT NULL,
  parent    BIGINT DEFAULT 0,
  path      BIGINT [],
  root_id   BIGINT
);
CREATE INDEX IF NOT EXISTS posts_thread_index
  ON posts (thread);
CREATE INDEX IF NOT EXISTS posts_thread_path_index
  ON posts (thread, path DESC);
CREATE INDEX IF NOT EXISTS posts_root_id_index
  ON posts (root_id);
CREATE INDEX IF NOT EXISTS posts_thread_id_index
  ON posts (thread, id DESC);
CREATE INDEX IF NOT EXISTS posts_thread_parent_index
  ON posts (thread, parent, path DESC);
CREATE INDEX IF NOT EXISTS posts_thread_parent_path_index
  ON posts (thread, parent);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_parent_path() RETURNS TRIGGER AS
$update_parent_path$
  BEGIN
    IF (NEW.parent = 0)
      THEN
        NEW.path = ARRAY[NEW.id];
        NEW.root_id = NEW.id;
      ELSE
        NEW.path = (SELECT posts.path || NEW.id FROM posts WHERE id = NEW.parent);
        NEW.root_id = NEW.path[1];
    END IF;
    RETURN NEW;
  END;
$update_parent_path$
LANGUAGE plpgsql;
-- +migrate StatementEnd

-- +migrate Up
CREATE TRIGGER IF NOT EXISTS parent_path_tgr BEFORE INSERT ON posts
FOR EACH ROW EXECUTE PROCEDURE update_parent_path();

-- +migrate Up
CREATE TABLE IF NOT EXISTS votes (
  id        SERIAL NOT NULL PRIMARY KEY,
  author    TEXT REFERENCES users (nickname),
  thread    BIGINT,
  voice     INTEGER NOT NULL
);
CREATE UNIQUE INDEX votes_user_thread_index
  ON votes (lower(author), thread);
CREATE INDEX votes_thread_index
  ON votes (thread);

-- +migrate Up
CREATE TABLE IF NOT EXISTS forum_users (
  author_id  BIGINT,
  forum_id   BIGINT
);
CREATE INDEX forum_users_slug_index
  ON forum_users (forum_id);
