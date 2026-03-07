package socketio

import (
	"encoding/json"
	"fmt"
	"log"	

	socketio "github.com/googollee/go-socket.io"
)

var SocketIOServer *socketio.Server


type testMsgPayload struct {
	name string `json:"name"`
	msg string `json:"msg"`
}


func Init(){
	log.Println("Inside socketio.Init()")
	SocketIOServer = socketio.NewServer(nil)

	SocketIOServer.OnConnect("/", func(s socketio.Conn) error{
		s.SetContext("")
		fmt.Println("Server connected: ", s.ID())
		return nil
	})

	// Make sure that the Client's payload is a stringified JSON
	SocketIOServer.OnEvent("/", "testMsg", func(s socketio.Conn, payload string) string{
		s.SetContext(payload)
		fmt.Println("testMsg received: ", payload)		

		var jsonObj map[string]string
		err := json.Unmarshal([]byte(payload), &jsonObj)
		if err != nil {
		 fmt.Println(err)		 
		}
	   
		fmt.Println("name: ", jsonObj["name"])
		replyToTestMsg(s)
		return payload
	})

	
	SocketIOServer.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	SocketIOServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("SocketID: ", s.ID(), " has disconnected")
		log.Println("closed", reason)
	})
}

func replyToTestMsg(s socketio.Conn){
	payload := map[string]string{"favoriteMove": "Star Wars", "msg":"Hello there"}
	s.Emit("replyTestMsg", payload)

	
}