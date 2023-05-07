package application

import (
	"L0/internal/domain"
	"encoding/json"
	"github.com/nats-io/stan.go"
)

type Publisher struct {
	sc stan.Conn
}

func NewPublisher(sc stan.Conn) *Publisher {
	return &Publisher{
		sc: sc,
	}
}
func (p *Publisher) WriteToCh(order domain.Order) {
	data, _ := json.Marshal(order)
	_ = p.sc.Publish("orders", data)
}
