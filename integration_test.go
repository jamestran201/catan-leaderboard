package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4"
)

var testDbConn *pgx.Conn

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	db, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		fmt.Println("Error preparing database connection for migration: ", err)
		return 1
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("Error preparing database driver for migration: ", err)
		return 1
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		fmt.Println("Error during migration: ", err)
		return 1
	}

	migration.Up()
	err = db.Close()
	if err != nil {
		fmt.Println("Error closing connection after migration: ", err)
		return 1
	}

	testDbConn, err = pgx.Connect(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
		return 1
	}

	defer testDbConn.Close(context.Background())

	return m.Run()
}

func registerCleanup(t *testing.T) {
	t.Cleanup(func() {
		_, err := testDbConn.Exec(context.Background(), "TRUNCATE TABLE users;")

		if err != nil {
			t.Error("Error cleaning up after test: ", err)
		}
	})
}

func TestAddUser(t *testing.T) {
	registerCleanup(t)

	message := &discordgo.Message{Content: "catan! adduser hinata", GuildID: "1"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := &MessageSenderMock{}
	messageParser := &discordMessageParser{discordMessage: messageCreate}
	db := &postgresDataLayer{testDbConn}
	bot := catanBot{sender, messageParser, db}
	bot.handleCommand()

	expected := "Successfully added user: hinata"
	if sender.messageSent != expected {
		t.Errorf("\nGot: %s\nExpect: %s\n", sender.messageSent, expected)
	}
}

func TestAddWin(t *testing.T) {
	registerCleanup(t)

	db := &postgresDataLayer{testDbConn}
	db.addUser("kageyama", "1")

	message := &discordgo.Message{Content: "catan! addwin kageyama", GuildID: "1"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := &MessageSenderMock{}
	messageParser := &discordMessageParser{discordMessage: messageCreate}
	bot := catanBot{sender, messageParser, db}
	bot.handleCommand()

	expectedLeaderboard := `+------+----------+-----------+
| RANK | USERNAME | VICTORIES |
+------+----------+-----------+
|    1 | kageyama |         1 |
+------+----------+-----------+
`
	expectedMessage := fmt.Sprintf("Congrats kageyama on the win! :tada:\n```%s```", expectedLeaderboard)
	if sender.messageSent != expectedMessage {
		t.Errorf("\nGot: %s\nExpect: %s\n", sender.messageSent, expectedMessage)
	}
}

func TestShowLeaderboard(t *testing.T) {
	registerCleanup(t)

	db := &postgresDataLayer{testDbConn}
	db.addUser("kageyama", "1")
	db.addUser("oikawa", "1")

	db.addWin("kageyama", "1")
	db.addWin("oikawa", "1")
	db.addWin("oikawa", "1")

	message := &discordgo.Message{Content: "catan! leaderboard", GuildID: "1"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := &MessageSenderMock{}
	messageParser := &discordMessageParser{discordMessage: messageCreate}
	bot := catanBot{sender, messageParser, db}
	bot.handleCommand()

	expectedLeaderboard := `+------+----------+-----------+
| RANK | USERNAME | VICTORIES |
+------+----------+-----------+
|    1 | oikawa   |         2 |
|    2 | kageyama |         1 |
+------+----------+-----------+
`
	expectedValue := fmt.Sprintf("```%s```", expectedLeaderboard)

	if sender.messageSent != expectedValue {
		t.Errorf("\nGot: %s\nExpect: %s\n", sender.messageSent, expectedValue)
	}
}

type MessageSenderMock struct {
	messageSent  string
	messageEmbed *discordgo.MessageEmbed
}

func (sender *MessageSenderMock) sendMessage(message string) {
	sender.messageSent = message
}

func (sender *MessageSenderMock) sendEmbedMessage(embed *discordgo.MessageEmbed) {
	sender.messageEmbed = embed
}
