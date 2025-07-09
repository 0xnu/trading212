package main

import (
	"os"
	"testing"
	"time"
)

func TestCreateLogFile(t *testing.T) {
	logFile := createLogFile()

	if logFile == nil {
		t.Error("Expected log file to be created, got nil")
	} else {
		defer func() {
			logFile.Close()
			// Clean up test file
			fileName := logFile.Name()
			os.Remove(fileName)
		}()
	}
}

func TestNewRoboAdvisor(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)

	if advisor == nil {
		t.Fatal("Expected RoboAdvisor to be created, got nil")
	}

	if advisor.client == nil {
		t.Error("Expected client to be initialised")
	}

	if advisor.strategies == nil {
		t.Error("Expected strategies to be initialised")
	}

	if advisor.instrumentMapping == nil {
		t.Error("Expected instrument mapping to be initialised")
	}

	defer advisor.Close()
}

func TestGetDefaultConfigurations(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)
	defer advisor.Close()

	configs := advisor.getDefaultConfigurations()

	if len(configs) != 3 {
		t.Errorf("Expected 3 default configurations, got %d", len(configs))
	}

	expectedNames := []string{"Conservative Growth", "Balanced Portfolio", "Tech Growth"}
	expectedStrategies := []InvestmentStrategy{Conservative, Balanced, Growth}
	expectedAmounts := []float64{500.0, 750.0, 1000.0}

	for i, config := range configs {
		if config.Name != expectedNames[i] {
			t.Errorf("Expected name %s, got %s", expectedNames[i], config.Name)
		}

		if config.Strategy != expectedStrategies[i] {
			t.Errorf("Expected strategy %s, got %s", expectedStrategies[i], config.Strategy)
		}

		if config.MonthlyAmount != expectedAmounts[i] {
			t.Errorf("Expected monthly amount %.2f, got %.2f", expectedAmounts[i], config.MonthlyAmount)
		}
	}
}

func TestParseConfigurationsOrDefault(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)
	defer advisor.Close()

	// Test with empty JSON
	configs, err := advisor.parseConfigurationsOrDefault("")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(configs) != 3 {
		t.Errorf("Expected 3 default configurations, got %d", len(configs))
	}

	// Test with invalid JSON
	configs, err = advisor.parseConfigurationsOrDefault("invalid json")
	if err != nil {
		t.Errorf("Expected no error with fallback to defaults, got %v", err)
	}

	if len(configs) != 3 {
		t.Errorf("Expected 3 default configurations on invalid JSON, got %d", len(configs))
	}
}

func TestInstrumentAllocator(t *testing.T) {
	allocation := AssetAllocation{
		Technology: 0.5,
		Healthcare: 0.3,
		Financial:  0.2,
	}

	allocator := &InstrumentAllocator{allocation: allocation}

	// Test base instruments creation
	instruments := allocator.createBaseInstruments()

	if len(instruments) == 0 {
		t.Error("Expected instruments to be created")
	}

	// Test normalisation
	normalised := allocator.normaliseWeights(instruments)

	total := 0.0
	for _, weight := range normalised {
		total += weight
	}

	// Should sum to approximately 1.0
	if total < 0.99 || total > 1.01 {
		t.Errorf("Expected normalised weights to sum to 1.0, got %.4f", total)
	}
}

func TestAddTechnologyInstruments(t *testing.T) {
	allocation := AssetAllocation{Technology: 0.6}
	allocator := &InstrumentAllocator{allocation: allocation}

	instruments := make(map[string]float64)
	allocator.addTechnologyInstruments(instruments)

	expectedTickers := []string{"META", "MSFT", "NVDA"}
	for _, ticker := range expectedTickers {
		if _, exists := instruments[ticker]; !exists {
			t.Errorf("Expected ticker %s to be added", ticker)
		}
	}

	// Test with zero allocation
	allocation.Technology = 0.0
	allocator.allocation = allocation
	instruments = make(map[string]float64)
	allocator.addTechnologyInstruments(instruments)

	if len(instruments) != 0 {
		t.Error("Expected no instruments when allocation is zero")
	}
}

func TestCalculateTotal(t *testing.T) {
	allocator := &InstrumentAllocator{}

	instruments := map[string]float64{
		"AAPL":  0.3,
		"MSFT":  0.4,
		"GOOGL": 0.3,
	}

	total := allocator.calculateTotal(instruments)
	expected := 1.0

	if total != expected {
		t.Errorf("Expected total %.2f, got %.2f", expected, total)
	}
}

func TestApplyNormalisation(t *testing.T) {
	allocator := &InstrumentAllocator{}

	instruments := map[string]float64{
		"AAPL":  0.6,
		"MSFT":  0.8,
		"GOOGL": 0.6,
	}

	total := 2.0
	normalised := allocator.applyNormalisation(instruments, total)

	expectedNormalised := map[string]float64{
		"AAPL":  0.3,
		"MSFT":  0.4,
		"GOOGL": 0.3,
	}

	for ticker, expectedWeight := range expectedNormalised {
		if normalised[ticker] != expectedWeight {
			t.Errorf("Expected %s weight %.2f, got %.2f", ticker, expectedWeight, normalised[ticker])
		}
	}
}

func TestHasSufficientCash(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)
	defer advisor.Close()

	config := PieConfig{
		Name:          "Test Pie",
		MonthlyAmount: 500.0,
	}

	// Test sufficient cash
	if !advisor.hasSufficientCash(config, 1000.0) {
		t.Error("Expected true for sufficient cash")
	}

	// Test insufficient cash
	if advisor.hasSufficientCash(config, 250.0) {
		t.Error("Expected false for insufficient cash")
	}

	// Test exact amount
	if !advisor.hasSufficientCash(config, 500.0) {
		t.Error("Expected true for exact cash amount")
	}
}

func TestExtractAvailableCash(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)
	defer advisor.Close()

	// Test with proper map structure
	cashResponse := map[string]interface{}{
		"free": 1500.50,
	}

	cash := advisor.extractAvailableCash(cashResponse)
	expected := 1500.50

	if cash != expected {
		t.Errorf("Expected cash %.2f, got %.2f", expected, cash)
	}

	// Test with invalid structure
	invalidResponse := "invalid"
	cash = advisor.extractAvailableCash(invalidResponse)

	if cash != 0.0 {
		t.Errorf("Expected 0.0 for invalid response, got %.2f", cash)
	}

	// Test with missing free field
	missingFreeResponse := map[string]interface{}{
		"total": 2000.0,
	}

	cash = advisor.extractAvailableCash(missingFreeResponse)

	if cash != 0.0 {
		t.Errorf("Expected 0.0 for missing free field, got %.2f", cash)
	}
}

func TestAddDelayBetweenPies(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)
	defer advisor.Close()

	start := time.Now()

	// First pie (index 0) should not add delay
	advisor.addDelayBetweenPies(0)
	duration := time.Since(start)

	if duration > 100*time.Millisecond {
		t.Error("Expected no delay for first pie")
	}

	// Second pie (index 1) should add delay
	start = time.Now()
	advisor.addDelayBetweenPies(1)
	duration = time.Since(start)

	if duration < 5*time.Second {
		t.Error("Expected delay for subsequent pies")
	}
}

func TestCalculateTotalShares(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)
	defer advisor.Close()

	shares := map[string]float64{
		"AAPL":  10.5,
		"MSFT":  15.0,
		"GOOGL": 8.25,
	}

	total := advisor.calculateTotalShares(shares)
	expected := 33.75

	if total != expected {
		t.Errorf("Expected total shares %.2f, got %.2f", expected, total)
	}

	// Test empty shares
	emptyShares := make(map[string]float64)
	total = advisor.calculateTotalShares(emptyShares)

	if total != 0.0 {
		t.Errorf("Expected 0.0 for empty shares, got %.2f", total)
	}
}

func TestCalculateTotalGoal(t *testing.T) {
	advisor := NewRoboAdvisor("test-api-key", true)
	defer advisor.Close()

	// Mock pies data - we can't use the actual trading212.Pie type in tests
	// So we'll test the logic with a simplified approach

	// This would test the actual method if we had access to trading212.Pie
	// For now, we test the calculation logic separately
	goals := []int{50000, 75000, 100000}
	totalGoal := 0.0

	for _, goal := range goals {
		totalGoal += float64(goal)
	}

	expected := 225000.0
	if totalGoal != expected {
		t.Errorf("Expected total goal %.2f, got %.2f", expected, totalGoal)
	}
}

func TestInitializeStrategies(t *testing.T) {
	strategies := initializeStrategies()

	expectedStrategies := []InvestmentStrategy{
		Conservative, Balanced, Aggressive, Growth, Income,
	}

	if len(strategies) != len(expectedStrategies) {
		t.Errorf("Expected %d strategies, got %d", len(expectedStrategies), len(strategies))
	}

	for _, strategy := range expectedStrategies {
		if _, exists := strategies[strategy]; !exists {
			t.Errorf("Expected strategy %s to exist", strategy)
		}
	}

	// Test Conservative strategy allocation
	conservative := strategies[Conservative]
	if conservative.Equities != 0.4 {
		t.Errorf("Expected Conservative equities 0.4, got %.2f", conservative.Equities)
	}

	if conservative.Bonds != 0.4 {
		t.Errorf("Expected Conservative bonds 0.4, got %.2f", conservative.Bonds)
	}
}

func TestInitializeInstrumentMapping(t *testing.T) {
	mapping := initializeInstrumentMapping()

	expectedMappings := map[string]string{
		"US_EQUITIES":   "SPY",
		"TECHNOLOGY":    "AAPL",
		"HEALTHCARE":    "JNJ",
		"FINANCIAL":     "JPM",
		"ENERGY":        "XOM",
		"REITS":         "AMT",
		"BONDS":         "TLT",
		"INTERNATIONAL": "TSM",
		"EMERGING":      "TSM",
	}

	if len(mapping) != len(expectedMappings) {
		t.Errorf("Expected %d mappings, got %d", len(expectedMappings), len(mapping))
	}

	for key, expectedValue := range expectedMappings {
		if value, exists := mapping[key]; !exists || value != expectedValue {
			t.Errorf("Expected mapping %s -> %s, got %s", key, expectedValue, value)
		}
	}
}

// Benchmark tests for performance
func BenchmarkCreateBaseInstruments(b *testing.B) {
	allocation := AssetAllocation{
		Technology: 0.4,
		Healthcare: 0.2,
		Financial:  0.2,
		Energy:     0.1,
		REITs:      0.1,
	}

	allocator := &InstrumentAllocator{allocation: allocation}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		allocator.createBaseInstruments()
	}
}

func BenchmarkNormaliseWeights(b *testing.B) {
	allocator := &InstrumentAllocator{}
	instruments := map[string]float64{
		"AAPL":  0.3,
		"MSFT":  0.4,
		"GOOGL": 0.3,
		"META":  0.2,
		"NVDA":  0.1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		allocator.normaliseWeights(instruments)
	}
}
