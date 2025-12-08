package domain

import "time"

// OrderRequest : 주문 요청 (들어오는 데이터) - order.created 토픽에서 받음
type OrderRequest struct {
	OrderID   string  `json:"order_id"`
	AccountID int64   `json:"account_id"` // account-command가 체결 처리할 때 필요
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`       // BUY or SELL
	Price     int64   `json:"price"`      // micro units (1,000,000 = 1.0)
	Quantity  int64   `json:"quantity"`
}

// FillEvent : 체결 결과 (나가는 데이터) - order.filled 토픽으로 발행
type FillEvent struct {
	EventID       string    `json:"event_id"`
	ClientOrderID string    `json:"client_order_id"`
	AccountID     int64     `json:"account_id"`     // 필수! account-command, streaming-service가 사용
	SecurityID    int32     `json:"security_id"`
	Side          string    `json:"side"`           // BUY or SELL
	Fills         []Fill    `json:"fills"`          // 체결 내역
	Status        string    `json:"status"`         // FILLED
	TransactTime  time.Time `json:"transact_time"`
}

// Fill : 체결 정보
type Fill struct {
	PriceMicroUnits int64 `json:"price_micro_units"`
	Quantity        int64 `json:"quantity"`
}

// Filler : 비즈니스 로직 인터페이스 (행동 정의)
// 나중에 Mock 테스트하기 쉽도록 인터페이스로 정의합니다.
type Filler interface {
	ProcessOrder(req OrderRequest) error
}
