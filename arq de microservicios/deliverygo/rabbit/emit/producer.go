package rabbit

import (
	"deliverygo/rabbit"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// PublishMessage publica un mensaje en un exchange
func PublishMessage(exchange, routingKey string, body interface{}) error {
	channel := rabbit.GetChannel()

	// Serializa el mensaje en JSON
	message, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = channel.Publish(
		exchange,   // Nombre del exchange
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Printf("Error al publicar mensaje: %v", err)
		return err
	}

	log.Printf("Mensaje publicado: %s", message)
	return nil
}
