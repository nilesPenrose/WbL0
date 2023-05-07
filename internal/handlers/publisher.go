package handlers

import (
	"L0/internal/application"
	"L0/internal/domain"
	"encoding/json"
	"io"
	"net/http"
)

type PublishHandler struct {
	publisher *application.Publisher
}

func NewPublishHandler(publisher *application.Publisher) *PublishHandler {
	return &PublishHandler{
		publisher: publisher,
	}
}

func (p *PublishHandler) Publish(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var order domain.Order

	if err = json.Unmarshal(data, &order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	p.publisher.WriteToCh(order)
}
