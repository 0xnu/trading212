package main

import (
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/0xnu/trading212"
)

type PriceData struct {
	Price float64
	High  float64
	Low   float64
	Close float64
}

type TradingIndicators struct {
	ATR        float64
	UpperBand  float64
	MiddleBand float64
	LowerBand  float64
	StdDev     float64
}

type StockData struct {
	Ticker       string
	PriceHistory []PriceData
	Indicators   TradingIndicators
	CurrentPrice float64
	LastUpdate   time.Time
}

type TradingBot struct {
	client          *trading212.Client
	stocks          map[string]*StockData
	riskPercent     float64
	atrPeriod       int
	bollingerPeriod int
	mutex           sync.RWMutex
}

func main() {
	log.Println("Starting Multi-Stock ATR + Bollinger Bands Trading Bot...")

	client := trading212.NewClient("your_api_key", true) // true for demo; false for live

	// Define stock universe
	tickers := []string{"NVDA", "PLTR", "TSLA", "AAPL", "GOOGL"}

	bot := &TradingBot{
		client:          client,
		stocks:          make(map[string]*StockData),
		riskPercent:     0.4, // 0.4% per stock (2% total across 5 stocks)
		atrPeriod:       14,
		bollingerPeriod: 20,
	}

	// Initialize stock data
	for _, ticker := range tickers {
		bot.stocks[ticker] = &StockData{
			Ticker:       ticker,
			PriceHistory: make([]PriceData, 0),
			LastUpdate:   time.Now(),
		}
	}

	for {
		log.Println("=== Starting new trading cycle ===")

		// Get portfolio data
		positions, portfolioValue, err := bot.getPortfolioData()
		if err != nil {
			log.Printf("Portfolio error: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		log.Printf("Portfolio Value: £%.2f", portfolioValue)

		// Update price data for all stocks
		bot.updateAllPrices()

		// Calculate indicators for all stocks
		bot.calculateAllIndicators()

		// Execute trading logic for each stock
		for ticker, stockData := range bot.stocks {
			bot.executeTrading(ticker, stockData, positions, portfolioValue)
		}

		// Display portfolio summary
		bot.displayPortfolioSummary(positions)

		time.Sleep(5 * time.Minute)
	}
}

func (bot *TradingBot) getPortfolioData() ([]trading212.Position, float64, error) {
	var positions []trading212.Position
	var err error
	maxRetries := 5
	baseDelay := time.Second

	for i := 0; i < maxRetries; i++ {
		positions, err = bot.client.Portfolio()
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
		return nil, 0, err
	}

	// Calculate total portfolio value
	portfolioValue := 0.0
	for _, pos := range positions {
		portfolioValue += pos.Value
	}

	return positions, portfolioValue, nil
}

func (bot *TradingBot) updateAllPrices() {
	bot.mutex.Lock()
	defer bot.mutex.Unlock()

	for ticker, stockData := range bot.stocks {
		currentPrice := bot.getCurrentPrice(ticker)
		if currentPrice == 0 {
			log.Printf("Could not get price for %s", ticker)
			continue
		}

		stockData.CurrentPrice = currentPrice
		stockData.LastUpdate = time.Now()

		// Add to price history (simplified OHLC approximation)
		priceData := PriceData{
			Price: currentPrice,
			High:  currentPrice * 1.005,
			Low:   currentPrice * 0.995,
			Close: currentPrice,
		}

		stockData.PriceHistory = append(stockData.PriceHistory, priceData)

		// Keep only required history
		maxHistory := int(math.Max(float64(bot.atrPeriod), float64(bot.bollingerPeriod)))
		if len(stockData.PriceHistory) > maxHistory {
			stockData.PriceHistory = stockData.PriceHistory[len(stockData.PriceHistory)-maxHistory:]
		}
	}
}

func (bot *TradingBot) calculateAllIndicators() {
	bot.mutex.Lock()
	defer bot.mutex.Unlock()

	for _, stockData := range bot.stocks {
		maxHistory := int(math.Max(float64(bot.atrPeriod), float64(bot.bollingerPeriod)))
		if len(stockData.PriceHistory) >= maxHistory {
			stockData.Indicators = bot.calculateIndicators(stockData.PriceHistory)
		}
	}
}

func (bot *TradingBot) executeTrading(ticker string, stockData *StockData, positions []trading212.Position, portfolioValue float64) {
	bot.mutex.RLock()
	indicators := stockData.Indicators
	currentPrice := stockData.CurrentPrice
	hasData := len(stockData.PriceHistory) >= int(math.Max(float64(bot.atrPeriod), float64(bot.bollingerPeriod)))
	bot.mutex.RUnlock()

	if !hasData || currentPrice == 0 {
		return
	}

	// Calculate position size for this stock
	positionSize := int((portfolioValue * bot.riskPercent / 100) / currentPrice)
	if positionSize < 1 {
		positionSize = 1
	}

	// Check existing position
	hasPosition := false
	var existingPosition trading212.Position
	for _, pos := range positions {
		if pos.Ticker == ticker {
			hasPosition = true
			existingPosition = pos
			break
		}
	}

	// Trading logic
	if !hasPosition {
		bot.evaluateEntry(ticker, currentPrice, indicators, positionSize)
	} else {
		bot.evaluateExit(ticker, currentPrice, indicators, existingPosition)
	}
}

func (bot *TradingBot) evaluateEntry(ticker string, currentPrice float64, indicators TradingIndicators, positionSize int) {
	// Entry conditions
	touchesLowerBand := currentPrice <= indicators.LowerBand
	lowVolatility := indicators.ATR < (currentPrice * 0.025) // ATR less than 2.5% of price
	belowMiddleBand := currentPrice < indicators.MiddleBand

	// Additional filters
	bollingerWidth := (indicators.UpperBand - indicators.LowerBand) / indicators.MiddleBand
	normalVolatility := bollingerWidth > 0.02 && bollingerWidth < 0.15 // 2% to 15% width

	if touchesLowerBand && lowVolatility && belowMiddleBand && normalVolatility {
		stopLoss := currentPrice - (indicators.ATR * 1.5)

		_, err := bot.client.EquityOrderPlaceMarket(ticker, positionSize)
		if err != nil {
			log.Printf("❌ BUY ERROR %s: %v", ticker, err)
		} else {
			log.Printf("🟢 BOUGHT %s: %d shares @ £%.2f | Stop: £%.2f | ATR: %.4f",
				ticker, positionSize, currentPrice, stopLoss, indicators.ATR)
		}
	} else {
		log.Printf("📊 %s: £%.2f | BB(%.2f-%.2f) | ATR: %.4f | No Entry",
			ticker, currentPrice, indicators.LowerBand, indicators.UpperBand, indicators.ATR)
	}
}

func (bot *TradingBot) evaluateExit(ticker string, currentPrice float64, indicators TradingIndicators, position trading212.Position) {
	entryPrice := position.Value / position.Quantity
	stopLoss := entryPrice - (indicators.ATR * 1.5)

	// Exit conditions
	var shouldSell bool
	var sellReason string

	// Take profit at upper band
	if currentPrice >= indicators.UpperBand {
		shouldSell = true
		sellReason = "Take Profit (Upper Band)"
	}

	// ATR-based stop loss
	if currentPrice <= stopLoss {
		shouldSell = true
		sellReason = "Stop Loss (ATR)"
	}

	// Profit protection: if 5%+ profit and price drops below middle band
	profitPercent := ((currentPrice - entryPrice) / entryPrice) * 100
	if profitPercent >= 5 && currentPrice < indicators.MiddleBand {
		shouldSell = true
		sellReason = "Profit Protection"
	}

	// High volatility exit
	if indicators.ATR > (entryPrice*0.04) && profitPercent > 0 {
		shouldSell = true
		sellReason = "High Volatility Exit"
	}

	// Time-based exit (if holding for too long without profit)
	if profitPercent < 0 && indicators.ATR < (entryPrice*0.01) {
		shouldSell = true
		sellReason = "Low Volatility Cut"
	}

	if shouldSell {
		_, err := bot.client.EquityOrderPlaceMarket(ticker, -int(position.Quantity))
		if err != nil {
			log.Printf("❌ SELL ERROR %s: %v", ticker, err)
		} else {
			pnl := (currentPrice - entryPrice) * position.Quantity
			log.Printf("🔴 SOLD %s: %d shares @ £%.2f | %s | P&L: £%.2f (%.1f%%)",
				ticker, int(position.Quantity), currentPrice, sellReason, pnl, profitPercent)
		}
	} else {
		log.Printf("🔵 HOLDING %s: Entry £%.2f | Current £%.2f | P&L: %.1f%% | Stop: £%.2f",
			ticker, entryPrice, currentPrice, profitPercent, stopLoss)
	}
}

func (bot *TradingBot) displayPortfolioSummary(positions []trading212.Position) {
	log.Println("=== Portfolio Summary ===")
	totalValue := 0.0
	totalPnL := 0.0

	for _, pos := range positions {
		if _, exists := bot.stocks[pos.Ticker]; exists {
			entryPrice := pos.Value / pos.Quantity
			bot.mutex.RLock()
			currentPrice := bot.stocks[pos.Ticker].CurrentPrice
			bot.mutex.RUnlock()

			if currentPrice > 0 {
				pnl := (currentPrice - entryPrice) * pos.Quantity
				pnlPercent := ((currentPrice - entryPrice) / entryPrice) * 100
				totalValue += pos.Value
				totalPnL += pnl

				log.Printf("%s: %d shares | Entry: £%.2f | Current: £%.2f | P&L: £%.2f (%.1f%%)",
					pos.Ticker, int(pos.Quantity), entryPrice, currentPrice, pnl, pnlPercent)
			}
		}
	}

	log.Printf("Total Portfolio Value: £%.2f | Total P&L: £%.2f", totalValue, totalPnL)
	log.Println("========================")
}

func (bot *TradingBot) getCurrentPrice(ticker string) float64 {
	positions, err := bot.client.Portfolio()
	if err != nil {
		log.Printf("Portfolio error for %s: %v", ticker, err)
		return 0
	}

	for _, pos := range positions {
		if pos.Ticker == ticker {
			return pos.Value / pos.Quantity
		}
	}

	// If no position exists, return simulated price
	// In production, you'd use market data API
	return 0
}

func (bot *TradingBot) calculateIndicators(priceHistory []PriceData) TradingIndicators {
	if len(priceHistory) < int(math.Max(float64(bot.atrPeriod), float64(bot.bollingerPeriod))) {
		return TradingIndicators{}
	}

	atr := bot.calculateATR(priceHistory)
	upperBand, middleBand, lowerBand, stdDev := bot.calculateBollingerBands(priceHistory)

	return TradingIndicators{
		ATR:        atr,
		UpperBand:  upperBand,
		MiddleBand: middleBand,
		LowerBand:  lowerBand,
		StdDev:     stdDev,
	}
}

func (bot *TradingBot) calculateATR(priceHistory []PriceData) float64 {
	if len(priceHistory) < bot.atrPeriod+1 {
		return 0
	}

	var trueRanges []float64

	for i := 1; i < len(priceHistory); i++ {
		current := priceHistory[i]
		previous := priceHistory[i-1]

		tr1 := current.High - current.Low
		tr2 := math.Abs(current.High - previous.Close)
		tr3 := math.Abs(current.Low - previous.Close)

		trueRange := math.Max(tr1, math.Max(tr2, tr3))
		trueRanges = append(trueRanges, trueRange)
	}

	if len(trueRanges) < bot.atrPeriod {
		return 0
	}

	sum := 0.0
	start := len(trueRanges) - bot.atrPeriod
	for i := start; i < len(trueRanges); i++ {
		sum += trueRanges[i]
	}

	return sum / float64(bot.atrPeriod)
}

func (bot *TradingBot) calculateBollingerBands(priceHistory []PriceData) (float64, float64, float64, float64) {
	if len(priceHistory) < bot.bollingerPeriod {
		return 0, 0, 0, 0
	}

	sum := 0.0
	start := len(priceHistory) - bot.bollingerPeriod
	for i := start; i < len(priceHistory); i++ {
		sum += priceHistory[i].Close
	}
	sma := sum / float64(bot.bollingerPeriod)

	variance := 0.0
	for i := start; i < len(priceHistory); i++ {
		diff := priceHistory[i].Close - sma
		variance += diff * diff
	}
	variance /= float64(bot.bollingerPeriod)
	stdDev := math.Sqrt(variance)

	upperBand := sma + (2 * stdDev)
	lowerBand := sma - (2 * stdDev)

	return upperBand, sma, lowerBand, stdDev
}
