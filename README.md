## Written as base code for running server with numerous concurrent players with realtime interactions

## Features
- Player registration/login
- JWT auth tokens
- MongoDB module
- RabbitMQ module
- SocketIO module

## .env entries:
- JWT_SECRET
- AUTH_HEADER

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
- go mod init example.com/PACKAGE-NAME
    - create a new module

## Create a new module:
- mkdir new-folder-name
- cd new-folder-name
- go mod init util-module  (typically same name as folder)
- then in the folder that will use your new module (eg 'main'), execute commands:
    - go mod edit -replace example.com/util-module=../util-module
    - go mod tidy
