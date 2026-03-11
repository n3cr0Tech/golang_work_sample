package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	auth "golang_work_sample/internal/auth"
	middleware "golang_work_sample/internal/middleware"
	mockData "golang_work_sample/internal/mockData"
	mongodb "golang_work_sample/internal/mongodb"
	rabbitmq "golang_work_sample/internal/rabbitmq"
	socketio "golang_work_sample/internal/socketio"

	// socketio "github.com/googollee/go-socket.io"
	types "golang_work_sample/internal/types"
	utils "golang_work_sample/internal/utils"
)

var router *gin.Engine

func main() {
	readConfigFile()

	router = gin.Default()
	// mockData.CreateMockUsers()
	utils.GetEnvEntries()
	USERS_COLLECTIONS_NAME := utils.EnvEntries["MONGO_USERS_DB"]

	// Init MongoDB
	mongoClient, err := mongodb.NewClient(context.Background(), utils.EnvEntries["MONGO_URL"], utils.EnvEntries["MONGO_USER"], utils.EnvEntries["MONGO_PWD"], utils.EnvEntries["DB_NAME"])
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.CreateCollections(context.Background(), USERS_COLLECTIONS_NAME); err != nil {
		log.Printf("Warning: failed to ensure collections: %v", err)
	}

	// Testing DB
	collectionsName := "test"
	testRecordName := "player-machine-4-lab"
	testCreateDBData(mongoClient, collectionsName, testRecordName) //DEBUG
	testGetDBData(mongoClient, collectionsName, testRecordName)    //DEBUG
	testDeleteDBData(mongoClient, collectionsName, testRecordName) //DEBUG

	// Init RabbitMQ
	rabbitmq.Init(utils.GetRabbitMQURL())
	log.Println("rabbitmq INIT done")
	sendTestJSONMsg()
	sendTestJSONMsgB()

	//-- SocketIO --
	socketio.Init()
	go func() {
		if err := socketio.SocketIOServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer socketio.SocketIOServer.Close()
	//---------------

	// Init handlers
	authHandler := auth.NewAuthHandler(mongoClient)
	mwHandler := middleware.NewMiddlewareHandler(mongoClient)

	router.GET("/socket.io/*any", gin.WrapH(socketio.SocketIOServer))
	router.POST("/socket.io/*any", gin.WrapH(socketio.SocketIOServer))
	router.GET("/albums", mwHandler.AuthChecker, mockData.GetAlbums)
	router.GET("/albums/:id", mwHandler.AuthChecker, mockData.GetAlbum)
	router.POST("/greeting", mockData.PostGreeting)
	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)

	log.Println("API Server at localhost:8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatal("failed run app: ", err)
	}

	rabbitmq.ListenToQueue(rabbitMQCallback, rabbitmq.SERVER_QNAME)
	fmt.Println("Called rabbitmq.Listen()")

}

func rabbitMQCallback(msg string, qName string) {
	log.Printf(" %s Received a message: %s", qName, msg)
}

func sendTestJSONMsg() {
	payload := types.TestPayload{
		Msg: "test123",
		Data: types.TestPayloadChild{
			MockCounter: 96,
			MockId:      1337,
		},
	}
	rabbitmq.SendJSONMsg(payload, rabbitmq.CLIENT_QNAME)
}

// Demoes rabbitmq.SendJSONMsg() can take any struct type
func sendTestJSONMsgB() {
	payload := types.TestPayloadB{
		Foo:  "hello-there",
		Blah: 72,
	}
	rabbitmq.SendJSONMsg(payload, rabbitmq.CLIENT_QNAME)
}

func testGetDBData(mongoClient *mongodb.MongoClient, collections string, recordName string) {
	data := map[string]interface{}{
		"name": recordName,
	}
	mongoClient.FindOne(context.Background(), collections, data)
}

func testCreateDBData(mongoClient *mongodb.MongoClient, collections string, recordName string) {
	if err := mongoClient.CreateCollections(context.Background(), collections); err != nil {
		log.Printf("Test error: %v", err)
	}
	filter := map[string]interface{}{"uuid": "abc-0"}
	data := map[string]any{
		"name":   recordName,
		"points": 123,
	}
	if _, err := mongoClient.UpsertRecord(context.Background(), collections, filter, data); err != nil {
		log.Printf("Test error: %v", err)
	}
}

func testDeleteDBData(mongoClient *mongodb.MongoClient, collections string, recordName string) {
	if err := mongoClient.DeleteRecord(context.Background(), collections, "name", recordName); err != nil {
		log.Printf("Test error: %v", err)
	}
}

func readConfigFile() {
	config := utils.GetConfigFromJSON("", "config.json")
	log.Println("- config: ")
	log.Println(config)
	tmp := config["foo-entry"]
	log.Println("test value: " + tmp)
}
