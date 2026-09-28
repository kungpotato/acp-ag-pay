package acp

import "net/http"

type OrderHandler struct {
	orders *OrderStore
}

func NewOrderHandler(orders *OrderStore) *OrderHandler {
	return &OrderHandler{orders: orders}
}

func (h *OrderHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/orders/{id}", h.get)
}

func (h *OrderHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, ok := h.orders.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "order not found"})
		return
	}
	writeJSON(w, http.StatusOK, order)
}
