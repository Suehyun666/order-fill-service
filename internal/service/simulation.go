package service

import (
	"fmt"
	"log/slog"
	"math/rand"
	"order-fill-service/internal/domain"
	"order-fill-service/internal/infra/kafka"
	"time"
)

// SimulationService : Filler 인터페이스 구현체
type SimulationService struct {
	producer *kafka.EventProducer // 의존성 주입 (Producer를 알고 있음)
}

func NewSimulationService(p *kafka.EventProducer) *SimulationService {
	return &SimulationService{producer: p}
}

// ProcessOrder : 주문을 받아서 랜덤 딜레이 후 체결 이벤트 발송
func (s *SimulationService) ProcessOrder(req domain.OrderRequest) error {
	// Goroutine으로 실행하여 Non-blocking 처리 (Kafka Consumer가 막히지 않게)
	go func(order domain.OrderRequest) {
		// 1. 랜덤 딜레이 (100ms ~ 500ms) - 거래소 네트워크 지연 시뮬레이션
		delay := time.Duration(rand.Intn(400)+100) * time.Millisecond
		time.Sleep(delay)

		// 2. 체결 이벤트 생성 (100% 체결 가정)
		fillEvent := domain.FillEvent{
			EventID:       generateEventID(),
			ClientOrderID: order.OrderID,
			AccountID:     order.AccountID, // 필수! account-command가 사용
			SecurityID:    1,                // TODO: Symbol을 SecurityID로 매핑
			Side:          order.Side,
			Fills: []domain.Fill{
				{
					PriceMicroUnits: order.Price,
					Quantity:        order.Quantity,
				},
			},
			Status:       "FILLED",
			TransactTime: time.Now(),
		}

		slog.Info("Fill simulated",
			"order_id", order.OrderID,
			"account_id", order.AccountID,
			"delay_ms", delay.Milliseconds())

		// 3. Producer를 통해 전송
		if err := s.producer.SendFillEvent(fillEvent); err != nil {
			slog.Error("Failed to send fill event", "error", err)
		}
	}(req)

	return nil
}

// generateEventID : 유니크한 이벤트 ID 생성
func generateEventID() string {
	return fmt.Sprintf("fill-%d-%d", time.Now().UnixNano(), rand.Intn(10000))
}
