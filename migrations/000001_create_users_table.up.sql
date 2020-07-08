BEGIN;

CREATE TABLE IF NOT EXISTS users(
  id serial primary key,
  username VARCHAR (300) UNIQUE NOT NULL,
  guild_id VARCHAR (100) NOT NULL
);

CREATE INDEX IF NOT EXISTS users_id ON users using btree (id);
CREATE INDEX IF NOT EXISTS users_guild_id ON users using btree (guild_id);

COMMIT;
