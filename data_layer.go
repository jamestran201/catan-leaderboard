package main

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type dataLayer interface {
	addUser(username string, guildID string) error
	getTopTwentyUsers(guildID string) ([]user, error)
	checkUserExists(username string, guildID string) (int, error)
	addWin(username string, guildID string) error
	updateGameStats(username string, points string, guildID string) error
}

type postgresDataLayer struct {
	dbConn *pgx.Conn
}

func (db *postgresDataLayer) addUser(username string, guildID string) error {
	_, err := db.dbConn.Exec(context.Background(), "INSERT INTO users (username, guild_id) VALUES ($1, $2)", username, guildID)

	return err
}

func (db *postgresDataLayer) getTopTwentyUsers(guildID string) ([]user, error) {
	rows, err := db.dbConn.Query(
		context.Background(),
		`SELECT
			CAST(RANK() OVER (ORDER BY games_won DESC) AS TEXT), username, CAST(games_won AS TEXT) ,
			CAST(points AS TEXT), CAST(games AS TEXT), points_per_game
		FROM users
		WHERE guild_id = ($1)
		LIMIT 20`,
		guildID,
	)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	users := make([]user, 0, 20)

	for i := 0; rows.Next(); i++ {
		user := user{}
		err = rows.Scan(
			&user.rank,
			&user.username,
			&user.victories,
			&user.points,
			&user.games,
			&user.pointsPerGame,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (db *postgresDataLayer) checkUserExists(username string, guildID string) (int, error) {
	row := db.dbConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username = ($1) AND guild_id = ($2);", username, guildID)

	var recordExists int
	err := row.Scan(&recordExists)

	if err != nil {
		return -1, err
	}

	return recordExists, nil
}

func (db *postgresDataLayer) addWin(username string, guildID string) error {
	_, err := db.dbConn.Exec(context.Background(), "UPDATE users SET games_won = games_won + 1 WHERE username = ($1) AND guild_id = ($2)", username, guildID)

	return err
}

func (db *postgresDataLayer) updateGameStats(username string, points string, guildID string) error {
	tx, err := db.dbConn.Begin(context.Background())

	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		`UPDATE users
		SET
			points = points + ($1),
			games = games + 1
		WHERE username = ($2) AND guild_id = ($3)`,
		points, username, guildID,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), "UPDATE users SET points_per_game = CAST(points as real) / games WHERE username = ($1) AND guild_id = ($2)", username, guildID)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
