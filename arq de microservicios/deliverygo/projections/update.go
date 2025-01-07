// Contiene la lógica de negocio que actualiza las proyecciones en base a los eventos recibidos.
// Aplica cambios en los datos procesados dependiendo del tipo de evento (confirmación, cancelación, etc.).
package projections

import (
	"github.com/virlop/ecommerce-ams/deliverygo/events"
)

func Update(deliveryId string, ev []*events.Event, ctx ...interface{}) error {
	projection, err := FindByDeliveryId(deliveryId, ctx...)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if projection == nil {
		projection = &DeliveryProjection{
			DeliveryId: deliveryId,
		}
	}

	for _, e := range ev {
		projection = projection.update(e)
	}

	if _, err := insert(projection, ctx...); err != nil {
		return err
	}

	return nil
}

func (dp *DeliveryProjection) update(event *events.Event) *DeliveryProjection {
	switch event.Type {
	case events.ConfirmDelivery:
		dp.Status = "confirmed"
	case events.CancelledDelivery:
		dp.Status = "cancelled"
	case events.SetOnTheGoDelivery:
		dp.Status = "on_the_go"
	case events.SetDeliveredDelivery:
		dp.Status = "delivered"
	}
	dp.LastModified = event.Created
	return dp
}
