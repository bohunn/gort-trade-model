// Package models provides shared data structures for the liquidation heatmap service
package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Exchange represents supported exchanges
type Exchange string

const (
	ExchangeBinance  Exchange = "binance"
	ExchangeOKX      Exchange = "okx"
	ExchangeBybit    Exchange = "bybit"
	ExchangeCoinbase Exchange = "coinbase"
	ExchangeKraken   Exchange = "kraken"
	ExchangeDeribit  Exchange = "deribit"
	ExchangeBitfinex Exchange = "bitfinex"
)

// Symbol represents a trading pair
type Symbol string

const (
	SymbolBTCUSDT Symbol = "BTCUSDT"
	SymbolETHUSDT Symbol = "ETHUSDT"
	SymbolBNBUSDT Symbol = "BNBUSDT"
	SymbolSOLUSDT Symbol = "SOLUSDT"
	SymbolXRPUSDT Symbol = "XRPUSDT"
)

// OrderType represents the type of liquidation order
type OrderType string

const (
	OrderTypeLiquidation OrderType = "liquidation"
	OrderTypeADL         OrderType = "adl" // Auto-deleveraging
	OrderTypeBankruptcy  OrderType = "bankruptcy"
)

// Side represents position side
type Side string

const (
	SideLong  Side = "long"
	SideShort Side = "short"
	SideBuy   Side = "BUY"  // Binance format
	SideSell  Side = "SELL" // Binance format
)

// Interval represents time intervals for aggregation
type Interval string

const (
	Interval1s  Interval = "1s"
	Interval1m  Interval = "1m"
	Interval5m  Interval = "5m"
	Interval15m Interval = "15m"
	Interval1h  Interval = "1h"
	Interval4h  Interval = "4h"
	Interval1d  Interval = "1d"
)

// ===========================================
// RAW MARKET DATA STRUCTURES
// ===========================================

// MarketSnapshot represents current market state
type MarketSnapshot struct {
	Exchange        Exchange `json:"exchange"`
	Symbol          Symbol   `json:"symbol"`
	Timestamp       int64    `json:"timestamp"`
	MarkPrice       float64  `json:"mark_price"`
	IndexPrice      float64  `json:"index_price"`
	FundingRate     float64  `json:"funding_rate"`
	OpenInterest    float64  `json:"open_interest"`     // in contracts
	OpenInterestUSD float64  `json:"open_interest_usd"` // in USD
	Volume24h       float64  `json:"volume_24h"`        // in USD
	Turnover24h     float64  `json:"turnover_24h"`      // in USD
	NextFundingTime int64    `json:"next_funding_time"`
}

// LiquidationEvent represents a single liquidation from exchange
type LiquidationEvent struct {
	Exchange       Exchange  `json:"exchange"`
	Symbol         Symbol    `json:"symbol"`
	Timestamp      int64     `json:"timestamp"`
	Side           Side      `json:"side"`     // BUY/SELL or long/short
	Price          float64   `json:"price"`    // Liquidation price
	Quantity       float64   `json:"quantity"` // Contract quantity
	Value          float64   `json:"value"`    // USD value
	OrderType      OrderType `json:"order_type"`
	AvgPrice       float64   `json:"avg_price,omitempty"`        // Average fill price
	FilledQty      float64   `json:"filled_qty,omitempty"`       // Filled quantity
	OrderStatus    string    `json:"order_status,omitempty"`     // Order status
	OrderTradeTime int64     `json:"order_trade_time,omitempty"` // Trade execution time
}

// OrderBookSnapshot represents order book state
type OrderBookSnapshot struct {
	Exchange     Exchange     `json:"exchange"`
	Symbol       Symbol       `json:"symbol"`
	Timestamp    int64        `json:"timestamp"`
	Bids         []PriceLevel `json:"bids"`
	Asks         []PriceLevel `json:"asks"`
	LastUpdateID int64        `json:"last_update_id,omitempty"`
	Spread       float64      `json:"spread,omitempty"`
	MidPrice     float64      `json:"mid_price,omitempty"`
	Imbalance    float64      `json:"imbalance,omitempty"` // -1 to 1
}

// PriceLevel represents a price and size at that level
type PriceLevel struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Count    int     `json:"count,omitempty"` // Number of orders at this level
}

// ===========================================
// CALCULATED HEATMAP STRUCTURES
// ===========================================

// HeatmapData represents the complete liquidation heatmap
type HeatmapData struct {
	Symbol       Symbol               `json:"symbol"`
	Exchange     Exchange             `json:"exchange,omitempty"`
	Timestamp    int64                `json:"timestamp"`
	Interval     Interval             `json:"interval"`
	CurrentPrice float64              `json:"current_price"`
	Levels       []LiquidationLevel   `json:"levels"`
	Clusters     []LiquidationCluster `json:"clusters"`
	Summary      HeatmapSummary       `json:"summary"`
}

// LiquidationLevel represents liquidations at a specific price
type LiquidationLevel struct {
	Price             float64 `json:"price"`
	LongLiquidations  float64 `json:"long_liquidations"`  // USD volume
	ShortLiquidations float64 `json:"short_liquidations"` // USD volume
	TotalVolume       float64 `json:"total_volume"`       // Total USD volume
	Intensity         float64 `json:"intensity"`          // 0-100 score
	Timestamp         int64   `json:"timestamp"`
}

// LiquidationCluster represents a cluster of significant liquidation levels
type LiquidationCluster struct {
	Symbol          Symbol             `json:"symbol"`
	PriceRangeStart float64            `json:"price_range_start"`
	PriceRangeEnd   float64            `json:"price_range_end"`
	Levels          []LiquidationLevel `json:"levels"`
	TotalVolume     float64            `json:"total_volume"`
	PeakIntensity   float64            `json:"peak_intensity"`
	UpdatedAt       int64              `json:"updated_at"`
}

// HeatmapSummary contains aggregated heatmap statistics
type HeatmapSummary struct {
	TotalLongLiquidations  float64        `json:"total_long_liquidations"`
	TotalShortLiquidations float64        `json:"total_short_liquidations"`
	MaxLiquidationPrice    float64        `json:"max_liquidation_price"`
	MaxLiquidationVolume   float64        `json:"max_liquidation_volume"`
	WeightedAvgLongPrice   float64        `json:"weighted_avg_long_price"`
	WeightedAvgShortPrice  float64        `json:"weighted_avg_short_price"`
	SignificantLevels      int            `json:"significant_levels"`
	CriticalZones          []CriticalZone `json:"critical_zones"`
}

// CriticalZone represents a high-risk liquidation zone
type CriticalZone struct {
	PriceStart float64 `json:"price_start"`
	PriceEnd   float64 `json:"price_end"`
	Type       string  `json:"type"` // "long", "short", or "mixed"
	Intensity  float64 `json:"intensity"`
	Volume     float64 `json:"volume"`
}

// ===========================================
// STREAM MESSAGE STRUCTURES
// ===========================================

// StreamMessage represents a message for Redis Streams
type StreamMessage struct {
	ID        string                 `json:"id"`     // Stream message ID
	Stream    string                 `json:"stream"` // Stream name
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ToStreamMessage converts any model to a StreamMessage
func ToStreamMessage(streamName string, v interface{}) (*StreamMessage, error) {
	data, err := structToMap(v)
	if err != nil {
		return nil, err
	}

	return &StreamMessage{
		Stream:    streamName,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}, nil
}

// structToMap converts a struct to a map for Redis
func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	// Flatten the map for Redis (convert nested objects to JSON strings)
	result := make(map[string]interface{})
	for k, v := range m {
		switch val := v.(type) {
		case string, int, int64, float64, bool:
			result[k] = fmt.Sprintf("%v", val)
		default:
			// For complex types, store as JSON
			jsonBytes, _ := json.Marshal(val)
			result[k] = string(jsonBytes)
		}
	}

	return result, nil
}

// ===========================================
// STREAM NAME GENERATORS
// ===========================================

// GetStreamName generates the stream name for different data types
func GetStreamName(dataType string, exchange Exchange, symbol Symbol) string {
	if exchange == "" {
		return fmt.Sprintf("%s:%s", dataType, symbol)
	}
	return fmt.Sprintf("%s:%s:%s", dataType, exchange, symbol)
}

// Stream name generators
func GetLiquidationStreamName(exchange Exchange, symbol Symbol) string {
	return GetStreamName("liquidations", exchange, symbol)
}

func GetMarketStreamName(exchange Exchange, symbol Symbol) string {
	return GetStreamName("market", exchange, symbol)
}

func GetOrderBookStreamName(exchange Exchange, symbol Symbol) string {
	return GetStreamName("orderbook", exchange, symbol)
}

func GetHeatmapStreamName(symbol Symbol) string {
	return fmt.Sprintf("heatmap:%s", symbol)
}

func GetHeatmapCacheKey(symbol Symbol, interval Interval) string {
	return fmt.Sprintf("heatmap:cache:%s:%s", symbol, interval)
}

// ===========================================
// VALIDATION METHODS
// ===========================================

// Validate checks if MarketSnapshot is valid
func (m *MarketSnapshot) Validate() error {
	if m.Exchange == "" {
		return fmt.Errorf("exchange is required")
	}
	if m.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if m.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp")
	}
	if m.MarkPrice <= 0 {
		return fmt.Errorf("invalid mark price")
	}
	return nil
}

// Validate checks if LiquidationEvent is valid
func (l *LiquidationEvent) Validate() error {
	if l.Exchange == "" {
		return fmt.Errorf("exchange is required")
	}
	if l.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if l.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp")
	}
	if l.Price <= 0 {
		return fmt.Errorf("invalid price")
	}
	if l.Quantity <= 0 {
		return fmt.Errorf("invalid quantity")
	}
	return nil
}

// Validate checks if HeatmapData is valid
func (h *HeatmapData) Validate() error {
	if h.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if h.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp")
	}
	if h.CurrentPrice <= 0 {
		return fmt.Errorf("invalid current price")
	}
	if len(h.Levels) == 0 {
		return fmt.Errorf("no liquidation levels")
	}
	return nil
}

// ===========================================
// HELPER FUNCTIONS
// ===========================================

// GetLiquidationType returns the liquidation type based on side
func (l *LiquidationEvent) GetLiquidationType() string {
	switch l.Side {
	case SideSell:
		return "LONG" // Long positions get liquidated with sell orders
	case SideBuy:
		return "SHORT" // Short positions get liquidated with buy orders
	default:
		if l.Side == SideLong {
			return "LONG"
		}
		return "SHORT"
	}
}

// GetEstimatedLeverage estimates the leverage used based on liquidation price
func (l *LiquidationEvent) GetEstimatedLeverage(markPrice float64) float64 {
	maintenanceMargin := 0.004 // 0.4% for Binance

	if l.GetLiquidationType() == "LONG" {
		if markPrice > 0 && l.Price < markPrice {
			return 1 / (1 - l.Price/markPrice + maintenanceMargin)
		}
	} else { // SHORT
		if markPrice > 0 && l.Price > markPrice {
			return 1 / (l.Price/markPrice - 1 + maintenanceMargin)
		}
	}

	return 0 // Unable to calculate
}

// CalculateIntensity calculates the intensity score for a liquidation level
func (ll *LiquidationLevel) CalculateIntensity(maxVolume float64) {
	if maxVolume > 0 {
		ll.Intensity = (ll.TotalVolume / maxVolume) * 100
	}
}

// IsSignificant determines if a liquidation level is significant
func (ll *LiquidationLevel) IsSignificant(threshold float64) bool {
	return ll.Intensity >= threshold
}

// GetIntervalDuration returns the duration for an interval
func GetIntervalDuration(interval Interval) time.Duration {
	switch interval {
	case Interval1s:
		return time.Second
	case Interval1m:
		return time.Minute
	case Interval5m:
		return 5 * time.Minute
	case Interval15m:
		return 15 * time.Minute
	case Interval1h:
		return time.Hour
	case Interval4h:
		return 4 * time.Hour
	case Interval1d:
		return 24 * time.Hour
	default:
		return time.Minute
	}
}

// RoundToInterval rounds a timestamp to the nearest interval
func RoundToInterval(timestamp int64, interval Interval) int64 {
	duration := GetIntervalDuration(interval)
	t := time.UnixMilli(timestamp)
	rounded := t.Truncate(duration)
	return rounded.UnixMilli()
}
