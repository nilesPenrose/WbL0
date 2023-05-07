package main

import (
	"L0/internal/application"
	"L0/internal/handlers"
	"L0/internal/infrastructure/cache"
	postgres "L0/internal/infrastructure/storage"
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func main() {
	pgsStorage, err := postgres.NewPgsStorage(context.Background())
	if err != nil {
		log.Fatal("new pgsStorage ", err)
	}

	logger := logrus.New()

	orderCache, err := cache.NewOrderCache(pgsStorage)
	if err != nil {
		log.Fatal("new orderCache ", err)
	}
	fmt.Println(orderCache)

	publisherSC, _ := stan.Connect("test-cluster", "2")
	publisher := application.NewPublisher(publisherSC)

	handler := handlers.NewPublishHandler(publisher)

	orderHandler := handlers.NewOrderHandler(orderCache)

	service := application.NewService(orderCache, pgsStorage, logger)

	sc, err := stan.Connect("test-cluster", "1")
	if err != nil {
		log.Fatal("stan connect ", err)
	}

	sub, err := sc.Subscribe("orders", service.Run)
	if err != nil {
		log.Fatal("stan subs ", err)
	}
	defer sub.Close()

	http.HandleFunc("/", handler.Publish)

	http.HandleFunc("/order/get", orderHandler.ReceiveOrderByID)
	http.HandleFunc("/order/list", orderHandler.ReceiveAllOrders)
	http.HandleFunc("/order/uid/list", orderHandler.ReceiveAllOrderUIDs)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("can't up the http server ", err)
	}
}
