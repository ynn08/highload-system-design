package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/highload-system-design/order-service/internal/domain"
	"github.com/user/highload-system-design/order-service/internal/usecase"
)

type OrderHandler struct {
	createOrderUC *usecase.CreateOrderUseCase
	manageCartUC  *usecase.ManageCartUseCase
}

func NewOrderHandler(couc *usecase.CreateOrderUseCase, mcuc *usecase.ManageCartUseCase) *OrderHandler {
	return &OrderHandler{createOrderUC: couc, manageCartUC: mcuc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input domain.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.createOrderUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) AddToCart(c *gin.Context) {
	customerID := c.Param("customerId")
	restaurantID := c.Query("restaurantId")
	var item domain.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.manageCartUC.AddToCart(c.Request.Context(), customerID, restaurantID, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *OrderHandler) GetCart(c *gin.Context) {
	customerID := c.Param("customerId")
	cart, err := h.manageCartUC.GetCart(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}
	c.JSON(http.StatusOK, cart)
}
