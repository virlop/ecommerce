
package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"deliverygo/projections"
	"deliverygo/rest/server"
	"deliverygo/security"
)

// Define las rutas del servicio REST
func init() {
	server.Router().PUT("/v1/delivery/:deliveryId", server.ValidateAuthentication, updateDeliveryStatus)
	server.Router().GET("/v1/delivery", server.ValidateAuthentication, listDeliveriesByStatus)
	server.Router().GET("/v1/delivery/:orderId", server.ValidateAuthentication, getDeliveryByOrderId)
}

// Estructura para el cuerpo de la solicitud de actualización de estado
type UpdateDeliveryRequest struct {
	Status string `json:"status" binding:"required"`
}

// Estructura para la respuesta de un delivery
type DeliveryResponse struct {
	DeliveryId string `json:"deliveryId"`
	OrderId    string `json:"orderId"`
	Status     string `json:"status"`
	Created    time.Time `json:"created"`
	UserId     string `json:"userId"`
}

// Actualizar estado de un delivery
func updateDeliveryStatus(c *gin.Context) {
	deliveryId := c.Param("deliveryId")
	var req UpdateDeliveryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		server.AbortWithError(c, err)
		return
	}

	ctx := server.GinCtx(c)
	delivery, err := projections.FindDeliveryById(deliveryId, ctx...)
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	if err := projections.UpdateDeliveryStatus(deliveryId, req.Status, ctx...); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, DeliveryResponse{
		DeliveryId: deliveryId,
		OrderId:    delivery.OrderId,
		Status:     req.Status,
		Created:    time.Now(),
		UserId:     delivery.UserId,
	})
}

// Listar deliveries por estado
func listDeliveriesByStatus(c *gin.Context) {
	status := c.Query("status")

	if !security.IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "El usuario no tiene permisos de admin"})
		return
	}

	ctx := server.GinCtx(c)
	deliveries, err := projections.FindDeliveriesByStatus(status, ctx...)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deliveries NO encontrados"})
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

// Obtener detalles de un delivery por orderId
func getDeliveryByOrderId(c *gin.Context) {
	orderId := c.Param("orderId")
	ctx := server.GinCtx(c)
	delivery, err := projections.FindDeliveryByOrderId(orderId, ctx...)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery NO encontrado"})
		return
	}

	user, err := security.GetAuthenticatedUser(c)
	if err != nil || user.ID != delivery.UserId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Delivery NO corresponde al usuario autenticado"})
		return
	}

	c.JSON(http.StatusOK, DeliveryResponse{
		DeliveryId: delivery.DeliveryId,
		OrderId:    orderId,
		Status:     delivery.Status,
		Created:    delivery.Created,
		UserId:     delivery.UserId,
	})
}
