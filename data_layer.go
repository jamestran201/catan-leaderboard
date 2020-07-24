package main

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DataLayer interface {
	AddUser(username string, guild_id string) error
	GetTopFiveUsers(guild_id string) ([]User, error)
}

type PostgresDataLayer struct {
	dbConn *pgx.Conn
}

func (db *PostgresDataLayer) AddUser(username string, guild_id string) error {
	_, err := db.dbConn.Exec(context.Background(), "INSERT INTO users (username, guild_id) VALUES ($1, $2)", username, guild_id)

	return err
}

func (db *PostgresDataLayer) GetTopFiveUsers(guild_id string) ([]User, error) {
	rows, err := db.dbConn.Query(context.Background(), "SELECT CAST(RANK() OVER (ORDER BY games_won DESC) AS TEXT), username, CAST(games_won AS TEXT) FROM users WHERE guild_id = ($1) LIMIT 5", guild_id)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	users := make([]User, 0, 5)

	for i := 0; rows.Next(); i++ {
		user := User{}
		err = rows.Scan(&user.Rank, &user.Username, &user.Victories)
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
