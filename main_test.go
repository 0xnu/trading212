package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient tests the client creation
func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		demo   bool
		want   string
	}{
		{
			name:   "demo client",
			apiKey: "test-api-key",
			demo:   true,
			want:   "https://demo.trading212.com",
		},
		{
			name:   "live client",
			apiKey: "test-api-key",
			demo:   false,
			want:   "https://live.trading212.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.demo)
			if client.host != tt.want {
				t.Errorf("NewClient() host = %v, want %v", client.host, tt.want)
			}
			if client.apiKey != tt.apiKey {
				t.Errorf("NewClient() apiKey = %v, want %v", client.apiKey, tt.apiKey)
			}
			if client.httpClient == nil {
				t.Error("NewClient() httpClient is nil")
			}
		})
	}
}

// TestValidateTimeValidity tests time validity validation
func TestValidateTimeValidity(t *testing.T) {
	tests := []struct {
		name         string
		timeValidity string
		wantErr      bool
	}{
		{
			name:         "valid GTC",
			timeValidity: "GTC",
			wantErr:      false,
		},
		{
			name:         "valid DAY",
			timeValidity: "DAY",
			wantErr:      false,
		},
		{
			name:         "invalid value",
			timeValidity: "INVALID",
			wantErr:      true,
		},
		{
			name:         "empty value",
			timeValidity: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTimeValidity(tt.timeValidity)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTimeValidity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateDate tests date validation
func TestValidateDate(t *testing.T) {
	tests := []struct {
		name     string
		dateText string
		wantErr  bool
	}{
		{
			name:     "valid date",
			dateText: "2023-12-25T10:30:00Z",
			wantErr:  false,
		},
		{
			name:     "invalid format",
			dateText: "2023-12-25",
			wantErr:  true,
		},
		{
			name:     "empty date",
			dateText: "",
			wantErr:  true,
		},
		{
			name:     "invalid date",
			dateText: "invalid-date",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDate(tt.dateText)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateDividendCashAction tests dividend cash action validation
func TestValidateDividendCashAction(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{
			name:    "valid REINVEST",
			action:  "REINVEST",
			wantErr: false,
		},
		{
			name:    "valid TO_ACCOUNT_CASH",
			action:  "TO_ACCOUNT_CASH",
			wantErr: false,
		},
		{
			name:    "invalid action",
			action:  "INVALID",
			wantErr: true,
		},
		{
			name:    "empty action",
			action:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDividendCashAction(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDividendCashAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateIcon tests icon validation
func TestValidateIcon(t *testing.T) {
	tests := []struct {
		name    string
		icon    string
		wantErr bool
	}{
		{
			name:    "valid icon Home",
			icon:    "Home",
			wantErr: false,
		},
		{
			name:    "valid icon Tech",
			icon:    "Tech",
			wantErr: false,
		},
		{
			name:    "invalid icon",
			icon:    "InvalidIcon",
			wantErr: true,
		},
		{
			name:    "empty icon",
			icon:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIcon(tt.icon)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIcon() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateInstrumentShares tests instrument shares validation
func TestValidateInstrumentShares(t *testing.T) {
	tests := []struct {
		name             string
		instrumentShares map[string]float64
		wantErr          bool
	}{
		{
			name: "valid shares",
			instrumentShares: map[string]float64{
				"AAPL":  50.0,
				"GOOGL": 30.0,
			},
			wantErr: false,
		},
		{
			name:             "empty shares",
			instrumentShares: map[string]float64{},
			wantErr:          true,
		},
		{
			name: "zero shares",
			instrumentShares: map[string]float64{
				"AAPL": 0.0,
			},
			wantErr: true,
		},
		{
			name: "negative shares",
			instrumentShares: map[string]float64{
				"AAPL": -10.0,
			},
			wantErr: true,
		},
		{
			name: "empty ticker",
			instrumentShares: map[string]float64{
				"": 50.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInstrumentShares(tt.instrumentShares)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInstrumentShares() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestClientString tests the String method
func TestClientString(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		demo     bool
		expected string
	}{
		{
			name:     "demo client",
			apiKey:   "test-api-key-1234",
			demo:     true,
			expected: "Trading212(api_key=****1234, demo=true)",
		},
		{
			name:     "live client",
			apiKey:   "live-api-key-5678",
			demo:     false,
			expected: "Trading212(api_key=****5678, demo=false)",
		},
		{
			name:     "short api key",
			apiKey:   "abc",
			demo:     false,
			expected: "Trading212(api_key=****, demo=false)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.demo)
			result := client.String()
			if result != tt.expected {
				t.Errorf("Client.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestClientCash tests the Cash method with mock server
func TestClientCash(t *testing.T) {
	mockResponse := CashInfo{
		Free:              1000.50,
		Total:             5000.00,
		PieOrders:         200.00,
		Interest:          10.25,
		CashForInvestment: 800.25,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/equity/account/cash" {
			t.Errorf("Expected path /api/v0/equity/account/cash, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("Expected Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-api-key", true)
	client.host = server.URL

	cash, err := client.Cash()
	if err != nil {
		t.Fatalf("Cash() error = %v", err)
	}

	if cash.Free != mockResponse.Free {
		t.Errorf("Cash.Free = %v, want %v", cash.Free, mockResponse.Free)
	}
	if cash.Total != mockResponse.Total {
		t.Errorf("Cash.Total = %v, want %v", cash.Total, mockResponse.Total)
	}
}

// TestClientAccountInfo tests the AccountInfo method with mock server
func TestClientAccountInfo(t *testing.T) {
	mockResponse := AccountInfo{
		CurrencyCode: "GBP",
		ID:           12345,
		Type:         "LIVE",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/equity/account/info" {
			t.Errorf("Expected path /api/v0/equity/account/info, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-api-key", true)
	client.host = server.URL

	info, err := client.AccountInfo()
	if err != nil {
		t.Fatalf("AccountInfo() error = %v", err)
	}

	if info.CurrencyCode != mockResponse.CurrencyCode {
		t.Errorf("AccountInfo.CurrencyCode = %v, want %v", info.CurrencyCode, mockResponse.CurrencyCode)
	}
	if info.ID != mockResponse.ID {
		t.Errorf("AccountInfo.ID = %v, want %v", info.ID, mockResponse.ID)
	}
}

// TestClientPortfolio tests the Portfolio method with mock server
func TestClientPortfolio(t *testing.T) {
	mockResponse := []Position{
		{
			Ticker:   "AAPL",
			Quantity: 10.5,
			Value:    1500.75,
		},
		{
			Ticker:   "GOOGL",
			Quantity: 5.0,
			Value:    2500.00,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/equity/portfolio" {
			t.Errorf("Expected path /api/v0/equity/portfolio, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-api-key", true)
	client.host = server.URL

	portfolio, err := client.Portfolio()
	if err != nil {
		t.Fatalf("Portfolio() error = %v", err)
	}

	if len(portfolio) != len(mockResponse) {
		t.Errorf("Portfolio length = %v, want %v", len(portfolio), len(mockResponse))
	}

	if len(portfolio) > 0 {
		if portfolio[0].Ticker != mockResponse[0].Ticker {
			t.Errorf("Portfolio[0].Ticker = %v, want %v", portfolio[0].Ticker, mockResponse[0].Ticker)
		}
		if portfolio[0].Quantity != mockResponse[0].Quantity {
			t.Errorf("Portfolio[0].Quantity = %v, want %v", portfolio[0].Quantity, mockResponse[0].Quantity)
		}
	}
}

// TestClientErrorHandling tests error handling
func TestClientErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code":"InvalidRequest","message":"Bad request"}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key", true)
	client.host = server.URL

	_, err := client.Cash()
	if err == nil {
		t.Error("Expected error for 400 status code")
	}

	expectedError := "API error 400"
	if !contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %v, got %v", expectedError, err.Error())
	}
}

// TestExportCSV tests the ExportCSV method
func TestExportCSV(t *testing.T) {
	mockResponse := map[string]interface{}{
		"exportId": "12345",
		"status":   "PENDING",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/api/v0/history/exports" {
			t.Errorf("Expected path /api/v0/history/exports, got %s", r.URL.Path)
		}

		var req ExportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Error decoding request: %v", err)
		}

		if !req.DataIncluded.IncludeDividends {
			t.Error("Expected IncludeDividends to be true")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-api-key", true)
	client.host = server.URL

	timeFrom := time.Now().AddDate(0, -1, 0)
	timeTo := time.Now()

	result, err := client.ExportCSV(timeFrom, timeTo, true, true, true, true)
	if err != nil {
		t.Fatalf("ExportCSV() error = %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && contains(s[1:], substr))
}
