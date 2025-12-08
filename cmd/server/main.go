package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"order-fill-service/internal/infra/kafka"
	"order-fill-service/internal/service"
)

func main() {
	// 1. 설정 (하드코딩 대신 환경변수로 빼는 게 좋음)
	brokers := []string{"localhost:9092"}
	reqTopic := "order.created"
	resTopic := "order.filled"

	slog.Info("Starting Order Fill Service...")

	// 2. Producer 초기화 (Infrastructure)
	producer, err := kafka.NewEventProducer(brokers, resTopic)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	// 3. Service 초기화 (Business Logic) - Producer 주입
	svc := service.NewSimulationService(producer)

	// 4. Consumer 초기화 (Infrastructure) - Service 주입
	consumer, err := kafka.NewOrderConsumer(brokers, reqTopic, svc)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// 5. 컨슈머 실행 (별도 고루틴이 아니므로 메인 스레드 블로킹 방지 위해 고루틴으로)
	go consumer.Start()

	// 6. Graceful Shutdown 대기
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down service...")
}
