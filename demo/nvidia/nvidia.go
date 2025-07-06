package main

import (
	"github.com/0xnu/trading212"
	"log"
	"math"
	"strings"
	"time"
)

type TradingBot struct {
	client      *trading212.Client
	ticker      string
	riskPercent float64
}

type RetryConfig struct {
	maxRetries int
	baseDelay  time.Duration
}

func NewTradingBot(apiKey string, isDemo bool, ticker string, riskPercent float64) *TradingBot {
	return &TradingBot{
		client:      trading212.NewClient(apiKey, isDemo),
		ticker:      ticker,
		riskPercent: riskPercent,
	}
}

func main() {
	log.Println("Starting trading bot...")

	bot := NewTradingBot("your_api_key", true, "NVDA", 1.0)
	bot.Run()
}

func (bot *TradingBot) Run() {
	for {
		log.Println("Starting new iteration...")

		if bot.executeTradeLogic() {
			log.Println("Trade logic executed successfully")
		} else {
			log.Println("Trade logic failed, waiting before retry")
			time.Sleep(5 * time.Minute)
			continue
		}

		time.Sleep(5 * time.Minute)
	}
}

func (bot *TradingBot) executeTradeLogic() bool {
	positions := bot.getPositionsWithRetry()
	if positions == nil {
		return false
	}

	currentPrice := bot.getCurrentPrice()
	if currentPrice == 0 {
		log.Println("Failed to get current price")
		return false
	}

	upperBand, lowerBand := bot.calculateBollingerBands(currentPrice)
	portfolioValue := bot.calculatePortfolioValue(positions)
	positionSize := bot.calculatePositionSize(portfolioValue, currentPrice)
	hasPosition := bot.hasExistingPosition(positions)

	return bot.executeTrade(hasPosition, currentPrice, upperBand, lowerBand, positionSize)
}

func (bot *TradingBot) getPositionsWithRetry() []trading212.Position {
	config := RetryConfig{maxRetries: 5, baseDelay: time.Second}

	for i := 0; i < config.maxRetries; i++ {
		positions, err := bot.client.Portfolio()
		if err == nil {
			return positions
		}

		if !bot.shouldRetry(err, i, config) {
			break
		}
	}

	return nil
}

func (bot *TradingBot) shouldRetry(err error, attempt int, config RetryConfig) bool {
	if !strings.Contains(err.Error(), "TooManyRequests") {
		log.Printf("API error: %v", err)
		return false
	}

	if attempt >= config.maxRetries-1 {
		log.Printf("Failed after retries: %v", err)
		return false
	}

	delay := time.Duration(math.Pow(2, float64(attempt))) * config.baseDelay
	log.Printf("Rate limited, retry %d/%d in %v", attempt+1, config.maxRetries, delay)
	time.Sleep(delay)
	return true
}

func (bot *TradingBot) getCurrentPrice() float64 {
	positions, err := bot.client.Portfolio()
	if err != nil {
		log.Printf("Portfolio error: %v", err)
		return 0
	}

	log.Printf("Got %d positions", len(positions))

	for _, pos := range positions {
		if pos.Ticker == bot.ticker {
			return pos.Value / pos.Quantity
		}
	}

	return 0
}

func (bot *TradingBot) calculateBollingerBands(currentPrice float64) (float64, float64) {
	stdDev := currentPrice * 0.02 // 2% of price as approximation
	upperBand := currentPrice + (2 * stdDev)
	lowerBand := currentPrice - (2 * stdDev)
	return upperBand, lowerBand
}

func (bot *TradingBot) calculatePortfolioValue(positions []trading212.Position) float64 {
	portfolioValue := 0.0
	for _, pos := range positions {
		portfolioValue += pos.Value
	}
	return portfolioValue
}

func (bot *TradingBot) calculatePositionSize(portfolioValue, currentPrice float64) int {
	positionSize := int((portfolioValue * bot.riskPercent / 100) / currentPrice)
	if positionSize < 1 {
		positionSize = 1
	}
	return positionSize
}

func (bot *TradingBot) hasExistingPosition(positions []trading212.Position) bool {
	for _, pos := range positions {
		if pos.Ticker == bot.ticker {
			return true
		}
	}
	return false
}

func (bot *TradingBot) executeTrade(hasPosition bool, currentPrice, upperBand, lowerBand float64, positionSize int) bool {
	if bot.shouldBuy(hasPosition, currentPrice, lowerBand) {
		return bot.placeBuyOrder(positionSize, currentPrice)
	}

	if bot.shouldSell(hasPosition, currentPrice, upperBand) {
		return bot.placeSellOrder(positionSize, currentPrice)
	}

	return true
}

func (bot *TradingBot) shouldBuy(hasPosition bool, currentPrice, lowerBand float64) bool {
	return !hasPosition && currentPrice < lowerBand
}

func (bot *TradingBot) shouldSell(hasPosition bool, currentPrice, upperBand float64) bool {
	return hasPosition && currentPrice > upperBand
}

func (bot *TradingBot) placeBuyOrder(positionSize int, currentPrice float64) bool {
	_, err := bot.client.EquityOrderPlaceMarket(bot.ticker, positionSize)
	if err != nil {
		log.Printf("Buy error: %v", err)
		return false
	}

	log.Printf("Bought %d %s @ %.2f (Bollinger Band strategy)", positionSize, bot.ticker, currentPrice)
	return true
}

func (bot *TradingBot) placeSellOrder(positionSize int, currentPrice float64) bool {
	_, err := bot.client.EquityOrderPlaceMarket(bot.ticker, -positionSize)
	if err != nil {
		log.Printf("Sell error: %v", err)
		return false
	}

	log.Printf("Sold %d %s @ %.2f (Bollinger Band strategy)", positionSize, bot.ticker, currentPrice)
	return true
}
