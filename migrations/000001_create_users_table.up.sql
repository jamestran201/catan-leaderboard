BEGIN;

CREATE TABLE IF NOT EXISTS users(
  id serial primary key,
  username VARCHAR (300) UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS users_id ON users using btree (id);
CREATE INDEX IF NOT EXISTS users_username ON users using btree (username);

COMMIT;
