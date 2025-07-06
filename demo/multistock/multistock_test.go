package main

import (
	"math"
	"testing"
	"time"
)

func TestPriceDataStruct(t *testing.T) {
	price := PriceData{
		Price: 100.0,
		High:  102.0,
		Low:   98.0,
		Close: 100.5,
	}
	
	if price.Price != 100.0 {
		t.Errorf("Expected price 100.0, got %f", price.Price)
	}
	if price.High != 102.0 {
		t.Errorf("Expected high 102.0, got %f", price.High)
	}
	if price.Low != 98.0 {
		t.Errorf("Expected low 98.0, got %f", price.Low)
	}
	if price.Close != 100.5 {
		t.Errorf("Expected close 100.5, got %f", price.Close)
	}
}

func TestTradingIndicatorsStruct(t *testing.T) {
	indicators := TradingIndicators{
		ATR:        2.5,
		UpperBand:  105.0,
		MiddleBand: 100.0,
		LowerBand:  95.0,
		StdDev:     2.5,
	}
	
	if indicators.UpperBand <= indicators.MiddleBand {
		t.Error("Upper band should be greater than middle band")
	}
	if indicators.LowerBand >= indicators.MiddleBand {
		t.Error("Lower band should be less than middle band")
	}
	if indicators.ATR <= 0 {
		t.Error("ATR should be positive")
	}
	if indicators.StdDev <= 0 {
		t.Error("Standard deviation should be positive")
	}
}

func TestStockDataStruct(t *testing.T) {
	stockData := &StockData{
		Ticker:       "TEST",
		PriceHistory: make([]PriceData, 0),
		CurrentPrice: 100.0,
		LastUpdate:   time.Now(),
	}
	
	if stockData.Ticker != "TEST" {
		t.Errorf("Expected ticker TEST, got %s", stockData.Ticker)
	}
	if stockData.PriceHistory == nil {
		t.Error("Price history should not be nil")
	}
	if stockData.CurrentPrice != 100.0 {
		t.Errorf("Expected current price 100.0, got %f", stockData.CurrentPrice)
	}
}

func TestTradingBotStruct(t *testing.T) {
	bot := &TradingBot{
		stocks:          make(map[string]*StockData),
		riskPercent:     1.0,
		atrPeriod:       14,
		bollingerPeriod: 20,
	}
	
	if bot.stocks == nil {
		t.Error("Stocks map should not be nil")
	}
	if bot.riskPercent != 1.0 {
		t.Errorf("Expected risk percent 1.0, got %f", bot.riskPercent)
	}
	if bot.atrPeriod != 14 {
		t.Errorf("Expected ATR period 14, got %d", bot.atrPeriod)
	}
	if bot.bollingerPeriod != 20 {
		t.Errorf("Expected Bollinger period 20, got %d", bot.bollingerPeriod)
	}
}

func TestCalculateATR(t *testing.T) {
	bot := &TradingBot{atrPeriod: 3}
	
	priceHistory := []PriceData{
		{High: 102, Low: 98, Close: 100},
		{High: 104, Low: 99, Close: 103},
		{High: 105, Low: 101, Close: 102},
		{High: 103, Low: 99, Close: 101},
	}
	
	atr := bot.calculateATR(priceHistory)
	if atr <= 0 {
		t.Error("ATR should be positive")
	}
	
	// Test with insufficient data
	shortHistory := []PriceData{{High: 100, Low: 98, Close: 99}}
	atr = bot.calculateATR(shortHistory)
	if atr != 0 {
		t.Error("ATR should be 0 with insufficient data")
	}
	
	// Test edge case - exactly minimum data
	minHistory := []PriceData{
		{High: 102, Low: 98, Close: 100},
		{High: 104, Low: 99, Close: 103},
		{High: 105, Low: 101, Close: 102},
		{High: 103, Low: 99, Close: 101},
	}
	atr = bot.calculateATR(minHistory)
	if atr <= 0 {
		t.Error("ATR should be positive with minimum required data")
	}
}

func TestCalculateBollingerBands(t *testing.T) {
	bot := &TradingBot{bollingerPeriod: 3}
	
	priceHistory := []PriceData{
		{Close: 100},
		{Close: 102},
		{Close: 98},
	}
	
	upper, middle, lower, stdDev := bot.calculateBollingerBands(priceHistory)
	
	expectedMiddle := 100.0 // (100 + 102 + 98) / 3
	if middle != expectedMiddle {
		t.Errorf("Expected middle band %f, got %f", expectedMiddle, middle)
	}
	
	if stdDev <= 0 {
		t.Error("Standard deviation should be positive")
	}
	
	if upper <= middle || lower >= middle {
		t.Error("Upper band should be above middle, lower band below middle")
	}
	
	// Test with insufficient data
	shortHistory := []PriceData{{Close: 100}}
	upper, middle, lower, stdDev = bot.calculateBollingerBands(shortHistory)
	if upper != 0 || middle != 0 || lower != 0 || stdDev != 0 {
		t.Error("All values should be 0 with insufficient data")
	}
}

func TestCalculateIndicators(t *testing.T) {
	bot := &TradingBot{
		atrPeriod:       14,
		bollingerPeriod: 20,
	}
	
	// Create sufficient price history
	priceHistory := make([]PriceData, 25)
	basePrice := 100.0
	for i := 0; i < 25; i++ {
		price := basePrice + float64(i%5-2) // Creates some variation
		priceHistory[i] = PriceData{
			High:  price + 1,
			Low:   price - 1,
			Close: price,
		}
	}
	
	indicators := bot.calculateIndicators(priceHistory)
	
	if indicators.ATR <= 0 {
		t.Error("ATR should be positive")
	}
	if indicators.MiddleBand <= 0 {
		t.Error("Middle band should be positive")
	}
	if indicators.UpperBand <= indicators.MiddleBand {
		t.Error("Upper band should be above middle band")
	}
	if indicators.LowerBand >= indicators.MiddleBand {
		t.Error("Lower band should be below middle band")
	}
	
	// Test with insufficient data
	shortHistory := []PriceData{{High: 100, Low: 98, Close: 99}}
	indicators = bot.calculateIndicators(shortHistory)
	if indicators.ATR != 0 || indicators.MiddleBand != 0 {
		t.Error("Indicators should be zero with insufficient data")
	}
}

func TestStockDataInitialisation(t *testing.T) {
	stocks := make(map[string]*StockData)
	tickers := []string{"AAPL", "GOOGL", "MSFT"}
	
	for _, ticker := range tickers {
		stocks[ticker] = &StockData{
			Ticker:       ticker,
			PriceHistory: make([]PriceData, 0),
			LastUpdate:   time.Now(),
		}
	}
	
	if len(stocks) != 3 {
		t.Errorf("Expected 3 stocks, got %d", len(stocks))
	}
	
	for _, ticker := range tickers {
		if stocks[ticker] == nil {
			t.Errorf("Stock data for %s should not be nil", ticker)
		}
		if stocks[ticker].Ticker != ticker {
			t.Errorf("Expected ticker %s, got %s", ticker, stocks[ticker].Ticker)
		}
	}
}

func TestPriceHistoryManagement(t *testing.T) {
	stockData := &StockData{
		Ticker:       "TEST",
		PriceHistory: make([]PriceData, 0),
	}
	
	maxHistory := 20
	
	// Add more than max history
	for i := 0; i < 25; i++ {
		priceData := PriceData{
			Price: float64(100 + i),
			Close: float64(100 + i),
		}
		stockData.PriceHistory = append(stockData.PriceHistory, priceData)
		
		// Simulate trimming logic
		if len(stockData.PriceHistory) > maxHistory {
			stockData.PriceHistory = stockData.PriceHistory[len(stockData.PriceHistory)-maxHistory:]
		}
	}
	
	if len(stockData.PriceHistory) != maxHistory {
		t.Errorf("Expected %d price history entries, got %d", maxHistory, len(stockData.PriceHistory))
	}
	
	// Check that latest prices are kept
	latestPrice := stockData.PriceHistory[len(stockData.PriceHistory)-1].Price
	expectedLatest := 124.0 // 100 + 24
	if latestPrice != expectedLatest {
		t.Errorf("Expected latest price %f, got %f", expectedLatest, latestPrice)
	}
}

func TestATRCalculationEdgeCases(t *testing.T) {
	bot := &TradingBot{atrPeriod: 2}
	
	// Test with identical prices (no volatility)
	stablePrices := []PriceData{
		{High: 100, Low: 100, Close: 100},
		{High: 100, Low: 100, Close: 100},
		{High: 100, Low: 100, Close: 100},
	}
	
	atr := bot.calculateATR(stablePrices)
	if atr != 0 {
		t.Errorf("Expected ATR 0 for stable prices, got %f", atr)
	}
	
	// Test with extreme volatility
	volatilePrices := []PriceData{
		{High: 100, Low: 90, Close: 95},
		{High: 110, Low: 95, Close: 105},
		{High: 120, Low: 100, Close: 110},
	}
	
	atr = bot.calculateATR(volatilePrices)
	if atr <= 0 {
		t.Error("ATR should be positive for volatile prices")
	}
}

func TestBollingerBandsEdgeCases(t *testing.T) {
	bot := &TradingBot{bollingerPeriod: 3}
	
	// Test with identical prices
	stablePrices := []PriceData{
		{Close: 100},
		{Close: 100},
		{Close: 100},
	}
	
	upper, middle, lower, stdDev := bot.calculateBollingerBands(stablePrices)
	
	if middle != 100.0 {
		t.Errorf("Expected middle band 100.0, got %f", middle)
	}
	if stdDev != 0 {
		t.Errorf("Expected stdDev 0 for identical prices, got %f", stdDev)
	}
	if upper != middle || lower != middle {
		t.Error("Upper and lower bands should equal middle band when stdDev is 0")
	}
}

// Benchmark tests
func BenchmarkCalculateATR(b *testing.B) {
	bot := &TradingBot{atrPeriod: 14}
	
	priceHistory := make([]PriceData, 100)
	for i := 0; i < 100; i++ {
		price := 100.0 + float64(i%10)
		priceHistory[i] = PriceData{
			High:  price + 1,
			Low:   price - 1,
			Close: price,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bot.calculateATR(priceHistory)
	}
}

func BenchmarkCalculateBollingerBands(b *testing.B) {
	bot := &TradingBot{bollingerPeriod: 20}
	
	priceHistory := make([]PriceData, 100)
	for i := 0; i < 100; i++ {
		price := 100.0 + float64(i%10)
		priceHistory[i] = PriceData{Close: price}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bot.calculateBollingerBands(priceHistory)
	}
}

func BenchmarkCalculateIndicators(b *testing.B) {
	bot := &TradingBot{
		atrPeriod:       14,
		bollingerPeriod: 20,
	}
	
	priceHistory := make([]PriceData, 50)
	for i := 0; i < 50; i++ {
		price := 100.0 + math.Sin(float64(i)*0.1)*10 // Sine wave price movement
		priceHistory[i] = PriceData{
			High:  price + 1,
			Low:   price - 1,
			Close: price,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bot.calculateIndicators(priceHistory)
	}
}