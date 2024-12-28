// Configura el logger, carga las variables de entorno, inicializa la base de datos y RabbitMQ
// Inicia los diferentes servidores: REST y rabbit
package main

import (
	routes "deliverygo/rest"
	"deliverygo/rabbit/consume"
	"deliverygo/tools/db"
)

func main() {
	consume.Init()
	routes.Start()
}
