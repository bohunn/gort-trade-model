# Gort Trade Model

A Go module providing shared data structures for cryptocurrency liquidation heatmap services.

## Installation

```bash
go get github.com/bohunn/gort-trade-model
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/bohunn/gort-trade-model"
)

func main() {
    // Create a liquidation event
    event := models.LiquidationEvent{
        Exchange:  models.ExchangeBinance,
        Symbol:    models.SymbolBTCUSDT,
        Timestamp: time.Now().UnixMilli(),
        Side:      models.SideLong,
        Price:     45000.0,
        Quantity:  1.5,
        Value:     67500.0,
        OrderType: models.OrderTypeLiquidation,
    }
    
    // Validate the event
    if err := event.Validate(); err != nil {
        fmt.Printf("Invalid event: %v\n", err)
        return
    }
    
    fmt.Printf("Liquidation: %+v\n", event)
}
```

## Supported Exchanges

- Binance
- OKX
- Bybit
- Coinbase
- Kraken
- Deribit
- Bitfinex

## Data Structures

### Core Types
- `LiquidationEvent` - Individual liquidation data
- `MarketSnapshot` - Current market state
- `PositionDistribution` - Position data at price levels
- `HeatmapData` - Aggregated liquidation heatmap
- `OrderBookSnapshot` - Order book state

### Enums
- `Exchange` - Supported exchanges
- `Symbol` - Trading pairs
- `OrderType` - Liquidation order types
- `Side` - Position sides (long/short)
- `Interval` - Time intervals for aggregation

## Validation

All major data structures include validation methods:

```go
event := &models.LiquidationEvent{...}
if err := event.Validate(); err != nil {
    // handle validation error
}
```

## Stream Integration

The module includes Redis Streams integration utilities:

```go
// Convert to stream message
msg, err := models.ToStreamMessage("liquidations", event)

// Get stream names
streamName := models.GetLiquidationStreamName(models.ExchangeBinance, models.SymbolBTCUSDT)
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

[MIT](LICENSE.md)