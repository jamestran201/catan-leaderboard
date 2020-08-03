BEGIN;

ALTER TABLE users ADD COLUMN IF NOT EXISTS points_per_game REAL DEFAULT 0;

UPDATE users
SET points_per_game =
CASE WHEN games = 0 THEN 0
ELSE CAST(points as real) / games END;

COMMIT;
