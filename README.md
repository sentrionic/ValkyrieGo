[![Go Report Card](https://goreportcard.com/badge/github.com/sentrionic/ValkyrieGo)](https://goreportcard.com/report/github.com/sentrionic/ValkyrieGo)

This repository got merged into the main repository and only exists to view the commit history prior to the merge.

For the newest and most up to date version go to [Valkyrie](https://github.com/sentrionic/Valkyrie).

# ValkyrieGO

<p align="center">
  <img src="https://harmony-cdn.s3.eu-central-1.amazonaws.com/logo.png" alt="Valkyrie Logo">
</p>

A [Discord](https://discord.com) backend clone written in Golang.
(For the original Typescript version including the client see [Valkyrie](https://github.com/sentrionic/Valkyrie))


## Features

- Message, Channel, Server CRUD
- Authentication using Sessions
- Channel / Websockets Member Protection
- Realtime Events
- File Upload (Avatar, Icon, Messages) to S3
- Direct Messaging
- Private Channels
- Friend System
- Notification System
- Basic Moderation for the guild owner (delete messages, kick & ban members)

## Stack

- [Gin](https://gin-gonic.com/) for the HTTP server
- [Gorilla Websockets](https://github.com/gorilla/websocket) for WS communication
- [Gorm](https://gorm.io/) as the database ORM
- PostgreSQL
- Redis
- S3 for storing files and Gmail for sending emails
- [React Client](https://github.com/sentrionic/Valkyrie/tree/websocket)
- [Flutter Application](https://github.com/sentrionic/ValkyrieApp/tree/websocket)
---

## Installation
If you are familiar with `make`, take a look at the `Makefile` to quickly setup the following steps
or alternatively copy the commands into your CLI.

1. Install Docker and get the Postgresql and Redis containers (`make postgres` && `make redis`)
2. Start both containers (`make start`) and create a DB (`make createdb`)
3. Install Golang and get all the dependencies (`go mod tidy`)
4. Rename `.env.example` to `.env` and fill in the values

- `Required`

        PORT=4000
        DATABASE_URL=postgresql://<username>:<password>@localhost:5432/valkyrie
        REDIS_URL=redis://localhost:6379
        CORS_ORIGIN=http://localhost:3000
        SECRET=SUPERSECRET
        HANDLER_TIMEOUT=5
        MAX_BODY_BYTES=4194304 # 4MB in Bytes = 4 * 1024 * 1024

- `Optional: Not needed to run the app, but you won't be able to upload files or send emails.`

        AWS_ACCESS_KEY=ACCESS_KEY
        AWS_SECRET_ACCESS_KEY=SECRET_ACCESS_KEY
        AWS_STORAGE_BUCKET_NAME=STORAGE_BUCKET_NAME
        AWS_S3_REGION=S3_REGION
        GMAIL_USER=GMAIL_USER
        GMAIL_PASSWORD=GMAIL_PASSWORD

5. Run `go run github.com/sentrionic/valkyrie` to run the server

## Endpoints

Once the server is running go to `localhost:<PORT>/swagger/index.html` to see all the HTTP endpoints
and `localhost:<PORT>` for all the websockets events.

## Tests
All routes in `handler` have tests written for them.

Only services that do not just delegate work to the repository have tests written for them.

Run `go test -v -cover ./service/... ./handler/...` (`make test`) to run all tests

Additionally this repository includes E2E tests for all successful requests. To run them you
have to have Postgres and Redis running in Docker and then run `go test github.com/sentrionic/valkyrie` (`make e2e`).

## Credits
[Jacob Goodwin](https://github.com/JacobSNGoodwin/memrizr): This backend is built upon his tutorial series and uses his backend structure

[Jeroen de Kok](https://dev.to/jeroendk/building-a-simple-chat-application-with-websockets-in-go-and-vue-js-gao): The websockets structure is based on his tutorial
