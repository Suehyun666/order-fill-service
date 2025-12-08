package kafka

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"order-fill-service/internal/domain"

	"github.com/IBM/sarama"
)

// EventProducer : 메시지 발송 역할 정의
type EventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewEventProducer : 생성자 (Constructor)
func NewEventProducer(brokers []string, topic string) (*EventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &EventProducer{producer: p, topic: topic}, nil
}

// SendFillEvent : 도메인 객체를 받아서 카프카로 쏨
func (ep *EventProducer) SendFillEvent(event domain.FillEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling failed: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: ep.topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", event.AccountID)), // 파티셔닝 키 (같은 계좌는 같은 파티션)
		Value: sarama.ByteEncoder(bytes),
	}

	partition, offset, err := ep.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka send failed: %w", err)
	}

	slog.Info("Message Sent",
		"partition", partition,
		"offset", offset,
		"order_id", event.ClientOrderID,
		"account_id", event.AccountID)
	return nil
}

func (ep *EventProducer) Close() error {
	return ep.producer.Close()
}
