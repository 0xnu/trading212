package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Initialize Trading212 client
	client := NewClient("your_api_key", true) // true for demo

	// Fetch historical orders
	fmt.Println("Fetching orders...")
	orders, err := client.Orders(0, "", 50)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d orders\n", len(orders))

	// Fetch account cash information
	fmt.Println("\nFetching cash information...")
	cash, err := client.Cash()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Free cash: £%.2f\n", cash.Free)
	fmt.Printf("Total cash: £%.2f\n", cash.Total)

	// Fetch portfolio positions
	fmt.Println("\nFetching portfolio...")
	portfolio, err := client.Portfolio()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Portfolio has %d positions\n", len(portfolio))
	for _, position := range portfolio {
		fmt.Printf("- %s: %.2f shares (£%.2f)\n", position.Ticker, position.Quantity, position.Value)
	}

	// Fetch account information
	fmt.Println("\nFetching account info...")
	accountInfo, err := client.AccountInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Account ID: %d\n", accountInfo.ID)
	fmt.Printf("Currency: %s\n", accountInfo.CurrencyCode)
	fmt.Printf("Account Type: %s\n", accountInfo.Type)

	// Fetch available instruments
	fmt.Println("\nFetching instruments...")
	instruments, err := client.Instruments()
	if err != nil {
		fmt.Printf("Error fetching instruments: %v\n", err)
	} else {
		fmt.Printf("Instruments data retrieved successfully\n")
		// Try to show some sample tickers if possible
		if instrumentList, ok := instruments.([]interface{}); ok && len(instrumentList) > 0 {
			fmt.Printf("Sample instruments available: %d total\n", len(instrumentList))
			// Show first few instruments
			for i, instrument := range instrumentList {
				if i >= 3 {
					break
				}
				fmt.Printf("- %+v\n", instrument)
			}
		}
	}

	// Fetch all pies
	fmt.Println("\nFetching pies...")
	pies, err := client.Pies()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d pies\n", len(pies))
	for _, pie := range pies {
		fmt.Printf("- %s (ID: %d, Goal: £%d)\n", pie.Name, pie.ID, pie.Goal)
	}

	// Fetch dividends
	fmt.Println("\nFetching dividends...")
	dividends, err := client.Dividends(0, "", 50)
	if err != nil {
		fmt.Printf("Error fetching dividends (this is normal for demo accounts): %v\n", err)
	} else {
		fmt.Printf("Found %d dividend payments\n", len(dividends))
	}

	// Fetch transactions
	fmt.Println("\nFetching transactions...")
	transactions, err := client.Transactions(0, 50)
	if err != nil {
		fmt.Printf("Error fetching transactions: %v\n", err)
	} else {
		fmt.Printf("Found %d transactions\n", len(transactions))
	}

	// Place a limit order (commented out for safety)
	/*
		fmt.Println("\nPlacing limit order...")
		order, err := client.EquityOrderPlaceLimit("AAPL", 1, 150.00, "GTC")
		if err != nil {
			log.Printf("Error placing order: %v", err)
		} else {
			fmt.Printf("Order placed successfully: ID %d\n", order.ID)
		}
	*/

	// Create a CSV export request
	fmt.Println("\nRequesting CSV export...")
	timeFrom := time.Now().AddDate(0, -1, 0) // 1 month ago
	timeTo := time.Now()

	exportResult, err := client.ExportCSV(timeFrom, timeTo, true, true, true, true)
	if err != nil {
		fmt.Printf("Error requesting export (this is normal for demo accounts): %v\n", err)
	} else {
		fmt.Printf("Export requested successfully: %v\n", exportResult)
	}

	// Create a pie (commented out for safety)
	/*
		// First, let's find valid instruments for the demo
		fmt.Println("\nFetching available instruments to find valid tickers...")
		instruments, err := client.Instruments()
		if err == nil {
			// Print first few instruments to see format
			fmt.Printf("Sample instruments: %+v\n", instruments)
		}

		fmt.Println("\nCreating a demo pie...")
		// Use common UK/European tickers that should be available in demo
		instrumentShares := map[string]float64{
			"LLOY": 40.0,  // Lloyds Banking Group
			"VOD":  30.0,  // Vodafone
			"BP":   30.0,  // BP plc
		}

		endDate := time.Now().AddDate(1, 0, 0) // 1 year from now
		pie, err := client.PieCreate("REINVEST", endDate, 10000, "Tech", "UK Giants", instrumentShares)
		if err != nil {
			log.Printf("Error creating pie: %v", err)
		} else {
			fmt.Printf("Pie created successfully: %s (ID: %d)\n", pie.Name, pie.ID)
		}
	*/

	fmt.Println("\nDemo completed successfully!")
}
