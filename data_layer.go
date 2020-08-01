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
}

type postgresDataLayer struct {
	dbConn *pgx.Conn
}

func (db *postgresDataLayer) addUser(username string, guildID string) error {
	_, err := db.dbConn.Exec(context.Background(), "INSERT INTO users (username, guild_id) VALUES ($1, $2)", username, guildID)

	return err
}

func (db *postgresDataLayer) getTopTwentyUsers(guildID string) ([]user, error) {
	rows, err := db.dbConn.Query(context.Background(), "SELECT CAST(RANK() OVER (ORDER BY games_won DESC) AS TEXT), username, CAST(games_won AS TEXT) FROM users WHERE guild_id = ($1) LIMIT 20", guildID)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	users := make([]user, 0, 20)

	for i := 0; rows.Next(); i++ {
		user := user{}
		err = rows.Scan(&user.rank, &user.username, &user.victories)
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
