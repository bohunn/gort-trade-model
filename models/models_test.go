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
			name: "valid event",
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
			name: "missing exchange",
			event: LiquidationEvent{
				Symbol:    SymbolBTCUSDT,
				Timestamp: time.Now().UnixMilli(),
				Side:      SideLong,
				Price:     45000.0,
				Quantity:  1.5,
				Value:     67500.0,
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
				Value:     67500.0,
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

func TestCalculateLiquidationPrice(t *testing.T) {
	tests := []struct {
		name              string
		entryPrice        float64
		leverage          float64
		isLong            bool
		maintenanceMargin float64
		expected          float64
	}{
		{
			name:              "long position 10x leverage",
			entryPrice:        40000.0,
			leverage:          10.0,
			isLong:            true,
			maintenanceMargin: 0.005,
			expected:          36200.0, // 40000 * (1 - 1/10 + 0.005)
		},
		{
			name:              "short position 5x leverage",
			entryPrice:        40000.0,
			leverage:          5.0,
			isLong:            false,
			maintenanceMargin: 0.01,
			expected:          47600.0, // 40000 * (1 + 1/5 - 0.01)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateLiquidationPrice(tt.entryPrice, tt.leverage, tt.isLong, tt.maintenanceMargin)
			if result != tt.expected {
				t.Errorf("CalculateLiquidationPrice() = %v, expected %v", result, tt.expected)
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
			expected: "market:okx:ETHUSD",
		},
		{
			name:     "heatmap stream",
			function: func() string { return GetHeatmapStreamName(SymbolBTCUSDT) },
			expected: "heatmap:BTCUSDT",
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

	// Test round-trip
	var reconstructed LiquidationEvent
	err = FromStreamMessage(msg, &reconstructed)
	if err != nil {
		t.Fatalf("FromStreamMessage() error = %v", err)
	}

	if reconstructed.Exchange != event.Exchange {
		t.Errorf("Round-trip failed: Exchange = %v, expected %v", reconstructed.Exchange, event.Exchange)
	}
}
