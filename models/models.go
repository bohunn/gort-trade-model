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
	SymbolETHUSDT Symbol = "ETHUSD"
	SymbolBTCUSD  Symbol = "BTCUSD"
	// Add more as needed
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
)

// Interval represents time intervals for aggregation
type Interval string

const (
	Interval1m  Interval = "1m"
	Interval5m  Interval = "5m"
	Interval15m Interval = "15m"
	Interval1h  Interval = "1h"
	Interval4h  Interval = "4h"
	Interval1d  Interval = "1d"
)

// MarketSnapshot represents current market state
type MarketSnapshot struct {
	Exchange        Exchange `json:"exchange"`
	Symbol          Symbol   `json:"symbol"`
	Timestamp       int64    `json:"timestamp"`
	MarkPrice       float64  `json:"mark_price"`
	IndexPrice      float64  `json:"index_price"`
	FundingRate     float64  `json:"funding_rate"`
	OpenInterest    float64  `json:"open_interest"` // in USD
	Volume24h       float64  `json:"volume_24h"`    // in USD
	NextFundingTime int64    `json:"next_funding_time"`
}

// LiquidationEvent represents a single liquidation
type LiquidationEvent struct {
	Exchange  Exchange  `json:"exchange"`
	Symbol    Symbol    `json:"symbol"`
	Timestamp int64     `json:"timestamp"`
	Side      Side      `json:"side"`
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Value     float64   `json:"value"` // USD value
	OrderType OrderType `json:"order_type"`
}

// PositionDistribution represents positions at a price level
type PositionDistribution struct {
	Exchange       Exchange        `json:"exchange"`
	Symbol         Symbol          `json:"symbol"`
	Timestamp      int64           `json:"timestamp"`
	PriceLevel     float64         `json:"price_level"`
	LongPositions  PositionSummary `json:"long_positions"`
	ShortPositions PositionSummary `json:"short_positions"`
}

// PositionSummary contains aggregated position data
type PositionSummary struct {
	Count       int     `json:"count"`
	Volume      float64 `json:"volume"` // in USD
	AvgLeverage float64 `json:"avg_leverage"`
}

// OrderBookSnapshot represents order book state
type OrderBookSnapshot struct {
	Exchange   Exchange    `json:"exchange"`
	Symbol     Symbol      `json:"symbol"`
	Timestamp  int64       `json:"timestamp"`
	Bids       []PriceSize `json:"bids"`
	Asks       []PriceSize `json:"asks"`
	SequenceID int64       `json:"sequence_id,omitempty"`
}

// PriceSize represents a price and size tuple
type PriceSize struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// HeatmapData represents aggregated liquidation heatmap
type HeatmapData struct {
	Symbol            Symbol             `json:"symbol"`
	Timestamp         int64              `json:"timestamp"`
	Interval          Interval           `json:"interval"`
	CurrentPrice      float64            `json:"current_price"`
	LiquidationLevels []LiquidationLevel `json:"liquidation_levels"`
	Metadata          HeatmapMetadata    `json:"metadata"`
}

// LiquidationLevel represents liquidations at a specific price
type LiquidationLevel struct {
	Price                       float64    `json:"price"`
	CumulativeLongLiquidations  float64    `json:"cumulative_long_liquidations"`  // USD
	CumulativeShortLiquidations float64    `json:"cumulative_short_liquidations"` // USD
	EstimatedImpact             float64    `json:"estimated_impact"`              // percentage
	Exchanges                   []Exchange `json:"exchanges"`
}

// HeatmapMetadata contains metadata about the heatmap
type HeatmapMetadata struct {
	ExchangesCovered []Exchange `json:"exchanges_covered"`
	DataCompleteness float64    `json:"data_completeness"` // 0-100
	CalculationTime  int64      `json:"calculation_time"`  // ms
	LastUpdate       int64      `json:"last_update"`
}

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

// FromStreamMessage reconstructs a model from StreamMessage
func FromStreamMessage(msg *StreamMessage, v interface{}) error {
	// Reconstruct the original map
	originalMap := make(map[string]interface{})
	for k, v := range msg.Data {
		str, ok := v.(string)
		if !ok {
			originalMap[k] = v
			continue
		}

		// Try to parse JSON for complex types
		if (str[0] == '{' || str[0] == '[') && json.Valid([]byte(str)) {
			var parsed interface{}
			if err := json.Unmarshal([]byte(str), &parsed); err == nil {
				originalMap[k] = parsed
				continue
			}
		}

		// Try to parse as number
		var num float64
		if _, err := fmt.Sscanf(str, "%f", &num); err == nil {
			originalMap[k] = num
		} else {
			originalMap[k] = str
		}
	}

	// Marshal and unmarshal to target type
	data, err := json.Marshal(originalMap)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// GetStreamName generates the stream name for different data types
func GetStreamName(dataType string, exchange Exchange, symbol Symbol) string {
	return fmt.Sprintf("%s:%s:%s", dataType, exchange, symbol)
}

// Stream name generators
func GetLiquidationStreamName(exchange Exchange, symbol Symbol) string {
	return GetStreamName("liquidations", exchange, symbol)
}

func GetPositionStreamName(exchange Exchange, symbol Symbol) string {
	return GetStreamName("positions", exchange, symbol)
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

// Validation methods

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
	if l.Value <= 0 {
		return fmt.Errorf("invalid value")
	}
	if l.Side != SideLong && l.Side != SideShort {
		return fmt.Errorf("invalid side: %s", l.Side)
	}
	return nil
}

// Validate checks if PositionDistribution is valid
func (p *PositionDistribution) Validate() error {
	if p.Exchange == "" {
		return fmt.Errorf("exchange is required")
	}
	if p.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if p.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp")
	}
	if p.PriceLevel <= 0 {
		return fmt.Errorf("invalid price level")
	}
	return nil
}

// Helper functions for calculations

// CalculateLiquidationPrice calculates the liquidation price for a position
func CalculateLiquidationPrice(entryPrice, leverage float64, isLong bool, maintenanceMargin float64) float64 {
	if isLong {
		// Long liquidation: price drops
		// Liquidation Price = Entry Price × (1 - 1/leverage + maintenance margin)
		return entryPrice * (1 - 1/leverage + maintenanceMargin)
	}
	// Short liquidation: price rises
	// Liquidation Price = Entry Price × (1 + 1/leverage - maintenance margin)
	return entryPrice * (1 + 1/leverage - maintenanceMargin)
}

// EstimatePriceImpact estimates the price impact of liquidations
func EstimatePriceImpact(liquidationVolume, marketDepth float64) float64 {
	if marketDepth == 0 {
		return 0
	}
	// Simplified model: impact = liquidation_volume / market_depth * impact_factor
	impactFactor := 0.1 // 10% impact per 100% of depth consumed
	return (liquidationVolume / marketDepth) * impactFactor * 100
}

// GetIntervalDuration returns the duration for an interval
func GetIntervalDuration(interval Interval) time.Duration {
	switch interval {
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
