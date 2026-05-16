package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/user/highload-system-design/delivery-service/internal/domain"
	"github.com/user/highload-system-design/delivery-service/internal/usecase"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DeliveryHandler struct {
	assignUC      *usecase.AssignCourierUseCase
	saveCourierUC *usecase.SaveCourierUseCase
	deliveryRepo  usecase.DeliveryRepository
}

func NewDeliveryHandler(
	assignUC *usecase.AssignCourierUseCase,
	saveCourierUC *usecase.SaveCourierUseCase,
	deliveryRepo usecase.DeliveryRepository,
) *DeliveryHandler {
	return &DeliveryHandler{
		assignUC:      assignUC,
		saveCourierUC: saveCourierUC,
		deliveryRepo:  deliveryRepo,
	}
}

func (h *DeliveryHandler) GetStatus(c *gin.Context) {
	orderID := c.Param("orderId")
	delivery, err := h.deliveryRepo.FindByID(c.Request.Context(), orderID)
	if err != nil || delivery == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "delivery not found"})
		return
	}
	c.JSON(http.StatusOK, delivery)
}

func (h *DeliveryHandler) SaveCourier(c *gin.Context) {
	var courier domain.Courier
	if err := c.ShouldBindJSON(&courier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.saveCourierUC.Execute(c.Request.Context(), &courier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *DeliveryHandler) WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
