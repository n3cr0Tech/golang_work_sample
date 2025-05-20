package main

import (
	"fmt"
	"log"	
	"github.com/gin-gonic/gin"	

	auth "example.com/auth"
	middleware "example.com/middleware"
	mockData "example.com/mockData"
	mongodb "example.com/mongodb"
	rabbitmq "example.com/rabbitmq"
	"example.com/socketio"
	// socketio "github.com/googollee/go-socket.io"
	"example.com/types"
	utils "example.com/utils"		
)

var router *gin.Engine

func main() {
	readConfigFile()

	router = gin.Default()
	mockData.CreateMockUsers()
	utils.GetEnvEntries()
	USERS_COLLECTIONS_NAME := "users"
	mongodb.Init(utils.EnvEntries["MONGO_URL"], utils.EnvEntries["DB_NAME"], USERS_COLLECTIONS_NAME)
	collectionsName := "test"
	testRecordName := "player-machine-4-lab"
	testCreateDBData(collectionsName, testRecordName) //DEBUG
	testGetDBData(collectionsName, testRecordName)//DEBUG
	testDeleteDBData(collectionsName, testRecordName)//DEBUG
	
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

	router.GET("/socket.io/*any", gin.WrapH(socketio.SocketIOServer))
	router.POST("/socket.io/*any", gin.WrapH(socketio.SocketIOServer))
	router.GET("/albums", middleware.AuthChecker, mockData.GetAlbums)
	router.GET("/albums/:id", middleware.AuthChecker, mockData.GetAlbum)
	router.POST("/greeting", mockData.PostGreeting)
	router.POST("/login", auth.Login)
	router.POST("/register", auth.Register)
	
		
	log.Println("API Server at localhost:8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatal("failed run app: ", err)
	}

	rabbitmq.ListenToQueue(rabbitMQCallback, rabbitmq.SERVER_QNAME)
	fmt.Println("Called rabbitmq.Listen()")

}


func rabbitMQCallback(msg string, qName string){
	  log.Printf(" %s Received a message: %s", qName, msg)
}

func sendTestJSONMsg(){
	payload := types.TestPayload{
		Msg: "test123",
		Data: types.TestPayloadChild{
			MockCounter: 96,
			MockId: 1337,
		},
	}
	rabbitmq.SendJSONMsg(payload, rabbitmq.CLIENT_QNAME)
}

// Demoes rabbitmq.SendJSONMsg() can take any struct type
func sendTestJSONMsgB(){
	payload := types.TestPayloadB{
		Foo: "hello-there",
		Blah: 72,
	}
	rabbitmq.SendJSONMsg(payload, rabbitmq.CLIENT_QNAME)
}

func testGetDBData(collections string, recordName string){		
	data := map[string]interface{}{
		"name": recordName,
	}	
	mongodb.GetRecord(collections, data)
}

func testCreateDBData(collections string, recordName string){	
	mongodb.CreateCollections(collections)
	filter := map[string]interface{}{"uuid": "abc-0"}
	data := map[string]any{		
		"name": recordName,
		"points": 123,
	}	
	mongodb.UpsertRecord(collections, filter, data)
}

func testDeleteDBData(collections string, recordName string){	
	mongodb.DeleteRecord(collections, "name", recordName)
}

func readConfigFile(){
	config := utils.GetConfigFromJSON("", "config.json")	
	log.Println("- config: ")
	log.Println(config)
	tmp := config["foo-entry"]
	log.Println("test value: " + tmp)
}