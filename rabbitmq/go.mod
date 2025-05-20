module rabbitmq

go 1.22.5

require (
	example.com/types v0.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.10.0
)

replace example.com/types => ../types
