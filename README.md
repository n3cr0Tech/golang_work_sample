## Written as base code for running server with numerous concurrent players with realtime interactions


## Features
- Player registration/login
- JWT auth tokens
- MongoDB module
- RabbitMQ module
- SocketIO module

## Setup / Quickstart
- cp env_template .env
- docker compose up -d
- cd main
- go get .
- go run .

## Relevant commands:
- go run .
    - run program in baseS dir
- go get .
    - download all imported dependencies
- go mod tidy
    - ensures that go.mod matches the source code in the module, adds any missing modules (similar to npm install)
    - creates a go.sum file
- go get LIBRARY_NAME
    - installs a library/package
- go mod init golang_work_sample/internal/PACKAGE-NAME
    - create a new module


## Connecting to MongoDB
- mongodb://USER:PWD@localhost:27017/admin