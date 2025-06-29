package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client represents the Trading212 REST API client
type Client struct {
	apiKey     string
	host       string
	httpClient *http.Client
}

// Order represents an order structure
type Order struct {
	ID           int     `json:"id"`
	Ticker       string  `json:"ticker"`
	Quantity     int     `json:"quantity"`
	LimitPrice   float64 `json:"limitPrice,omitempty"`
	StopPrice    float64 `json:"stopPrice,omitempty"`
	TimeValidity string  `json:"timeValidity"`
}

// Position represents a portfolio position
type Position struct {
	Ticker   string  `json:"ticker"`
	Quantity float64 `json:"quantity"`
	Value    float64 `json:"value"`
}

// Pie represents a Trading212 pie
type Pie struct {
	ID                 int                `json:"id"`
	Name               string             `json:"name"`
	Icon               string             `json:"icon"`
	Goal               int                `json:"goal"`
	EndDate            string             `json:"endDate"`
	DividendCashAction string             `json:"dividendCashAction"`
	InstrumentShares   map[string]float64 `json:"instrumentShares"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Items        []interface{} `json:"items"`
	NextPagePath string        `json:"nextPagePath,omitempty"`
}

// CashInfo represents account cash information
type CashInfo struct {
	Free              float64 `json:"free"`
	Total             float64 `json:"total"`
	PieOrders         float64 `json:"pieOrders"`
	Interest          float64 `json:"interest"`
	CashForInvestment float64 `json:"cashForInvestment"`
}

// AccountInfo represents account information
type AccountInfo struct {
	CurrencyCode string `json:"currencyCode"`
	ID           int    `json:"id"`
	Type         string `json:"type"`
}

// ExportRequest represents a CSV export request
type ExportRequest struct {
	DataIncluded DataIncluded `json:"dataIncluded"`
	TimeFrom     string       `json:"timeFrom"`
	TimeTo       string       `json:"timeTo"`
}

// DataIncluded specifies what data to include in exports
type DataIncluded struct {
	IncludeDividends    bool `json:"includeDividends"`
	IncludeInterest     bool `json:"includeInterest"`
	IncludeOrders       bool `json:"includeOrders"`
	IncludeTransactions bool `json:"includeTransactions"`
}

// NewClient creates a new Trading212 client
func NewClient(apiKey string, demo bool) *Client {
	host := "https://live.trading212.com"
	if demo {
		host = "https://demo.trading212.com"
	}

	return &Client{
		apiKey:     apiKey,
		host:       host,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// get performs a GET request to the API
func (c *Client) get(endpoint string, params url.Values, apiVersion string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s/%s", c.host, apiVersion, endpoint)
	if params != nil {
		url += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)
	return c.processRequest(req)
}

// post performs a POST request to the API
func (c *Client) post(endpoint string, data interface{}, apiVersion string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s/%s", c.host, apiVersion, endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return c.processRequest(req)
}

// getURL performs a GET request to a full URL path
func (c *Client) getURL(urlPath string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.host, urlPath)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)
	return c.processRequest(req)
}

// deleteURL performs a DELETE request to a full URL path
func (c *Client) deleteURL(urlPath string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.host, urlPath)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)
	return c.processRequest(req)
}

// processRequest executes HTTP request and handles response
func (c *Client) processRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// processItems handles paginated responses
func (c *Client) processItems(initialResponse []byte) ([]interface{}, error) {
	var response PaginatedResponse
	if err := json.Unmarshal(initialResponse, &response); err != nil {
		return nil, err
	}

	items := response.Items

	for response.NextPagePath != "" {
		nextData, err := c.getURL(response.NextPagePath)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(nextData, &response); err != nil {
			return nil, err
		}

		items = append(items, response.Items...)
	}

	return items, nil
}

// validateTimeValidity validates time validity parameter
func validateTimeValidity(timeValidity string) error {
	if timeValidity != "GTC" && timeValidity != "DAY" {
		return fmt.Errorf("time_validity must be one of GTC or DAY")
	}
	return nil
}

// validateDate validates date format
func validateDate(dateText string) error {
	_, err := time.Parse("2006-01-02T15:04:05Z", dateText)
	if err != nil {
		return fmt.Errorf("incorrect date format, should be YYYY-MM-DDTHH:MM:SSZ")
	}
	return nil
}

// validateDividendCashAction validates dividend cash action
func validateDividendCashAction(action string) error {
	validActions := []string{"REINVEST", "TO_ACCOUNT_CASH"}
	for _, valid := range validActions {
		if action == valid {
			return nil
		}
	}
	return fmt.Errorf("dividendCashAction must be one of %v", validActions)
}

// validateIcon validates pie icon
func validateIcon(icon string) error {
	validIcons := []string{
		"Home", "PiggyBank", "Iceberg", "Airplane", "RV", "Unicorn", "Whale", "Convertable", "Family",
		"Coins", "Education", "BillsAndCoins", "Bills", "Water", "Wind", "Car", "Briefcase", "Medical",
		"Landscape", "Child", "Vault", "Travel", "Cabin", "Apartments", "Burger", "Bus", "Energy",
		"Factory", "Global", "Leaf", "Materials", "Pill", "Ring", "Shipping", "Storefront", "Tech", "Umbrella",
	}

	for _, valid := range validIcons {
		if icon == valid {
			return nil
		}
	}
	return fmt.Errorf("icon must be one of %v", validIcons)
}

// validateInstrumentShares validates instrument shares map
func validateInstrumentShares(instrumentShares map[string]float64) error {
	if len(instrumentShares) == 0 {
		return fmt.Errorf("instrument_shares cannot be empty")
	}

	for ticker, shares := range instrumentShares {
		if ticker == "" {
			return fmt.Errorf("instrument identifiers must be non-empty strings")
		}
		if shares <= 0 {
			return fmt.Errorf("number of shares must be greater than zero")
		}
	}
	return nil
}

// Orders fetches historical order data
func (c *Client) Orders(cursor int, ticker string, limit int) ([]interface{}, error) {
	params := url.Values{}
	params.Set("cursor", strconv.Itoa(cursor))
	params.Set("limit", strconv.Itoa(limit))
	if ticker != "" {
		params.Set("ticker", ticker)
	}

	response, err := c.get("equity/history/orders", params, "v0")
	if err != nil {
		return nil, err
	}

	return c.processItems(response)
}

// Dividends fetches dividends paid out
func (c *Client) Dividends(cursor int, ticker string, limit int) ([]interface{}, error) {
	params := url.Values{}
	params.Set("cursor", strconv.Itoa(cursor))
	params.Set("limit", strconv.Itoa(limit))
	if ticker != "" {
		params.Set("ticker", ticker)
	}

	response, err := c.get("history/dividends", params, "v0")
	if err != nil {
		return nil, err
	}

	return c.processItems(response)
}

// Export fetches all account exports as a list
func (c *Client) Export() ([]interface{}, error) {
	response, err := c.get("history/exports", nil, "v0")
	if err != nil {
		return nil, err
	}

	var exports []interface{}
	if err := json.Unmarshal(response, &exports); err != nil {
		return nil, err
	}

	return exports, nil
}

// ExportCSV requests a CSV export of account history
func (c *Client) ExportCSV(timeFrom, timeTo time.Time, includeDividends, includeInterest, includeOrders, includeTransactions bool) (interface{}, error) {
	exportReq := ExportRequest{
		DataIncluded: DataIncluded{
			IncludeDividends:    includeDividends,
			IncludeInterest:     includeInterest,
			IncludeOrders:       includeOrders,
			IncludeTransactions: includeTransactions,
		},
		TimeFrom: timeFrom.Format("2006-01-02T15:04:05Z"),
		TimeTo:   timeTo.Format("2006-01-02T15:04:05Z"),
	}

	response, err := c.post("history/exports", exportReq, "v0")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Transactions fetches transactions list
func (c *Client) Transactions(cursor, limit int) ([]interface{}, error) {
	params := url.Values{}
	if cursor > 0 {
		params.Set("cursor", strconv.Itoa(cursor))
	}
	params.Set("limit", strconv.Itoa(limit))

	response, err := c.get("history/transactions", params, "v0")
	if err != nil {
		return nil, err
	}

	return c.processItems(response)
}

// Instruments fetches tradeable instruments metadata
func (c *Client) Instruments() (interface{}, error) {
	response, err := c.get("equity/metadata/instruments", nil, "v0")
	if err != nil {
		return nil, err
	}

	var instruments interface{}
	if err := json.Unmarshal(response, &instruments); err != nil {
		return nil, err
	}

	return instruments, nil
}

// Cash fetches account cash information
func (c *Client) Cash() (*CashInfo, error) {
	response, err := c.get("equity/account/cash", nil, "v0")
	if err != nil {
		return nil, err
	}

	var cash CashInfo
	if err := json.Unmarshal(response, &cash); err != nil {
		return nil, err
	}

	return &cash, nil
}

// Portfolio fetches all open positions
func (c *Client) Portfolio() ([]Position, error) {
	response, err := c.get("equity/portfolio", nil, "v0")
	if err != nil {
		return nil, err
	}

	var portfolio []Position
	if err := json.Unmarshal(response, &portfolio); err != nil {
		return nil, err
	}

	return portfolio, nil
}

// Position fetches open position by ticker
func (c *Client) Position(ticker string) (*Position, error) {
	response, err := c.get(fmt.Sprintf("equity/portfolio/%s", ticker), nil, "v0")
	if err != nil {
		return nil, err
	}

	var position Position
	if err := json.Unmarshal(response, &position); err != nil {
		return nil, err
	}

	return &position, nil
}

// Exchanges fetches exchange list
func (c *Client) Exchanges() (interface{}, error) {
	response, err := c.get("equity/metadata/exchanges", nil, "v0")
	if err != nil {
		return nil, err
	}

	var exchanges interface{}
	if err := json.Unmarshal(response, &exchanges); err != nil {
		return nil, err
	}

	return exchanges, nil
}

// AccountInfo fetches account information
func (c *Client) AccountInfo() (*AccountInfo, error) {
	response, err := c.get("equity/account/info", nil, "v0")
	if err != nil {
		return nil, err
	}

	var info AccountInfo
	if err := json.Unmarshal(response, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// EquityOrders fetches all equity orders
func (c *Client) EquityOrders() ([]Order, error) {
	response, err := c.get("equity/orders", nil, "v0")
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(response, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// EquityOrder fetches equity order by ID
func (c *Client) EquityOrder(id int) (*Order, error) {
	response, err := c.get(fmt.Sprintf("equity/orders/%d", id), nil, "v0")
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// EquityOrderCancel cancels equity order by ID
func (c *Client) EquityOrderCancel(id int) error {
	_, err := c.deleteURL(fmt.Sprintf("/equity/orders/%d", id))
	return err
}

// EquityOrderPlaceLimit places a limit order
func (c *Client) EquityOrderPlaceLimit(ticker string, quantity int, limitPrice float64, timeValidity string) (*Order, error) {
	if err := validateTimeValidity(timeValidity); err != nil {
		return nil, err
	}

	orderData := map[string]interface{}{
		"quantity":     quantity,
		"limitPrice":   limitPrice,
		"ticker":       ticker,
		"timeValidity": timeValidity,
	}

	response, err := c.post("equity/orders/limit", orderData, "v0")
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// EquityOrderPlaceMarket places a market order
func (c *Client) EquityOrderPlaceMarket(ticker string, quantity int) (*Order, error) {
	orderData := map[string]interface{}{
		"quantity": quantity,
		"ticker":   ticker,
	}

	response, err := c.post("equity/orders/market", orderData, "v0")
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// EquityOrderPlaceStop places a stop order
func (c *Client) EquityOrderPlaceStop(ticker string, quantity int, stopPrice float64, timeValidity string) (*Order, error) {
	if err := validateTimeValidity(timeValidity); err != nil {
		return nil, err
	}

	orderData := map[string]interface{}{
		"quantity":     quantity,
		"stopPrice":    stopPrice,
		"ticker":       ticker,
		"timeValidity": timeValidity,
	}

	response, err := c.post("equity/orders/stop", orderData, "v0")
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// EquityOrderPlaceStopLimit places a stop-limit order
func (c *Client) EquityOrderPlaceStopLimit(ticker string, quantity int, stopPrice, limitPrice float64, timeValidity string) (*Order, error) {
	if err := validateTimeValidity(timeValidity); err != nil {
		return nil, err
	}

	orderData := map[string]interface{}{
		"quantity":     quantity,
		"stopPrice":    stopPrice,
		"limitPrice":   limitPrice,
		"ticker":       ticker,
		"timeValidity": timeValidity,
	}

	response, err := c.post("equity/orders/stop_limit", orderData, "v0")
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// Pies fetches all pies
func (c *Client) Pies() ([]Pie, error) {
	response, err := c.get("equity/pies", nil, "v0")
	if err != nil {
		return nil, err
	}

	var pies []Pie
	if err := json.Unmarshal(response, &pies); err != nil {
		return nil, err
	}

	return pies, nil
}

// PieCreate creates a new pie
func (c *Client) PieCreate(dividendCashAction string, endDate time.Time, goal int, icon, name string, instrumentShares map[string]float64) (*Pie, error) {
	if err := validateDividendCashAction(dividendCashAction); err != nil {
		return nil, err
	}
	if err := validateIcon(icon); err != nil {
		return nil, err
	}
	if err := validateInstrumentShares(instrumentShares); err != nil {
		return nil, err
	}

	endDateStr := endDate.Format("2006-01-02T15:04:05Z")
	if err := validateDate(endDateStr); err != nil {
		return nil, err
	}

	pieData := map[string]interface{}{
		"dividendCashAction": dividendCashAction,
		"endDate":            endDateStr,
		"goal":               goal,
		"icon":               icon,
		"instrumentShares":   instrumentShares,
		"name":               name,
	}

	response, err := c.post("equity/pies", pieData, "v0")
	if err != nil {
		return nil, err
	}

	var pie Pie
	if err := json.Unmarshal(response, &pie); err != nil {
		return nil, err
	}

	return &pie, nil
}

// PieDelete deletes pie by ID
func (c *Client) PieDelete(id int) error {
	_, err := c.deleteURL(fmt.Sprintf("/equity/pies/%d", id))
	return err
}

// Pie fetches pie by ID
func (c *Client) Pie(id int) (*Pie, error) {
	response, err := c.get(fmt.Sprintf("equity/pies/%d", id), nil, "v0")
	if err != nil {
		return nil, err
	}

	var pie Pie
	if err := json.Unmarshal(response, &pie); err != nil {
		return nil, err
	}

	return &pie, nil
}

// PieUpdate updates existing pie
func (c *Client) PieUpdate(id int, dividendCashAction, endDate string, goal int, icon, name string, instrumentShares map[string]float64) (*Pie, error) {
	if err := validateDividendCashAction(dividendCashAction); err != nil {
		return nil, err
	}
	if err := validateIcon(icon); err != nil {
		return nil, err
	}
	if err := validateDate(endDate); err != nil {
		return nil, err
	}
	if err := validateInstrumentShares(instrumentShares); err != nil {
		return nil, err
	}

	pieData := map[string]interface{}{
		"dividendCashAction": dividendCashAction,
		"endDate":            endDate,
		"goal":               goal,
		"icon":               icon,
		"instrumentShares":   instrumentShares,
		"name":               name,
	}

	response, err := c.post(fmt.Sprintf("equity/pies/%d", id), pieData, "v0")
	if err != nil {
		return nil, err
	}

	var pie Pie
	if err := json.Unmarshal(response, &pie); err != nil {
		return nil, err
	}

	return &pie, nil
}

// String returns string representation of the client
func (c *Client) String() string {
	demo := strings.Contains(c.host, "demo")
	apiKeySuffix := ""
	if len(c.apiKey) >= 4 {
		apiKeySuffix = c.apiKey[len(c.apiKey)-4:]
	}
	return fmt.Sprintf("Trading212(api_key=****%s, demo=%t)", apiKeySuffix, demo)
}
