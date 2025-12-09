package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"order-fill-service/internal/config"
	"order-fill-service/internal/infra/kafka"
	"order-fill-service/internal/logger"
	"order-fill-service/internal/service"
)

func main() {
	cfg := config.Load()

	logger.Setup(cfg.LogLevel)

	slog.Info("Starting Order Fill Service",
		"brokers", cfg.KafkaBrokers,
		"log_level", cfg.LogLevel,
	)

	reqTopic := "order.created"
	resTopic := "order.filled"

	producer, err := kafka.NewEventProducer(cfg.KafkaBrokers, resTopic)
	if err != nil {
		slog.Error("Failed to create producer", "error", err)
		os.Exit(1)
	}
	defer producer.Close()

	svc := service.NewSimulationService(producer)

	consumer, err := kafka.NewOrderConsumer(cfg.KafkaBrokers, reqTopic, svc)
	if err != nil {
		slog.Error("Failed to create consumer", "error", err)
		os.Exit(1)
	}
	defer consumer.Close()

	go consumer.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down service...")
}
