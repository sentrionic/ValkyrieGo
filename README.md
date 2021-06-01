# ValkyrieGO

<p align="center">
  <img src="https://harmony-cdn.s3.eu-central-1.amazonaws.com/logo.png">
</p>

A [Discord](https://discord.com) backend clone written in Golang.
(For the original Typescript version including the client see [Valkyrie](https://github.com/sentrionic/Valkyrie))


## Features

- Message, Channel, Server CRUD
- Authentication using Sessions
- Channel / Websocket Member Protection
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
- S3 for storing files and GMail for sending emails
- [React Client](https://github.com/sentrionic/Valkyrie/tree/websocket)
- [Flutter Application](https://github.com/sentrionic/ValkyrieApp/tree/websocket)
---

## Installation

### Server

1. Install PostgreSQL and create a DB
2. Install Redis
3. Install Golang and get all the dependencies
4. Rename `.env.example` to `.env` and fill in the values

- `Required`

        PORT=8080
        DATABASE_URL="postgresql://<username>:<password>@localhost:5432/db_name"
        REDIS_URL=localhost:6379
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

Once the server is running go to `localhost:8080/swagger/index.html` to see all the HTTP endpoints
and `localhost:8080` for all the websocket events.
