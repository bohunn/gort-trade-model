package models

import (
	"testing"
	"time"
)

func TestLiquidationEventValidation(t *testing.T) {
	tests := []struct {
		name    string
		event   LiquidationEvent
		wantErr bool
	}{
		{
			name: "valid event with all fields",
			event: LiquidationEvent{
				Exchange:  ExchangeBinance,
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				Side:      SideLong,
				Price:     45000.0,
				Quantity:  1.5,
				Value:     67500.0,
				OrderType: OrderTypeLiquidation,
			},
			wantErr: false,
		},
		{
			name: "valid event without value field",
			event: LiquidationEvent{
				Exchange:  ExchangeBinance,
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				Side:      SideSell, // Binance format
				Price:     45000.0,
				Quantity:  1.5,
				OrderType: OrderTypeLiquidation,
			},
			wantErr: false,
		},
		{
			name: "missing exchange",
			event: LiquidationEvent{
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				Side:      SideLong,
				Price:     45000.0,
				Quantity:  1.5,
				OrderType: OrderTypeLiquidation,
			},
			wantErr: true,
		},
		{
			name: "missing symbol",
			event: LiquidationEvent{
				Exchange:  ExchangeBinance,
				Timestamp: time.Now().UnixMilli(),
				Side:      SideLong,
				Price:     45000.0,
				Quantity:  1.5,
				OrderType: OrderTypeLiquidation,
			},
			wantErr: true,
		},
		{
			name: "invalid timestamp",
			event: LiquidationEvent{
				Exchange:  ExchangeBinance,
				Symbol:    SymbolBTCUSDT,
				Timestamp: 0,
				Side:      SideLong,
				Price:     45000.0,
				Quantity:  1.5,
				OrderType: OrderTypeLiquidation,
			},
			wantErr: true,
		},
		{
			name: "invalid price",
			event: LiquidationEvent{
				Exchange:  ExchangeBinance,
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				Side:      SideLong,
				Price:     -1000.0, // negative price
				Quantity:  1.5,
				OrderType: OrderTypeLiquidation,
			},
			wantErr: true,
		},
		{
			name: "invalid quantity",
			event: LiquidationEvent{
				Exchange:  ExchangeBinance,
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				Side:      SideLong,
				Price:     45000.0,
				Quantity:  0, // zero quantity
				OrderType: OrderTypeLiquidation,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("LiquidationEvent.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarketSnapshotValidation(t *testing.T) {
	tests := []struct {
		name    string
		market  MarketSnapshot
		wantErr bool
	}{
		{
			name: "valid market snapshot",
			market: MarketSnapshot{
				Exchange:  ExchangeBinance,
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				MarkPrice: 45000.0,
			},
			wantErr: false,
		},
		{
			name: "missing exchange",
			market: MarketSnapshot{
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				MarkPrice: 45000.0,
			},
			wantErr: true,
		},
		{
			name: "invalid mark price",
			market: MarketSnapshot{
				Exchange:  ExchangeBinance,
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				MarkPrice: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.market.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarketSnapshot.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHeatmapDataValidation(t *testing.T) {
	tests := []struct {
		name    string
		heatmap HeatmapData
		wantErr bool
	}{
		{
			name: "valid heatmap",
			heatmap: HeatmapData{
				Symbol:       SymbolBTCUSDT,
				Timestamp:    time.Now().UnixMilli(),
				CurrentPrice: 45000.0,
				Levels: []LiquidationLevel{
					{
						Price:       44000.0,
						TotalVolume: 100000.0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing symbol",
			heatmap: HeatmapData{
				Timestamp:    time.Now().UnixMilli(),
				CurrentPrice: 45000.0,
				Levels: []LiquidationLevel{
					{Price: 44000.0},
				},
			},
			wantErr: true,
		},
		{
			name: "no levels",
			heatmap: HeatmapData{
				Symbol:       SymbolBTCUSDT,
				Timestamp:    time.Now().UnixMilli(),
				CurrentPrice: 45000.0,
				Levels:       []LiquidationLevel{},
			},
			wantErr: true,
		},
		{
			name: "invalid current price",
			heatmap: HeatmapData{
				Symbol:       SymbolBTCUSDT,
				Timestamp:    time.Now().UnixMilli(),
				CurrentPrice: 0,
				Levels: []LiquidationLevel{
					{Price: 44000.0},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.heatmap.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("HeatmapData.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLiquidationType(t *testing.T) {
	tests := []struct {
		name     string
		event    LiquidationEvent
		expected string
	}{
		{
			name: "long liquidation with SideLong",
			event: LiquidationEvent{
				Side: SideLong,
			},
			expected: "LONG",
		},
		{
			name: "short liquidation with SideShort",
			event: LiquidationEvent{
				Side: SideShort,
			},
			expected: "SHORT",
		},
		{
			name: "long liquidation with SELL order",
			event: LiquidationEvent{
				Side: SideSell,
			},
			expected: "LONG",
		},
		{
			name: "short liquidation with BUY order",
			event: LiquidationEvent{
				Side: SideBuy,
			},
			expected: "SHORT",
		},
		{
			name: "SELL string",
			event: LiquidationEvent{
				Side: "SELL",
			},
			expected: "LONG",
		},
		{
			name: "BUY string",
			event: LiquidationEvent{
				Side: "BUY",
			},
			expected: "SHORT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.GetLiquidationType()
			if result != tt.expected {
				t.Errorf("GetLiquidationType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGetEstimatedLeverage(t *testing.T) {
	tests := []struct {
		name      string
		event     LiquidationEvent
		markPrice float64
		expected  float64
		valid     bool
	}{
		{
			name: "long liquidation valid",
			event: LiquidationEvent{
				Side:  SideSell, // Long liquidation
				Price: 36000.0,  // Liquidation price
			},
			markPrice: 40000.0,
			expected:  10.0, // Approximately 10x leverage
			valid:     true,
		},
		{
			name: "short liquidation valid",
			event: LiquidationEvent{
				Side:  SideBuy, // Short liquidation
				Price: 44000.0, // Liquidation price
			},
			markPrice: 40000.0,
			expected:  10.0, // Approximately 10x leverage
			valid:     true,
		},
		{
			name: "invalid long liquidation",
			event: LiquidationEvent{
				Side:  SideSell,
				Price: 45000.0, // Price above mark price for long
			},
			markPrice: 40000.0,
			expected:  0,
			valid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.GetEstimatedLeverage(tt.markPrice)
			if tt.valid {
				// Allow some tolerance for leverage calculation
				tolerance := 0.5
				if result < tt.expected-tolerance || result > tt.expected+tolerance {
					t.Errorf("GetEstimatedLeverage() = %v, expected around %v", result, tt.expected)
				}
			} else {
				if result != 0 {
					t.Errorf("GetEstimatedLeverage() = %v, expected 0 for invalid case", result)
				}
			}
		})
	}
}

func TestCalculateIntensity(t *testing.T) {
	level := LiquidationLevel{
		Price:       45000.0,
		TotalVolume: 50000.0,
	}

	level.CalculateIntensity(100000.0)
	if level.Intensity != 50.0 {
		t.Errorf("CalculateIntensity() = %v, expected 50.0", level.Intensity)
	}

	// Test with zero max volume
	level.CalculateIntensity(0)
	if level.Intensity != 0 {
		t.Errorf("CalculateIntensity() with zero max = %v, expected 0", level.Intensity)
	}
}

func TestIsSignificant(t *testing.T) {
	tests := []struct {
		name      string
		level     LiquidationLevel
		threshold float64
		expected  bool
	}{
		{
			name:      "significant level",
			level:     LiquidationLevel{Intensity: 75.0},
			threshold: 50.0,
			expected:  true,
		},
		{
			name:      "not significant",
			level:     LiquidationLevel{Intensity: 25.0},
			threshold: 50.0,
			expected:  false,
		},
		{
			name:      "exactly at threshold",
			level:     LiquidationLevel{Intensity: 50.0},
			threshold: 50.0,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.level.IsSignificant(tt.threshold)
			if result != tt.expected {
				t.Errorf("IsSignificant() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestStreamNameGeneration(t *testing.T) {
	tests := []struct {
		name     string
		function func() string
		expected string
	}{
		{
			name:     "liquidation stream",
			function: func() string { return GetLiquidationStreamName(ExchangeBinance, SymbolBTCUSDT) },
			expected: "liquidations:binance:BTCUSDT",
		},
		{
			name:     "market stream",
			function: func() string { return GetMarketStreamName(ExchangeOKX, SymbolETHUSDT) },
			expected: "market:okx:ETHUSDT",
		},
		{
			name:     "orderbook stream",
			function: func() string { return GetOrderBookStreamName(ExchangeBybit, SymbolBNBUSDT) },
			expected: "orderbook:bybit:BNBUSDT",
		},
		{
			name:     "heatmap stream",
			function: func() string { return GetHeatmapStreamName(SymbolBTCUSDT) },
			expected: "heatmap:BTCUSDT",
		},
		{
			name:     "heatmap cache key",
			function: func() string { return GetHeatmapCacheKey(SymbolBTCUSDT, Interval1m) },
			expected: "heatmap:cache:BTCUSDT:1m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			if result != tt.expected {
				t.Errorf("Stream name = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGetIntervalDuration(t *testing.T) {
	tests := []struct {
		name     string
		interval Interval
		expected time.Duration
	}{
		{
			name:     "1 second interval",
			interval: Interval1s,
			expected: time.Second,
		},
		{
			name:     "1 minute interval",
			interval: Interval1m,
			expected: time.Minute,
		},
		{
			name:     "5 minute interval",
			interval: Interval5m,
			expected: 5 * time.Minute,
		},
		{
			name:     "1 hour interval",
			interval: Interval1h,
			expected: time.Hour,
		},
		{
			name:     "4 hour interval",
			interval: Interval4h,
			expected: 4 * time.Hour,
		},
		{
			name:     "1 day interval",
			interval: Interval1d,
			expected: 24 * time.Hour,
		},
		{
			name:     "unknown interval defaults to minute",
			interval: Interval("unknown"),
			expected: time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetIntervalDuration(tt.interval)
			if result != tt.expected {
				t.Errorf("GetIntervalDuration() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestRoundToInterval(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 34, 56, 789000000, time.UTC)
	tests := []struct {
		name      string
		timestamp int64
		interval  Interval
		expected  time.Time
	}{
		{
			name:      "round to minute",
			timestamp: baseTime.UnixMilli(),
			interval:  Interval1m,
			expected:  time.Date(2024, 1, 1, 12, 34, 0, 0, time.UTC),
		},
		{
			name:      "round to 5 minutes",
			timestamp: baseTime.UnixMilli(),
			interval:  Interval5m,
			expected:  time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
		},
		{
			name:      "round to hour",
			timestamp: baseTime.UnixMilli(),
			interval:  Interval1h,
			expected:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "round to day",
			timestamp: baseTime.UnixMilli(),
			interval:  Interval1d,
			expected:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundToInterval(tt.timestamp, tt.interval)
			expected := tt.expected.UnixMilli()
			if result != expected {
				t.Errorf("RoundToInterval() = %v, expected %v",
					time.UnixMilli(result), time.UnixMilli(expected))
			}
		})
	}
}

func TestToStreamMessage(t *testing.T) {
	event := LiquidationEvent{
		Exchange:  ExchangeBinance,
		Symbol:    SymbolBTCUSDT,
		Timestamp: 1234567890,
		Side:      SideLong,
		Price:     45000.0,
		Quantity:  1.5,
		Value:     67500.0,
		OrderType: OrderTypeLiquidation,
	}

	msg, err := ToStreamMessage("test-stream", event)
	if err != nil {
		t.Fatalf("ToStreamMessage() error = %v", err)
	}

	if msg.Stream != "test-stream" {
		t.Errorf("Stream name = %v, expected %v", msg.Stream, "test-stream")
	}

	if msg.Data == nil {
		t.Error("Data should not be nil")
	}

	// Check that timestamp is set
	if msg.Timestamp <= 0 {
		t.Error("Timestamp should be set")
	}
}

func TestStructToMapConversion(t *testing.T) {
	// Test with nested structure
	heatmap := HeatmapData{
		Symbol:       SymbolBTCUSDT,
		Timestamp:    1234567890,
		CurrentPrice: 45000.0,
		Interval:     Interval1m,
		Levels: []LiquidationLevel{
			{
				Price:             44000.0,
				LongLiquidations:  100000.0,
				ShortLiquidations: 50000.0,
				TotalVolume:       150000.0,
				Intensity:         75.0,
			},
		},
		Summary: HeatmapSummary{
			TotalLongLiquidations:  1000000.0,
			TotalShortLiquidations: 500000.0,
			SignificantLevels:      10,
			CriticalZones: []CriticalZone{
				{
					PriceStart: 43000.0,
					PriceEnd:   44000.0,
					Type:       "long",
					Intensity:  80.0,
					Volume:     200000.0,
				},
			},
		},
	}

	msg, err := ToStreamMessage("heatmap-stream", heatmap)
	if err != nil {
		t.Fatalf("ToStreamMessage() error = %v", err)
	}

	// Verify complex fields are serialized as JSON strings
	if _, ok := msg.Data["levels"].(string); !ok {
		t.Error("Levels should be serialized as JSON string")
	}

	if _, ok := msg.Data["summary"].(string); !ok {
		t.Error("Summary should be serialized as JSON string")
	}
}

func BenchmarkToStreamMessage(b *testing.B) {
	event := LiquidationEvent{
		Exchange:  ExchangeBinance,
		Symbol:    SymbolBTCUSDT,
		Timestamp: time.Now().UnixMilli(),
		Side:      SideLong,
		Price:     45000.0,
		Quantity:  1.5,
		Value:     67500.0,
		OrderType: OrderTypeLiquidation,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ToStreamMessage("bench-stream", event)
	}
}

func BenchmarkValidation(b *testing.B) {
	event := LiquidationEvent{
		Exchange:  ExchangeBinance,
		Symbol:    SymbolBTCUSDT,
		Timestamp: time.Now().UnixMilli(),
		Side:      SideLong,
		Price:     45000.0,
		Quantity:  1.5,
		Value:     67500.0,
		OrderType: OrderTypeLiquidation,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = event.Validate()
	}
}
