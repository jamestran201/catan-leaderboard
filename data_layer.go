package main

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DataLayer interface {
	AddUser(username string, guild_id string) error
}

type PostgresDataLayer struct {
	dbConn *pgx.Conn
}

func (db *PostgresDataLayer) AddUser(username string, guild_id string) error {
	_, err := db.dbConn.Exec(context.Background(), "INSERT INTO users (username, guild_id) VALUES ($1, $2)", username, guild_id)

	return err
}
