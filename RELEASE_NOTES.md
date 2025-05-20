## 0.9.0
- Added functionality to utils to be able to read json files (for configs, etc.)

## 0.8.0
- MAJOR BUG FIX: turns out the API and the RabbitMQ were not being served simultaneously.
    - Removed the "<- forever" syntax that locked the main thread in rabbitmq module

## 0.7.0
- MAJOR BUG FIX: turns out the API and the SocketIO were not being served simultaneously. Fixed it by wrapping the socketio handlers with gin.WrapH() and using the router

## 0.6.0
- Refactored GetRecord() will now return an error if the requested record is not found on the DB
    - if record found, it returns a JSON string of the record
- Added EnsureUpsertRecord (to prevent overwriting existing records)
- Added /register endpoint with mongodb integration
- /login endpoint will now reply with a failure msg
- Refactored mongodb module to separate UpsertRecord() away from the /register code path

## 0.5.0
- BUG FIX: Server now able to listen to rabbitMQ AND serve API simultaneously
- Created MongoDB module
    - Creates collections
    - Upserts Record
    - Gets Record
    - Deletes Record

## 0.4.1
- Server now consumes rabbitmq messages
    - demos callback for printing incoming rabbitmq messages

## 0.4.0
- Added rabbitMQ module
    - Demos sending messages in RabbitMQ
    - explicit function for sending JSON
- TODO: rabbitmq needs to be able to consume messages
    
## 0.3.0
- Added socketIO module
    - Demos socketIO comms
    - Reading/Sending JSON payloads 

## 0.2.0
- .env file added
- created middleware
- created Login endpoint
    - jwt tokens sent out to successful logins    
    - mockData creates mock Users (bcrypt encrypted mock passwords)
- protected the album routes with middleware

## 0.1.0
- Created mockData package
- project demos basic API GET and POST endpoints

## 0.0.1
- project demos calling function from another package
- greetings package uses rsc.io/quote library