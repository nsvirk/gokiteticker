package kiteticker

import "github.com/nsvirk/gokiteticker/models"

// Order represents a individual order response.
type Order struct {
	AccountID string `json:"account_id"`
	PlacedBy  string `json:"placed_by"`

	OrderID                 string                 `json:"order_id"`
	ExchangeOrderID         string                 `json:"exchange_order_id"`
	ParentOrderID           string                 `json:"parent_order_id"`
	Status                  string                 `json:"status"`
	StatusMessage           string                 `json:"status_message"`
	StatusMessageRaw        string                 `json:"status_message_raw"`
	OrderTimestamp          models.Time            `json:"order_timestamp"`
	ExchangeUpdateTimestamp models.Time            `json:"exchange_update_timestamp"`
	ExchangeTimestamp       models.Time            `json:"exchange_timestamp"`
	Variety                 string                 `json:"variety"`
	Modified                bool                   `json:"modified"`
	Meta                    map[string]interface{} `json:"meta"`

	Exchange        string `json:"exchange"`
	TradingSymbol   string `json:"tradingsymbol"`
	InstrumentToken uint32 `json:"instrument_token"`

	OrderType         string  `json:"order_type"`
	TransactionType   string  `json:"transaction_type"`
	Validity          string  `json:"validity"`
	ValidityTTL       int     `json:"validity_ttl"`
	Product           string  `json:"product"`
	Quantity          float64 `json:"quantity"`
	DisclosedQuantity float64 `json:"disclosed_quantity"`
	Price             float64 `json:"price"`
	TriggerPrice      float64 `json:"trigger_price"`

	AveragePrice      float64 `json:"average_price"`
	FilledQuantity    float64 `json:"filled_quantity"`
	PendingQuantity   float64 `json:"pending_quantity"`
	CancelledQuantity float64 `json:"cancelled_quantity"`

	AuctionNumber string `json:"auction_number"`

	Tag  string   `json:"tag"`
	Tags []string `json:"tags"`
}
