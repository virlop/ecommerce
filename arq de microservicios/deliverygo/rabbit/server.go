package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

var connection *amqp.Connection
var channel *amqp.Channel

// Init inicializa la conexión a RabbitMQ
func Init(rabbitURL string) {
	var err error

	// Establece la conexión
	connection, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}

	// Crea un canal
	channel, err = connection.Channel()
	if err != nil {
		log.Fatalf("Error al crear el canal RabbitMQ: %v", err)
	}
}

// GetChannel retorna el canal actual
func GetChannel() *amqp.Channel {
	return channel
}

// Close cierra la conexión y el canal
func Close() {
	if channel != nil {
		channel.Close()
	}
	if connection != nil {
		connection.Close()
	}
}
