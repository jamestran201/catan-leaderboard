package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"os"
	"os/signal"
	"syscall"
)

var dbConn *pgx.Conn

func init() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error preparing database connection for migration: ", err)
		os.Exit(1)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("Error preparing database driver for migration: ", err)
		os.Exit(1)
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		fmt.Println("Error during migration: ", err)
		os.Exit(1)
	}

	migration.Up()
	err = db.Close()
	if err != nil {
		fmt.Println("Error closing connection after migration: ", err)
		os.Exit(1)
	}

	dbConn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
		os.Exit(1)
	}
}

func main() {
	defer dbConn.Close(context.Background())

	token := os.Getenv("BOT_TOKEN")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		os.Exit(1)
	}
	defer discord.Close()

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		os.Exit(1)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == session.State.User.ID {
		return
	}

	messageSender := &DiscordMessageSender{session, m}
	messageParser := &DiscordMessageParser{discordMessage: m}
	db := &PostgresDataLayer{dbConn}
	bot := CatanBot{session, m, nil, messageSender, messageParser, db}
	bot.handleCommand()
}
