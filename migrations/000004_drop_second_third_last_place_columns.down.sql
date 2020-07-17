BEGIN;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS second_place INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS third_place INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS last_place INTEGER DEFAULT 0;

CREATE INDEX IF NOT EXISTS users_ranking ON users using btree (games_won, second_place, third_place, last_place);

COMMIT;
