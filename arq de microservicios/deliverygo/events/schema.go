// estructuras y métodos necesarios para manejar eventos

package events

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeliveryStatus
type DeliveryStatus string

const (
	DeliveryStatusConfirmed DeliveryStatus = "confirmed"
	DeliveryStatusCancelled DeliveryStatus = "cancelled"
	DeliveryStatusOnTheGo   DeliveryStatus = "on_the_go"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
)

func (ds DeliveryStatus) IsValid() bool {
	switch ds {
	case DeliveryStatusConfirmed, DeliveryStatusCancelled, DeliveryStatusOnTheGo, DeliveryStatusDelivered:
		return true
	}
	return false
}

// EventType define los tipos de eventos manejados
type EventType string

const (
	ConfirmDelivery      EventType = "confirm_delivery"
	CancelledDelivery    EventType = "cancelled_delivery"
	SetOnTheGoDelivery   EventType = "set_onthego_delivery"
	SetDeliveredDelivery EventType = "set_delivered_delivery"
)

func (et EventType) IsValid() bool {
	switch et {
	case ConfirmDelivery, CancelledDelivery, SetOnTheGoDelivery, SetDeliveredDelivery:
		return true
	}
	return false
}

// Event representa un evento de RabbitMQ
type Event struct {
	ID                   primitive.ObjectID         `bson:"_id,omitempty"`                  // ID generado por MongoDB
	DeliveryId           string                     `bson:"deliveryId" validate:"required"` // ID del delivery
	OrderId              string                     `bson:"orderId" validate:"required"`    // ID de la orden asociada
	DeliveryStatus       DeliveryStatus             `bson:"deliveryStatus" validate:"required"`
	Type                 EventType                  `bson:"type" validate:"required"` // Tipo de evento
	ConfirmDelivery      *ConfirmDeliveryEvent      `bson:"confirmDeliveryEvent"`     // Datos del evento específico
	CancelledDelivery    *CancelledDeliveryEvent    `bson:"cancelledDeliveryEvent"`
	SetOnTheGoDelivery   *SetOnTheGoDeliveryEvent   `bson:"setOnTheGoDeliveryEvent"`
	SetDeliveredDelivery *SetDeliveredDeliveryEvent `bson:"setDeliveredDeliveryEvent"`
	Created              time.Time                  `bson:"created"` // Fecha de creación del evento
}

// ValidateSchema valida que los datos del evento sean correctos antes de insertarlo
func (e *Event) ValidateSchema() error {
	if err := validator.New().Struct(e); err != nil {
		return err
	}

	if !e.DeliveryStatus.IsValid() {
		return fmt.Errorf("invalid delivery status: %s", e.DeliveryStatus)
	}
	if !e.Type.IsValid() {
		return fmt.Errorf("invalid event type: %s", e.Type)
	}

	return nil
}

// ConfirmDeliveryEvent define los datos específicos de confirmación
type ConfirmDeliveryEvent struct {
	Timestamp time.Time `bson:"timestamp"` // Fecha de la confirmación
}
type CancelledDeliveryEvent struct {
	UserId    string    `bson:"userId" validate:"required"` // ID del usuario que cancela
	Timestamp time.Time `bson:"timestamp"`                  // Fecha del cambio
}

type SetOnTheGoDeliveryEvent struct {
	UserId    string    `bson:"userId" validate:"required"` // ID del usuario que cambia a "on the go"
	Timestamp time.Time `bson:"timestamp"`                  // Fecha del cambio
}

type SetDeliveredDeliveryEvent struct {
	UserId    string    `bson:"userId" validate:"required"` // ID del usuario que marca como entregado
	Timestamp time.Time `bson:"timestamp"`                  // Fecha del cambio
}
