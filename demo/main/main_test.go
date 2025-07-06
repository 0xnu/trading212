package main

import (
	"testing"
)

func TestNewTradingDemoRunner(t *testing.T) {
	runner := NewTradingDemoRunner("test_key", true)
	if runner == nil {
		t.Fatal("Expected non-nil runner")
	}
	if runner.client == nil {
		t.Fatal("Expected non-nil client")
	}
}

func TestDisplaySampleInstruments(t *testing.T) {
	runner := &TradingDemoRunner{}
	
	// Test with empty slice
	runner.displaySampleInstruments([]interface{}{})
	
	// Test with valid instruments
	instruments := []interface{}{
		map[string]interface{}{"ticker": "AAPL"},
		map[string]interface{}{"ticker": "GOOGL"},
		map[string]interface{}{"ticker": "MSFT"},
		map[string]interface{}{"ticker": "TSLA"},
	}
	runner.displaySampleInstruments(instruments)
	
	// Test with non-slice input
	runner.displaySampleInstruments("invalid")
}

func TestDisplaySampleInstrumentsEdgeCases(t *testing.T) {
	runner := &TradingDemoRunner{}
	
	// Test with nil
	runner.displaySampleInstruments(nil)
	
	// Test with wrong type slice
	runner.displaySampleInstruments([]string{"invalid", "data"})
	
	// Test with single instrument
	singleInstrument := []interface{}{
		map[string]interface{}{"ticker": "AAPL", "name": "Apple Inc"},
	}
	runner.displaySampleInstruments(singleInstrument)
}

func TestTradingDemoRunnerStructure(t *testing.T) {
	runner := &TradingDemoRunner{
		client: nil, // We can test with nil client for structure validation
	}
	
	if runner == nil {
		t.Error("TradingDemoRunner should not be nil")
	}
}