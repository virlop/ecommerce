# Etapa de construcción
FROM golang:1.19 AS builder

# Define el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos necesarios
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Construye el binario del servicio
RUN go build -o deliverygo .

# Etapa de ejecución
FROM alpine:latest

# Define el directorio de trabajo en el contenedor final
WORKDIR /root/

# Copia el binario construido en la etapa anterior
COPY --from=builder /app/deliverygo .

# Expone el puerto en el que el microservicio escucha
EXPOSE 8080

# Comando por defecto para ejecutar el servicio
CMD ["./deliverygo"]
