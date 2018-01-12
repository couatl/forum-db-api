-- +migrate Up
CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS LTREE;

CREATE TABLE IF NOT EXISTS users (
  id       BIGSERIAL NOT NULL PRIMARY KEY,
  about    TEXT,
  email    CITEXT UNIQUE NOT NULL,
  fullname VARCHAR(64) NOT NULL,
  nickname CITEXT UNIQUE NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS users_nickname_index
  ON users (nickname);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_index
  ON users (email);

CREATE TABLE IF NOT EXISTS forums (
  id      BIGSERIAL NOT NULL PRIMARY KEY,
  slug    CITEXT UNIQUE NOT NULL,
  user_id BIGINT REFERENCES users (id),
  title   VARCHAR(255) NOT NULL,
  posts   BIGINT  DEFAULT 0,
  threads INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS forums_user_index
  ON forums (user_id);

CREATE UNIQUE INDEX IF NOT EXISTS forums_slug_index
  ON forums (slug);

CREATE TABLE IF NOT EXISTS threads (
  id        SERIAL PRIMARY KEY,
  forum_id  BIGINT REFERENCES forums (id),
  author_id BIGINT REFERENCES users (id),
  created   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  message   TEXT NOT NULL,
  slug      CITEXT UNIQUE,
  title     VARCHAR(255) NOT NULL,
  votes     INTEGER DEFAULT 0 NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS threads_slug_index
  ON threads (lower(slug));

CREATE INDEX IF NOT EXISTS threads_forum_index
  ON threads (forum_id);

CREATE INDEX IF NOT EXISTS threads_owner_index
  ON threads (author_id);

CREATE TABLE IF NOT EXISTS posts (
  id        BIGSERIAL PRIMARY KEY,
  forum_id  BIGINT REFERENCES forums (id),
  thread_id INTEGER REFERENCES threads (id),
  author_id BIGINT REFERENCES users (id),
  created   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  is_edited BOOLEAN NOT NULL DEFAULT FALSE,
  message   TEXT NOT NULL,
  parent    BIGINT DEFAULT 0 NOT NULL,
  path      LTREE
);

CREATE INDEX IF NOT EXISTS posts_thread_index
  ON posts (thread_id);

CREATE INDEX IF NOT EXISTS posts_parent_index
  ON posts (parent);

CREATE INDEX IF NOT EXISTS posts_author_index
  ON posts (author_id);

CREATE INDEX IF NOT EXISTS posts_parent_path_index
  ON posts USING GIST (path);

CREATE TABLE IF NOT EXISTS votes (
  id        SERIAL NOT NULL PRIMARY KEY,
  user_id   BIGINT NOT NULL REFERENCES users (id),
  thread_id INTEGER REFERENCES threads (id),
  voice     INTEGER NOT NULL
);

CREATE UNIQUE INDEX votes_user_thread_index
  ON votes (user_id, thread_id);
