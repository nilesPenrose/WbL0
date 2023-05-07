package application

import (
	"L0/internal/domain"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type storage interface {
	ReceiveAll() ([]domain.Order, error)
	ReceiveOrder(orderUID string) (domain.Order, error)
	Insert(in domain.Order) error
}

type orderCache interface {
	Get(orderUID string) (domain.Order, error)
	Set(orderUID string, order domain.Order)
}

type Service struct {
	storage storage
	cache   orderCache
	logger  *logrus.Logger
}

func NewService(cache orderCache, storage storage, logger *logrus.Logger) *Service {
	return &Service{
		storage: storage,
		cache:   cache,
		logger:  logger,
	}
}

func (s *Service) Run(m *stan.Msg) {
	var order domain.Order

	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		s.logger.Printf("cant process msg:%s, got: %s\n", m.Data, err.Error())
		_ = m.Ack()
		return
	}

	s.insertOrder(order)

	_ = m.Ack()
}

func (s *Service) insertOrder(order domain.Order) {
	if !order.IsValid() {
		s.logger.Printf("order is invalid; recuire: order_uid")
		return
	}

	if err := s.storage.Insert(order); err != nil {
		s.logger.Printf("cant insert order:%v, got: %s\n", order, err.Error())
		return
	}

	s.cache.Set(order.OrderUid, order)

}
