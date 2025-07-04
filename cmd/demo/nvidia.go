package main

import (
	"log"
	"time"
	"math"
	"strings"
	"github.com/0xnu/trading212"
)

func main() {
    log.Println("Starting trading bot...")
	client := trading212.NewClient("your_api_key", true) // true for demo; false for live
	
	ticker := "NVDA"
	riskPercent := 1.0
	
	for {
		log.Println("Starting new iteration...")
		
		// Get portfolio with retry logic
		var positions []trading212.Position
		var err error
		maxRetries := 5
		baseDelay := time.Second
		
		for i := 0; i < maxRetries; i++ {
			positions, err = client.Portfolio()
			if err == nil {
				break
			}
			
			if strings.Contains(err.Error(), "TooManyRequests") {
				delay := time.Duration(math.Pow(2, float64(i))) * baseDelay
				log.Printf("Rate limited, retry %d/%d in %v", i+1, maxRetries, delay)
				time.Sleep(delay)
				continue
			}
			
			log.Printf("API error: %v", err)
			break
		}
		
		if err != nil {
			log.Printf("Failed after retries: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		
		currentPrice := getCurrentPrice(client, ticker)
		upperBand, _, lowerBand := calculateMovingAverage(client, ticker, 20)
		
		// Calculate total portfolio value for position sizing
		portfolioValue := 0.0
		for _, pos := range positions {
			portfolioValue += pos.Value
		}
		
		positionSize := int((portfolioValue * riskPercent/100) / currentPrice)
		if positionSize < 1 {
			positionSize = 1
		}
		
		hasPosition := false
		for _, pos := range positions {
			if pos.Ticker == ticker {
				hasPosition = true
				break
			}
		}
		
		// Trading logic
		if !hasPosition && currentPrice < lowerBand {
			_, err := client.EquityOrderPlaceMarket(ticker, positionSize)
			if err != nil {
				log.Printf("Buy error: %v", err)
			} else {
				log.Printf("Bought %d %s @ %.2f (Bollinger Band strategy)", 
					positionSize, ticker, currentPrice)
			}
		} else if hasPosition && currentPrice > upperBand {
			_, err := client.EquityOrderPlaceMarket(ticker, -positionSize)
			if err != nil {
				log.Printf("Sell error: %v", err)
			} else {
				log.Printf("Sold %d %s @ %.2f (Bollinger Band strategy)", 
					positionSize, ticker, currentPrice)
			}
		}
		
		time.Sleep(5 * time.Minute)
	}
}

func getCurrentPrice(client *trading212.Client, ticker string) float64 {
    positions, err := client.Portfolio()
    if err != nil {
        log.Printf("Portfolio error: %v", err)
    } else {
        log.Printf("Got %d positions", len(positions))
    }
    
    for _, pos := range positions {
        if pos.Ticker == ticker {
            return pos.Value / pos.Quantity
        }
    }
    return 0
}

func calculateMovingAverage(client *trading212.Client, ticker string, period int) (float64, float64, float64) {
    // Calculate current price as middle band
    currentPrice := getCurrentPrice(client, ticker)
    
    // Standard deviation (simplified approach)
    stdDev := currentPrice * 0.02 // 2% of price as approximation
    
    // Upper and lower bands (2 standard deviations)
    upperBand := currentPrice + (2 * stdDev)
    lowerBand := currentPrice - (2 * stdDev)
    
    return upperBand, currentPrice, lowerBand
}