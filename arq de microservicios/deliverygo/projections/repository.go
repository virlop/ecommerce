// Maneja la interacción con la base de datos MongoDB.
// Implementa funciones para insertar, actualizar o buscar documentos en la colección de proyecciones (delivery_projection).

package projections

import (
	"context"
	"deliverygo/tools/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func dbCollection(ctx ...interface{}) (*mongo.Collection, error) {
	if collection != nil {
		return collection, nil
	}

	// db.Get() para obtener la conexión a MongoDB
	database, err := db.Get(ctx...)
	if err != nil {
		return nil, err
	}

	collection = database.Collection("delivery_projection")
	return collection, nil
}

func insert(dp *DeliveryProjection, ctx ...interface{}) (*DeliveryProjection, error) {
	if err := dp.Validate(); err != nil {
		return nil, err
	}

	col, err := dbCollection(ctx...)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"deliveryId": dp.DeliveryId}
	update := bson.M{"$set": dp}
	options := options.Update().SetUpsert(true)

	if _, err := col.UpdateOne(context.Background(), filter, update, options); err != nil {
		return nil, err
	}

	return dp, nil
}

func FindByDeliveryId(deliveryId string, ctx ...interface{}) (*DeliveryProjection, error) {
	col, err := dbCollection(ctx...)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"deliveryId": deliveryId}
	dp := &DeliveryProjection{}
	if err := col.FindOne(context.Background(), filter).Decode(dp); err != nil {
		return nil, err
	}

	return dp, nil
}
