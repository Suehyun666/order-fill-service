package kafka

import (
	"encoding/json"
	"log/slog"
	"order-fill-service/internal/domain"

	"github.com/IBM/sarama"
)

type OrderConsumer struct {
	consumer sarama.Consumer
	service  domain.Filler // 인터페이스 의존! (구체적인 Service 몰라도 됨)
	topic    string
}

func NewOrderConsumer(brokers []string, topic string, service domain.Filler) (*OrderConsumer, error) {
	c, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &OrderConsumer{consumer: c, service: service, topic: topic}, nil
}

func (oc *OrderConsumer) Start() {
	partitionConsumer, err := oc.consumer.ConsumePartition(oc.topic, 0, sarama.OffsetNewest)
	if err != nil {
		slog.Error("Failed to start partition consumer", "error", err)
		return
	}
	defer partitionConsumer.Close()

	slog.Info("Consumer Started", "topic", oc.topic)

	for msg := range partitionConsumer.Messages() {
		var req domain.OrderRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			slog.Error("Invalid JSON", "error", err)
			continue
		}

		slog.Info("Order Received", "order_id", req.OrderID)

		// 서비스 로직 호출 (위임)
		oc.service.ProcessOrder(req)
	}
}

func (oc *OrderConsumer) Close() error {
	return oc.consumer.Close()
}
