package cache

import (
	"L0/internal/domain"
	"fmt"
	"sort"
	"sync"
)

type storage interface {
	ReceiveAll() ([]domain.Order, error)
	ReceiveOrder(orderUID string) (domain.Order, error)
}

type OrderCache struct {
	mtx     sync.RWMutex
	data    map[string]domain.Order
	storage storage
}

func NewOrderCache(storage storage) (*OrderCache, error) {
	orders := make(map[string]domain.Order)

	data, err := storage.ReceiveAll()
	if err != nil {
		fmt.Printf("can't read from storage: %s", err)
		return nil, err
	}

	for _, order := range data {
		orders[order.OrderUid] = order
	}

	return &OrderCache{
		data:    orders,
		storage: storage,
	}, nil
}

func (c *OrderCache) Set(orderUID string, order domain.Order) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[orderUID] = order
}

func (c *OrderCache) Get(orderUID string) (domain.Order, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	order, ok := c.data[orderUID]
	if !ok {
		order, err := c.storage.ReceiveOrder(orderUID)
		if err != nil {
			return domain.Order{}, err
		}
		return order, nil
	}
	return order, nil
}

func (c *OrderCache) GetKeyList() []string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
