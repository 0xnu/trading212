package main

import (
	"errors"
	"github.com/0xnu/trading212"
	"testing"
	"time"
)

func TestNewTradingBot(t *testing.T) {
	bot := NewTradingBot("test_key", true, "NVDA", 1.0)

	if bot.ticker != "NVDA" {
		t.Errorf("Expected ticker NVDA, got %s", bot.ticker)
	}
	if bot.riskPercent != 1.0 {
		t.Errorf("Expected risk percent 1.0, got %f", bot.riskPercent)
	}
	if bot.client == nil {
		t.Error("Client should not be nil")
	}
}

func TestCalculateBollingerBands(t *testing.T) {
	bot := &TradingBot{}
	currentPrice := 100.0

	upper, lower := bot.calculateBollingerBands(currentPrice)

	expectedUpper := 104.0 // 100 + (2 * 2.0)
	expectedLower := 96.0  // 100 - (2 * 2.0)

	if upper != expectedUpper {
		t.Errorf("Expected upper band %f, got %f", expectedUpper, upper)
	}
	if lower != expectedLower {
		t.Errorf("Expected lower band %f, got %f", expectedLower, lower)
	}
}

func TestCalculatePositionSize(t *testing.T) {
	bot := &TradingBot{riskPercent: 2.0}
	portfolioValue := 10000.0
	currentPrice := 100.0

	size := bot.calculatePositionSize(portfolioValue, currentPrice)
	expected := 2 // (10000 * 2 / 100) / 100 = 2

	if size != expected {
		t.Errorf("Expected position size %d, got %d", expected, size)
	}

	// Test minimum position size
	smallPortfolio := 50.0
	size = bot.calculatePositionSize(smallPortfolio, currentPrice)
	if size != 1 {
		t.Errorf("Expected minimum position size 1, got %d", size)
	}
}

func TestCalculatePortfolioValue(t *testing.T) {
	bot := &TradingBot{}
	positions := []trading212.Position{
		{Value: 1000.0},
		{Value: 1500.0},
		{Value: 2500.0},
	}

	value := bot.calculatePortfolioValue(positions)
	expected := 5000.0

	if value != expected {
		t.Errorf("Expected portfolio value %f, got %f", expected, value)
	}

	// Test empty positions
	emptyPositions := []trading212.Position{}
	value = bot.calculatePortfolioValue(emptyPositions)
	if value != 0.0 {
		t.Errorf("Expected portfolio value 0.0 for empty positions, got %f", value)
	}
}

func TestHasExistingPosition(t *testing.T) {
	bot := &TradingBot{ticker: "NVDA"}

	positions := []trading212.Position{
		{Ticker: "AAPL", Quantity: 10},
		{Ticker: "NVDA", Quantity: 5},
	}

	if !bot.hasExistingPosition(positions) {
		t.Error("Expected to find existing position")
	}

	positionsWithoutNVDA := []trading212.Position{
		{Ticker: "AAPL", Quantity: 10},
		{Ticker: "GOOGL", Quantity: 5},
	}

	if bot.hasExistingPosition(positionsWithoutNVDA) {
		t.Error("Expected no existing position")
	}

	// Test empty positions
	emptyPositions := []trading212.Position{}
	if bot.hasExistingPosition(emptyPositions) {
		t.Error("Expected no existing position for empty slice")
	}
}

func TestShouldBuyAndSell(t *testing.T) {
	bot := &TradingBot{}

	// Test buy conditions
	if !bot.shouldBuy(false, 95.0, 100.0) {
		t.Error("Should buy when no position and price below lower band")
	}
	if bot.shouldBuy(true, 95.0, 100.0) {
		t.Error("Should not buy when already have position")
	}
	if bot.shouldBuy(false, 105.0, 100.0) {
		t.Error("Should not buy when price above lower band")
	}

	// Test sell conditions
	if !bot.shouldSell(true, 105.0, 100.0) {
		t.Error("Should sell when have position and price above upper band")
	}
	if bot.shouldSell(false, 105.0, 100.0) {
		t.Error("Should not sell when no position")
	}
	if bot.shouldSell(true, 95.0, 100.0) {
		t.Error("Should not sell when price below upper band")
	}
}

func TestShouldRetry(t *testing.T) {
	bot := &TradingBot{}
	config := RetryConfig{maxRetries: 3, baseDelay: time.Millisecond}

	// Test rate limit error
	rateLimitErr := errors.New("TooManyRequests")
	if !bot.shouldRetry(rateLimitErr, 0, config) {
		t.Error("Should retry on rate limit error")
	}

	// Test max retries exceeded
	if bot.shouldRetry(rateLimitErr, 2, config) {
		t.Error("Should not retry when max retries reached")
	}

	// Test non-rate limit error
	otherErr := errors.New("Other error")
	if bot.shouldRetry(otherErr, 0, config) {
		t.Error("Should not retry on non-rate limit error")
	}
}

func TestExecuteTrade(t *testing.T) {
	bot := &TradingBot{}

	// Test no action case
	result := bot.executeTrade(false, 100.0, 105.0, 95.0, 1)
	if !result {
		t.Error("Should return true when no trade action needed")
	}
}

func TestRetryConfig(t *testing.T) {
	config := RetryConfig{
		maxRetries: 5,
		baseDelay:  time.Second,
	}

	if config.maxRetries != 5 {
		t.Errorf("Expected maxRetries 5, got %d", config.maxRetries)
	}
	if config.baseDelay != time.Second {
		t.Errorf("Expected baseDelay 1s, got %v", config.baseDelay)
	}
}

// Benchmark tests
func BenchmarkCalculateBollingerBands(b *testing.B) {
	bot := &TradingBot{}
	currentPrice := 100.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bot.calculateBollingerBands(currentPrice)
	}
}

func BenchmarkCalculatePositionSize(b *testing.B) {
	bot := &TradingBot{riskPercent: 1.0}
	portfolioValue := 10000.0
	currentPrice := 100.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bot.calculatePositionSize(portfolioValue, currentPrice)
	}
}
