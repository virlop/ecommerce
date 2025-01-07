// Define la estructura principal que representa las proyecciones de Delivery. Estas proyecciones son datos preprocesados, 
// optimizados para consultas rápidas y consistentes.
// Incluye validaciones para asegurar que los datos sean correctos antes de insertarlos o actualizarlos en la base de datos.
package projections

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeliveryProjection representa la proyección de un Delivery
type DeliveryProjection struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	DeliveryId   string             `bson:"deliveryId" validate:"required"`
	OrderId      string             `bson:"orderId" validate:"required"`
	UserId       string             `bson:"userId" validate:"required"`
	Status       string             `bson:"status" validate:"required"`
	CreatedAt    time.Time          `bson:"createdAt"`
	LastModified time.Time          `bson:"lastModified"`
}

// Validate valida la estructura
func (dp *DeliveryProjection) Validate() error {
	return validator.New().Struct(dp)
}
