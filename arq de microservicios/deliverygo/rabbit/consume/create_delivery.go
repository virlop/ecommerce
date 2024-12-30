package consume

import (
	"encoding/json"
	"time"

	"deliverygo/events"
	"deliverygo/tools/env"
	"deliverygo/tools/log"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

// Mensaje recibido para crear un Delivery
type CreateDeliveryMessage struct {
	CorrelationId string `json:"correlation_id"`
	OrderId       string `json:"order_id"`
	UserId        string `json:"user_id"`
}

// consumeCreateDelivery escucha mensajes para la creación de un delivery
func consumeCreateDelivery() error {
	logger := log.Get().
		WithField(log.LOG_FIELD_CONTROLLER, "Rabbit").
		WithField(log.LOG_FIELD_RABBIT_EXCHANGE, "delivery").
		WithField(log.LOG_FIELD_RABBIT_QUEUE, "create_delivery").
		WithField(log.LOG_FIELD_RABBIT_ACTION, "Consume")

	// Conexión a RabbitMQ
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		logger.Error("Error al conectar con RabbitMQ: ", err)
		return err
	}
	defer conn.Close()

	// Canal
	chn, err := conn.Channel()
	if err != nil {
		logger.Error("Error al crear el canal: ", err)
		return err
	}
	defer chn.Close()

	// Declarar el Exchange
	err = chn.ExchangeDeclare(
		"delivery", // Nombre del exchange
		"direct",   // Tipo
		true,       // Durable
		false,      // Auto-delete
		false,      // Interno
		false,      // No-wait
		nil,        // Args
	)
	if err != nil {
		logger.Error("Error al declarar el exchange: ", err)
		return err
	}

	// Declarar la Cola
	queue, err := chn.QueueDeclare(
		"create_delivery", // Nombre de la cola
		true,              // Durable
		false,             // Auto-delete
		false,             // Exclusivo
		false,             // No-wait
		nil,               // Args
	)
	if err != nil {
		logger.Error("Error al declarar la cola: ", err)
		return err
	}

	// Vincular la Cola con el Exchange
	err = chn.QueueBind(
		queue.Name,     // Nombre de la cola
		"create_order", // Routing Key
		"delivery",     // Exchange
		false,          // No-wait
		nil,            // Args
	)
	if err != nil {
		logger.Error("Error al vincular la cola: ", err)
		return err
	}

	// Consumir Mensajes
	mgs, err := chn.Consume(
		queue.Name, // Nombre de la cola
		"",         // Consumidor
		false,      // Auto-ack
		false,      // Exclusivo
		false,      // No-local
		false,      // No-wait
		nil,        // Args
	)
	if err != nil {
		logger.Error("Error al consumir mensajes: ", err)
		return err
	}

	logger.Info("RabbitMQ conectado para create_delivery")

	// Procesar Mensajes
	go func() {
		for d := range mgs {
			newMessage := &CreateDeliveryMessage{}
			err := json.Unmarshal(d.Body, newMessage)
			if err != nil {
				logger.Error("Error al deserializar mensaje: ", err)
				continue
			}

			// Procesar el mensaje
			processCreateDelivery(newMessage, logger)

			// Confirmar el mensaje (ACK)
			if err := d.Ack(false); err != nil {
				logger.Error("Error al confirmar mensaje: ", err)
			} else {
				logger.Info("Mensaje procesado correctamente: ", string(d.Body))
			}
		}
	}()

	logger.Info("Conexión cerrada: ", <-conn.NotifyClose(make(chan *amqp.Error)))
	return nil
}

// processCreateDelivery maneja la lógica para crear un delivery
func processCreateDelivery(newMessage *CreateDeliveryMessage, logger *log.Entry) {
	logger.Info("Procesando mensaje de creación de delivery")

	// Crear el Delivery
	delivery := &events.Delivery{
		OrderId: newMessage.OrderId,
		UserId:  newMessage.UserId,
		Status:  "CONFIRMED",
		Created: time.Now(),
	}

	err := events.SaveDelivery(delivery)
	if err != nil {
		logger.Error("Error al guardar el delivery: ", err)
		return
	}

	logger.Info("Delivery creado para la orden: ", newMessage.OrderId)
}

func ConsumeOrderCreatedEvents(ctx ...interface{}) {
    msgs := rabbit.GetChannel(ctx...).Consume(
        "order_payment_defined_queue", // Nombre de la cola
        "",                            // Nombre del consumidor
        true,                          // Auto-Ack
        false,
        false,
        false,
        nil,
    )

    for msg := range msgs {
        eventData := struct {
            DeliveryId string `json:"deliveryId"`
            OrderId    string `json:"orderId"`
            UserId     string `json:"userId"`
        }{}

        if err := json.Unmarshal(msg.Body, &eventData); err != nil {
            log.Get(ctx...).Error("Failed to unmarshal message: ", err)
            continue
        }

        event := NewConfirmDeliveryEvent(eventData.DeliveryId, eventData.OrderId, eventData.UserId)
        if _, err := InsertDeliveryEvent(event, ctx...); err != nil {
            log.Get(ctx...).Error("Failed to insert delivery event: ", err)
        }
    }
}
