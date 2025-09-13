package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kimvlry/simple-order-service/internal/interfaces"
)

type OrderHandler struct {
	orderService interfaces.OrderService
}

func NewOrderHandler(orderService interfaces.OrderService) *OrderHandler {
	return &OrderHandler{orderService}
}

func (h *OrderHandler) GetOrder(writer http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(writer, "missing id parameter", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrderByID(id, req.Context())
	if err != nil {
		http.Error(writer, fmt.Sprintf("internal error: %s", err), http.StatusInternalServerError)
		return
	}
	if order == nil {
		http.NotFound(writer, req)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(order)
}

func (h *OrderHandler) Ping(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(writer, "hello from simple order-service")
}
