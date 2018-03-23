CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS errors (
  message TEXT
);

CREATE TABLE IF NOT EXISTS users (
  username CITEXT UNIQUE CONSTRAINT pk__users_username PRIMARY KEY,
  email    CITEXT UNIQUE NOT NULL,
  fullname TEXT          NOT NULL,
  about    TEXT
);

CREATE INDEX IF NOT EXISTS idx__users_username
  ON users (username);

CREATE TABLE IF NOT EXISTS forums (
  id      SERIAL CONSTRAINT pk__forums_id PRIMARY KEY,
  slug    CITEXT UNIQUE NOT NULL,
  posts   BIGINT,
  threads INTEGER,
  title   TEXT,
  creator CITEXT CONSTRAINT fk__forums_creator__users_username REFERENCES users (username)
);

CREATE INDEX IF NOT EXISTS idx__forums_slug
  ON forums (slug);

CREATE TABLE IF NOT EXISTS threads (
  id      SERIAL4 UNIQUE CONSTRAINT pk__threads_id PRIMARY KEY,
  slug    CITEXT,
  author  CITEXT NOT NULL CONSTRAINT fk__threads_author__users_username REFERENCES users (username),
  created TIMESTAMP(3) DEFAULT now(),
  forum   CITEXT NOT NULL CONSTRAINT fk__threads_forum__forums_slug REFERENCES forums (slug),
  message TEXT   NOT NULL,
  title   TEXT   NOT NULL,
  votes   INTEGER      DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx__threads_slug
  ON threads (slug);
CREATE INDEX IF NOT EXISTS idx__threads_id
  ON threads (id);

CREATE TABLE IF NOT EXISTS posts (
  id      SERIAL8 UNIQUE CONSTRAINT pk__posts_id PRIMARY KEY,
  author  CITEXT NOT NULL CONSTRAINT fk__posts_author__users_username REFERENCES users (username),
  created TIMESTAMP(3) DEFAULT now(),
  forum   CITEXT CONSTRAINT fk__posts_forum__forums_slug REFERENCES forums (slug),
  edited  BOOLEAN      DEFAULT FALSE,
  message TEXT   NOT NULL,
  parent  BIGINT       DEFAULT 0,
  thread  INTEGER CONSTRAINT fk__posts_thread__threads_id REFERENCES threads (id)
);

CREATE INDEX IF NOT EXISTS idx__posts_id
  ON posts (id);

CREATE TABLE IF NOT EXISTS votes (
  username CITEXT   NOT NULL CONSTRAINT fo__votes_username__users_username REFERENCES users (username),
  voice    SMALLINT NOT NULL
);
