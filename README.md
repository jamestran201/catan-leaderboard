# Catan Leaderboard

## Requirements

1. Go version >= 1.14

2. Postgres database

3. Set the environment variable `DATABASE_URL` to `postgres://[user]:[password]@[hostname]:[port]/[dbname]`

4. Run `go get` to install dependencies

5. Set up Discord bot token

6. Give the bot permissions to `Send messages` and `Read message history`

7. Set the environment variable `BOT_TOKEN` to the Discord bot token

## Bot commands

The commands of this bot follow the format: `catan! [command] [arguments]`

| Command                           | Description                                                                        |
|-----------------------------------|------------------------------------------------------------------------------------|
| catan! adduser [username]         | Add a user to the leaderboard                                                      |
| catan! addwin [username]          | Add a win for the user                                                             |
| catan! record [username] [points] | Add points after a game for the user. This also updates the points per game column |
| catan! leaderboard                | Display the leaderboard                                                            |

![leaderboard screenshot](https://i.imgur.com/2o6sYrb.png)

## Running tests

1. Create a Postgres database for testing
2. Set the environment variable `TEST_DATABASE_URL` to `postgres://[user]:[password]@[hostname]:[port]/[dbname]`
3. Run `go test`
