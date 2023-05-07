package handlers

import (
	"L0/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type orderCache interface {
	Get(orderUID string) (domain.Order, error)
	Set(orderUID string, order domain.Order)
	GetKeyList() []string
}

type OrderHandler struct {
	cache orderCache
}

func NewOrderHandler(cache orderCache) *OrderHandler {
	return &OrderHandler{
		cache: cache,
	}
}

func (o *OrderHandler) ReceiveOrderByID(w http.ResponseWriter, req *http.Request) {
	data := req.URL.Query().Get("orderUID")

	order, err := o.cache.Get(data)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(order)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	_, _ = fmt.Fprintf(w, "%s", resp)
}

func (o *OrderHandler) ReceiveAllOrderUIDs(w http.ResponseWriter, req *http.Request) {
	keys := o.cache.GetKeyList()
	fmt.Println(keys)
	jsonKeys, _ := json.Marshal(keys)

	_, _ = fmt.Fprintf(w, "%s", jsonKeys)
}
