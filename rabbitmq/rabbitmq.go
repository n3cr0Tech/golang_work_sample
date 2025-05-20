package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel
var serverQ amqp.Queue
var clientQ amqp.Queue
var ctx context.Context

var CLIENT_QNAME = "mobile-0-q"
var SERVER_QNAME = "to-server-q"

func Init(url string){
	log.Println("INSIDE rabbitmq.Init()")
	var err error
	//Establish connection
	conn, err = amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	//Create a channel
	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	serverQ = createQueue(SERVER_QNAME, ch)
	clientQ = createQueue(CLIENT_QNAME, ch)
	log.Println("creating Q's done")		

}

func SendMsg(msg string, qName string){
	log.Println("INSIDE rabbitmq.SendMsg()")	
	var cancel context.CancelFunc
	//Create a context to use for publishing a message (body)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	body := msg
	err := ch.PublishWithContext(ctx,
	"",     // exchange
	qName, // routing key
	false,  // mandatory
	false,  // immediate
	amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s\n", body)
}

func SendJSONMsg(jsonPayload interface{}, qName string){
	log.Println("INSIDE rabbitmq.SendMsg()")	
	var cancel context.CancelFunc
	//Create a context to use for publishing a message (body)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	body, _ := json.Marshal(jsonPayload)
	err := ch.PublishWithContext(ctx,
	"",     // exchange
	qName, // routing key
	false,  // mandatory
	false,  // immediate
	amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s\n", body)
}

func ListenToQueue(callback func(string, string), qName string){
	
	//consume messages from q
	msgs, err := ch.Consume(
		qName, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	
	 
	go func() {
		for d := range msgs {
			callback(string(d.Body), qName)
		}
	}()
	
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

}

func createQueue(qName string, ch *amqp.Channel) amqp.Queue{
	//Create queue
	queue, err := ch.QueueDeclare(
		qName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}