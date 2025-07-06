package main

import (
	"fmt"
	"log"
	"time"
	"github.com/0xnu/trading212"
)

// TradingDemoRunner handles all trading operations
type TradingDemoRunner struct {
	client *trading212.Client
}

// NewTradingDemoRunner creates a new demo runner instance
func NewTradingDemoRunner(apiKey string, isDemo bool) *TradingDemoRunner {
	return &TradingDemoRunner{
		client: trading212.NewClient(apiKey, isDemo),
	}
}

// Run executes the complete trading demo
func (t *TradingDemoRunner) Run() {
	fmt.Println("Starting Trading212 Demo...")

	t.fetchAndDisplayOrders()
	t.fetchAndDisplayCash()
	t.fetchAndDisplayPortfolio()
	t.fetchAndDisplayAccountInfo()
	t.fetchAndDisplayInstruments()
	t.fetchAndDisplayPies()
	t.fetchAndDisplayDividends()
	t.fetchAndDisplayTransactions()
	t.placeTestOrder()
	t.createTestPie()

	fmt.Println("\nDemo completed successfully!")
}

// fetchAndDisplayOrders retrieves and displays order information
func (t *TradingDemoRunner) fetchAndDisplayOrders() {
	fmt.Println("Fetching orders...")
	orders, err := t.client.Orders(0, "", 50)
	if err != nil {
		log.Printf("Error fetching orders: %v", err)
		return
	}
	fmt.Printf("Found %d orders\n", len(orders))
}

// fetchAndDisplayCash retrieves and displays cash information
func (t *TradingDemoRunner) fetchAndDisplayCash() {
	fmt.Println("\nFetching cash information...")
	cash, err := t.client.Cash()
	if err != nil {
		log.Printf("Error fetching cash: %v", err)
		return
	}
	fmt.Printf("Free cash: £%.2f\n", cash.Free)
	fmt.Printf("Total cash: £%.2f\n", cash.Total)
}

// fetchAndDisplayPortfolio retrieves and displays portfolio positions
func (t *TradingDemoRunner) fetchAndDisplayPortfolio() {
	fmt.Println("\nFetching portfolio...")
	portfolio, err := t.client.Portfolio()
	if err != nil {
		log.Printf("Error fetching portfolio: %v", err)
		return
	}

	fmt.Printf("Portfolio has %d positions\n", len(portfolio))
	for _, position := range portfolio {
		fmt.Printf("- %s: %.2f shares (£%.2f)\n", position.Ticker, position.Quantity, position.Value)
	}
}

// fetchAndDisplayAccountInfo retrieves and displays account information
func (t *TradingDemoRunner) fetchAndDisplayAccountInfo() {
	fmt.Println("\nFetching account info...")
	accountInfo, err := t.client.AccountInfo()
	if err != nil {
		log.Printf("Error fetching account info: %v", err)
		return
	}

	fmt.Printf("Account ID: %d\n", accountInfo.ID)
	fmt.Printf("Currency: %s\n", accountInfo.CurrencyCode)
	fmt.Printf("Account Type: %s\n", accountInfo.Type)
}

// fetchAndDisplayInstruments retrieves and displays available instruments
func (t *TradingDemoRunner) fetchAndDisplayInstruments() {
	fmt.Println("\nFetching instruments...")
	instruments, err := t.client.Instruments()
	if err != nil {
		fmt.Printf("Error fetching instruments: %v\n", err)
		return
	}

	fmt.Printf("Instruments data retrieved successfully\n")
	t.displaySampleInstruments(instruments)
}

// displaySampleInstruments shows sample instruments if available
func (t *TradingDemoRunner) displaySampleInstruments(instruments interface{}) {
	instrumentList, ok := instruments.([]interface{})
	if !ok || len(instrumentList) == 0 {
		return
	}

	fmt.Printf("Sample instruments available: %d total\n", len(instrumentList))
	for i, instrument := range instrumentList {
		if i >= 3 {
			break
		}
		fmt.Printf("- %+v\n", instrument)
	}
}

// fetchAndDisplayPies retrieves and displays pie information
func (t *TradingDemoRunner) fetchAndDisplayPies() {
	fmt.Println("\nFetching pies...")
	pies, err := t.client.Pies()
	if err != nil {
		log.Printf("Error fetching pies: %v", err)
		return
	}

	fmt.Printf("Found %d pies\n", len(pies))
	for _, pie := range pies {
		fmt.Printf("- %s (ID: %d, Goal: £%d)\n", pie.Name, pie.ID, pie.Goal)
	}
}

// fetchAndDisplayDividends retrieves and displays dividend information
func (t *TradingDemoRunner) fetchAndDisplayDividends() {
	fmt.Println("\nFetching dividends...")
	dividends, err := t.client.Dividends(0, "", 50)
	if err != nil {
		fmt.Printf("Error fetching dividends (normal for demo accounts): %v\n", err)
		return
	}
	fmt.Printf("Found %d dividend payments\n", len(dividends))
}

// fetchAndDisplayTransactions retrieves and displays transaction information
func (t *TradingDemoRunner) fetchAndDisplayTransactions() {
	fmt.Println("\nFetching transactions...")
	transactions, err := t.client.Transactions(0, 50)
	if err != nil {
		fmt.Printf("Error fetching transactions: %v\n", err)
		return
	}
	fmt.Printf("Found %d transactions\n", len(transactions))
}

// requestCSVExport requests a CSV export for the last month
func (t *TradingDemoRunner) requestCSVExport() {
	fmt.Println("\nRequesting CSV export...")
	timeFrom := time.Now().AddDate(0, -1, 0) // 1 month ago
	timeTo := time.Now()

	exportResult, err := t.client.ExportCSV(timeFrom, timeTo, true, true, true, true)
	if err != nil {
		fmt.Printf("Error requesting export (normal for demo accounts): %v\n", err)
		return
	}
	fmt.Printf("Export requested successfully: %v\n", exportResult)
}

// placeTestOrder demonstrates placing a limit order
func (t *TradingDemoRunner) placeTestOrder() {
	fmt.Println("\nPlacing limit order...")
	order, err := t.client.EquityOrderPlaceLimit("AAPL", 1, 150.00, "GTC")
	if err != nil {
		log.Printf("Error placing order: %v", err)
		return
	}
	fmt.Printf("Order placed successfully: ID %d\n", order.ID)
}

// createTestPie demonstrates creating a new pie
func (t *TradingDemoRunner) createTestPie() {
	fmt.Println("\nCreating a demo pie...")
	instrumentShares := map[string]float64{
		"PLTR_US_EQ": 0.40, // Palantir Technologies - 40%
		"AAPL_US_EQ": 0.30, // Apple Inc - 30%
		"MSFT_US_EQ": 0.30, // Microsoft Corporation - 30%
	}

	endDate := time.Now().AddDate(1, 0, 0) // 1 year from now
	pie, err := t.client.PieCreate("REINVEST", endDate, 10000, "Tech", "Big Tech Portfolio", instrumentShares)
	if err != nil {
		log.Printf("Error creating pie: %v", err)
		return
	}
	fmt.Printf("Pie created successfully: %s (ID: %d)\n", pie.Name, pie.ID)
}

func main() {
	demoRunner := NewTradingDemoRunner("your_api_key", true)
	demoRunner.Run()
}