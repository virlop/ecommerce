package events

import (
	"context"
	"deliverygo/tools/db"
	"deliverygo/tools/log"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

// Configura y devuelve la colección deliveryEvents de MongoDB.
func dbCollection(deps ...interface{}) (*mongo.Collection, error) {
	if collection != nil {
		return collection, nil
	}

	database, err := db.Get(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	col := database.Collection("deliveryEvents")

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"deliveryId": 1, // Índice en deliveryId
			},
		},
	)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	collection = col
	return collection, nil
}

// Insertar un evento de Delivery en la collection de mongo
func InsertDeliveryEvent(event *Event, deps ...interface{}) (*Event, error) {
	if err := event.ValidateSchema(); err != nil { //valida el esquema del evento
		log.Get(deps...).Error(err)
		return nil, err
	}

	var collection, err = dbCollection(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}
	// Validar el status
	if event.DeliveryStatus != DeliveryStatusConfirmed && event.DeliveryStatus != DeliveryStatusCancelled &&
		event.DeliveryStatus != DeliveryStatusOnTheGo && event.DeliveryStatus != DeliveryStatusDelivered {
		err := fmt.Errorf("invalid status: %s", event.DeliveryStatus)
		log.Get(deps...).Error(err)
		return nil, err
	}

	if _, err := collection.InsertOne(context.Background(), event); err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	return event, nil
}

// Filtra los eventos de delivery por status
func FindDeliveryEventsByStatus(deliveryStatus string, deps ...interface{}) ([]*Event, error) {
	var collection, err = dbCollection(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	// Filtrar por estado
	filter := bson.M{
		"deliveryStatus": deliveryStatus,
	}
	cur, err := collection.Find(context.Background(), filter, nil)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	events := []*Event{}
	for cur.Next(context.Background()) {
		event := &Event{}
		if err := cur.Decode(event); err != nil {
			log.Get(deps...).Error(err)
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// Buscar eventos relacionados a un deliveryId
func FindDeliveryEventsByDeliveryId(deliveryId string, deps ...interface{}) ([]*Event, error) {
	var collection, err = dbCollection(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	filter := bson.M{"deliveryId": deliveryId}
	cur, err := collection.Find(context.Background(), filter, nil)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	events := []*Event{}
	for cur.Next(context.Background()) {
		event := &Event{}
		if err := cur.Decode(event); err != nil {
			log.Get(deps...).Error(err)
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
func FindDeliveryIdByOrderId(orderId string, ctx ...interface{}) (string, error) {
	var collection, err = dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return "", err
	}

	// Filtrar por orderId
	filter := bson.M{"orderId": orderId}
	cur, err := collection.Find(context.Background(), filter, nil)
	if err != nil {
		log.Get(ctx...).Error(err)
		return "", err
	}
	defer cur.Close(context.Background())

	// Buscar el deliveryId asociado
	var deliveryId string
	for cur.Next(context.Background()) {
		event := &Event{}
		if err := cur.Decode(event); err != nil {
			log.Get(ctx...).Error(err)
			return "", err
		}
		// Asignar el deliveryId si existe
		deliveryId = event.DeliveryId
		break
		// Solo necesitamos un `deliveryId`, salimos del bucle
	}

	// Verificar si se encontró un deliveryId
	if deliveryId == "" {
		return "", fmt.Errorf("no deliveryId found for orderId: %s", orderId)
	}

	return deliveryId, nil
}
