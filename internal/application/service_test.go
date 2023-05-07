package application

import (
	"L0/internal/domain"
	"L0/internal/infrastructure/cache"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

type StorageStub struct {
	data []domain.Order
}

func NewStorageStub() *StorageStub {
	return &StorageStub{
		data: []domain.Order{
			{
				OrderUid: "aaaa",
			},
		},
	}
}

func (s *StorageStub) ReceiveAll() ([]domain.Order, error) {
	return s.data, nil
}

func (s *StorageStub) ReceiveOrder(orderUID string) (domain.Order, error) {
	for _, order := range s.data {
		if order.OrderUid == orderUID {
			return order, nil
		}
	}
	return domain.Order{}, fmt.Errorf("cant find order")
}

func (s *StorageStub) Insert(in domain.Order) error {
	s.data = append(s.data, in)
	return nil
}

func TestRun(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		testStorage := NewStorageStub()
		testCache, err := cache.NewOrderCache(testStorage)
		require.NoError(t, err)

		service := NewService(testCache, testStorage, logrus.New())
		expectedOrder := domain.Order{OrderUid: "qweqw"}

		service.insertOrder(expectedOrder)
		receivedOrder, err := testStorage.ReceiveOrder(expectedOrder.OrderUid)
		require.NoError(t, err)

		require.Equal(t, expectedOrder, receivedOrder)
	})

}
