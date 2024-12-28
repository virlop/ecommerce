package db

import (
	"context"
	"fmt"
	"log" //maneja logs para registar errores o info general del sistema
	"os"  //da acceso a la variable de entorno, que uso para ibtener la URI de la conexión de mongo

	//del driver de mongo
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

// una referencia global a la db que estoy usando, delivery
var database *mongo.Database

// Get obtiene la base de datos MongoDB para Delivery
func Get(deps ...interface{}) (*mongo.Database, error) {
	if database == nil {
		// Lee la URL de conexión desde una variable de entorno
		mongoURL := os.Getenv("MONGO_URI")
		if mongoURL == "" {
			return nil, fmt.Errorf("MONGO_URI no está configurada")
		}

		// Configura las opciones del cliente
		clientOptions := options.Client().ApplyURI(mongoURL)

		// Conecta a MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatalf("Error al conectar con MongoDB: %v", err)
			return nil, err
		}

		// Selecciona la base de datos "delivery"
		database = client.Database("delivery")
	}
	return database, nil
}

// CheckError reinicia la base de datos en caso de error crítico
func CheckError(err interface{}) {
	// Si ocurre un error de tiempo de espera (timeout),
	// es decir, la BD no esta disponible,
	// resetea la variable global database a nil
	if err == topology.ErrServerSelectionTimeout {
		database = nil
	}
}

// IsUniqueKeyError retorna true si el error es por un índice único
func IsUniqueKeyError(err error) bool {
	if wErr, ok := err.(mongo.WriteException); ok {
		for i := 0; i < len(wErr.WriteErrors); i++ {
			if wErr.WriteErrors[i].Code == 11000 {
				return true
			}
		}
	}
	return false
}
