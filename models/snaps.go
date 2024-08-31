package models

// OHLC represents OHLC packets.
type OHLC struct {
	InstrumentToken uint32  `json:"-"`
	Open            float64 `json:"open"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Close           float64 `json:"close"`
}

// DepthItem represents a single market depth entry.
type DepthItem struct {
	Price    float64 `json:"price"`
	Quantity uint32  `json:"quantity"`
	Orders   uint32  `json:"orders"`
}

// Depth represents a group of buy/sell market depths.
type Depth struct {
	Buy  [5]DepthItem `json:"buy"`
	Sell [5]DepthItem `json:"sell"`
}

// Tick represents a single packet in the market feed.
type Tick struct {
	Mode               string  `json:"mode"`
	InstrumentToken    uint32  `json:"instrument_token"`
	IsTradable         bool    `json:"is_tradable"`
	IsIndex            bool    `json:"is_index"`
	Timestamp          Time    `json:"timestamp"` // Timestamp represents Exchange timestamp
	LastTradeTime      Time    `json:"last_trade_time"`
	LastPrice          float64 `json:"last_price"`
	LastTradedQuantity uint32  `json:"last_traded_quantity"`
	TotalBuyQuantity   uint32  `json:"total_buy_quantity"`
	TotalSellQuantity  uint32  `json:"total_sell_quantity"`
	VolumeTraded       uint32  `json:"volume_traded"`
	TotalBuy           uint32  `json:"total_buy"`
	TotalSell          uint32  `json:"total_sell"`
	AverageTradePrice  float64 `json:"average_trade_price"`
	OI                 uint32  `json:"oi"`
	OIDayHigh          uint32  `json:"oi_day_high"`
	OIDayLow           uint32  `json:"oi_day_low"`
	NetChange          float64 `json:"net_change"`
	OHLC               OHLC    `json:"ohlc"`
	Depth              Depth   `json:"depth"`
}
