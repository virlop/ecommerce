package events

//UPDATES DE LOS DELIVERIES
import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewConfirmDeliveryEvent crea un nuevo evento de confirmación
func NewConfirmDeliveryEvent(deliveryId, orderId, userId string) *Event {
	return &Event{
		ID:             primitive.NewObjectID(),
		DeliveryId:     deliveryId,
		OrderId:        orderId,
		DeliveryStatus: DeliveryStatusConfirmed, // Estado inicial del Delivery al crearse
		Type:           ConfirmDelivery,
		ConfirmDelivery: &ConfirmDeliveryEvent{
			Timestamp: time.Now(),
		},
		Created: time.Now(),
	}
}

func NewCancelledDeliveryEvent(deliveryId, orderId, userId string, deps ...interface{}) (*Event, error) {
	// Consultar el estado actual del Delivery
	events, err := FindDeliveryEventsByDeliveryId(deliveryId, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delivery events: %w", err)
	}

	// Verificar si el estado anterior es válido para la cancelación
	var currentStatus DeliveryStatus
	if len(events) > 0 {
		latestEvent := events[len(events)-1] // Último evento
		currentStatus = latestEvent.DeliveryStatus
	} else {
		return nil, fmt.Errorf("no events found for deliveryId: %s", deliveryId)
	}

	if currentStatus != DeliveryStatusConfirmed && currentStatus != DeliveryStatusOnTheGo {
		return nil, fmt.Errorf("cannot cancel delivery with current status: %s", currentStatus)
	}

	// Crear y devolver el evento de cancelación
	return &Event{
		ID:             primitive.NewObjectID(),
		DeliveryId:     deliveryId,
		OrderId:        orderId,
		DeliveryStatus: DeliveryStatusCancelled, // Actualiza el estado a 'cancelled'
		Type:           CancelledDelivery,
		CancelledDelivery: &CancelledDeliveryEvent{
			UserId:    userId,
			Timestamp: time.Now(),
		},
		Created: time.Now(),
	}, nil
}

func NewSetOnTheGoDeliveryEvent(deliveryId, orderId, userId string, deps ...interface{}) (*Event, error) {
	// Consultar los eventos relacionados con el Delivery para verificar el estado actual
	events, err := FindDeliveryEventsByDeliveryId(deliveryId, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delivery events: %w", err)
	}

	// Validar si el último estado es 'confirmed'
	var currentStatus DeliveryStatus
	if len(events) > 0 {
		latestEvent := events[len(events)-1] // Último evento registrado
		currentStatus = latestEvent.DeliveryStatus
	} else {
		return nil, fmt.Errorf("no events found for deliveryId: %s", deliveryId)
	}

	if currentStatus != DeliveryStatusConfirmed {
		return nil, fmt.Errorf("cannot set delivery to on_the_go with current status: %s", currentStatus)
	}

	// Crear el evento con el nuevo estado
	return &Event{
		ID:             primitive.NewObjectID(),
		DeliveryId:     deliveryId,
		OrderId:        orderId,
		DeliveryStatus: DeliveryStatusOnTheGo, // Cambiar el estado a 'on_the_go'
		Type:           SetOnTheGoDelivery,
		SetOnTheGoDelivery: &SetOnTheGoDeliveryEvent{
			UserId:    userId,
			Timestamp: time.Now(),
		},
		Created: time.Now(),
	}, nil
}

func NewSetDeliveredDeliveryEvent(deliveryId, orderId, userId string, deps ...interface{}) (*Event, error) {
	// Consultar los eventos relacionados con el Delivery para verificar el estado actual
	events, err := FindDeliveryEventsByDeliveryId(deliveryId, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delivery events: %w", err)
	}

	// Validar si el último estado es 'on_the_go'
	var currentStatus DeliveryStatus
	if len(events) > 0 {
		latestEvent := events[len(events)-1] // Último evento registrado
		currentStatus = latestEvent.DeliveryStatus
	} else {
		return nil, fmt.Errorf("no events found for deliveryId: %s", deliveryId)
	}

	if currentStatus != DeliveryStatusOnTheGo {
		return nil, fmt.Errorf("cannot set delivery to delivered with current status: %s", currentStatus)
	}

	// Crear el evento con el nuevo estado
	return &Event{
		ID:             primitive.NewObjectID(),
		DeliveryId:     deliveryId,
		OrderId:        orderId,
		DeliveryStatus: DeliveryStatusDelivered, // Cambiar el estado a 'delivered'
		Type:           SetDeliveredDelivery,
		SetDeliveredDelivery: &SetDeliveredDeliveryEvent{
			UserId:    userId,
			Timestamp: time.Now(),
		},
		Created: time.Now(),
	}, nil
}
