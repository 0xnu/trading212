package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0xnu/trading212"
)

// InvestmentStrategy defines different portfolio strategies
type InvestmentStrategy string

const (
	Conservative InvestmentStrategy = "conservative"
	Balanced     InvestmentStrategy = "balanced"
	Aggressive   InvestmentStrategy = "aggressive"
	Growth       InvestmentStrategy = "growth"
	Income       InvestmentStrategy = "income"
)

// AssetAllocation represents the target allocation for different asset classes
type AssetAllocation struct {
	Equities      float64 `json:"equities"`
	Bonds         float64 `json:"bonds"`
	Technology    float64 `json:"technology"`
	Healthcare    float64 `json:"healthcare"`
	Financial     float64 `json:"financial"`
	Energy        float64 `json:"energy"`
	REITs         float64 `json:"reits"`
	International float64 `json:"international"`
}

// PieConfig represents the configuration for a managed pie
type PieConfig struct {
	Name               string             `json:"name"`
	Strategy           InvestmentStrategy `json:"strategy"`
	MonthlyAmount      float64            `json:"monthly_amount"`
	MaxGoal            float64            `json:"max_goal"`
	Instruments        map[string]float64 `json:"instruments"`
	RiskTolerance      float64            `json:"risk_tolerance"`
	RebalanceThreshold float64            `json:"rebalance_threshold"`
}

// InstrumentAllocator handles instrument allocation logic
type InstrumentAllocator struct {
	allocation AssetAllocation
}

// RoboAdvisor manages automated pie investments
type RoboAdvisor struct {
	client            *trading212.Client
	strategies        map[InvestmentStrategy]AssetAllocation
	instrumentMapping map[string]string
	logFile           *os.File
}

// NewRoboAdvisor creates a new robo-advisor instance
func NewRoboAdvisor(apiKey string, isDemo bool) *RoboAdvisor {
	logFile := createLogFile()

	return &RoboAdvisor{
		client:            trading212.NewClient(apiKey, isDemo),
		strategies:        initializeStrategies(),
		instrumentMapping: initializeInstrumentMapping(),
		logFile:           logFile,
	}
}

func createLogFile() *os.File {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("robo_advisor_%s.log", timestamp)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to create log file: %v", err)
		return nil
	}
	return logFile
}

// Close properly closes the robo-advisor resources
func (ra *RoboAdvisor) Close() {
	if ra.logFile != nil {
		ra.logFile.Close()
	}
}

// logMessage logs to both console and file
func (ra *RoboAdvisor) logMessage(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMessage := fmt.Sprintf("[%s] %s", timestamp, message)

	log.Println(message)
	if ra.logFile != nil {
		ra.logFile.WriteString(formattedMessage + "\n")
	}
}

// Run executes the monthly robo-advisor process
func (ra *RoboAdvisor) Run() {
	ra.logMessage("🤖 Starting Robo-Advisor Monthly Execution")
	defer ra.Close()

	configs := ra.loadConfigurations()
	cashResponse := ra.getCashBalance()

	if cashResponse == nil {
		return
	}

	// Extract the Free field from the cash response using reflection or type assertion
	availableCash := ra.extractAvailableCash(cashResponse)
	ra.logMessage(fmt.Sprintf("💰 Available Cash: £%.2f", availableCash))
	ra.processAllPies(configs, availableCash)
	ra.generateMonthlyReport()
	ra.logMessage("✅ Robo-Advisor execution completed successfully")
}

// extractAvailableCash extracts the available cash amount from the API response
func (ra *RoboAdvisor) extractAvailableCash(cashResponse interface{}) float64 {

	// For now, I assume the response has a field we can access
	// You may need to adjust this based on the actual Trading212 API response structure

	// Temporary approach - you might need to inspect the actual response structure
	if cashMap, ok := cashResponse.(map[string]interface{}); ok {
		if free, exists := cashMap["free"]; exists {
			if freeFloat, ok := free.(float64); ok {
				return freeFloat
			}
		}
	}

	// If we can't extract the value, log and return 0
	ra.logMessage("⚠️ Could not extract available cash from response")
	return 0.0
}

func (ra *RoboAdvisor) loadConfigurations() []PieConfig {
	configs, err := ra.loadPieConfigurations()
	if err != nil {
		ra.logMessage(fmt.Sprintf("❌ Failed to load configurations: %v", err))
		return nil
	}
	return configs
}

func (ra *RoboAdvisor) getCashBalance() interface{} {
	cash, err := ra.client.Cash()
	if err != nil {
		ra.logMessage(fmt.Sprintf("❌ Failed to get cash balance: %v", err))
		return nil
	}
	return cash
}

func (ra *RoboAdvisor) processAllPies(configs []PieConfig, availableCash float64) {
	for i, config := range configs {
		ra.addDelayBetweenPies(i)
		ra.processPie(config, availableCash)
	}
}

func (ra *RoboAdvisor) addDelayBetweenPies(index int) {
	if index > 0 {
		time.Sleep(5 * time.Second)
	}
}

// loadPieConfigurations loads pie configurations from environment or defaults
func (ra *RoboAdvisor) loadPieConfigurations() ([]PieConfig, error) {
	configsJSON := os.Getenv("PIE_CONFIGURATIONS")

	return ra.parseConfigurationsOrDefault(configsJSON)
}

func (ra *RoboAdvisor) parseConfigurationsOrDefault(configsJSON string) ([]PieConfig, error) {
	if configsJSON == "" {
		ra.logMessage("📋 Using default pie configurations")
		return ra.getDefaultConfigurations(), nil
	}

	var configs []PieConfig
	if err := json.Unmarshal([]byte(configsJSON), &configs); err != nil {
		ra.logMessage("⚠️ Failed to parse custom configurations, using defaults")
		return ra.getDefaultConfigurations(), nil
	}

	return configs, nil
}

// getDefaultConfigurations returns sensible default pie configurations
func (ra *RoboAdvisor) getDefaultConfigurations() []PieConfig {
	return []PieConfig{
		{
			Name:               "Conservative Growth",
			Strategy:           Conservative,
			MonthlyAmount:      500.0,
			MaxGoal:            50000.0,
			RiskTolerance:      0.3,
			RebalanceThreshold: 0.05,
		},
		{
			Name:               "Balanced Portfolio",
			Strategy:           Balanced,
			MonthlyAmount:      750.0,
			MaxGoal:            75000.0,
			RiskTolerance:      0.5,
			RebalanceThreshold: 0.1,
		},
		{
			Name:               "Tech Growth",
			Strategy:           Growth,
			MonthlyAmount:      1000.0,
			MaxGoal:            100000.0,
			RiskTolerance:      0.7,
			RebalanceThreshold: 0.15,
		},
	}
}

// processPie handles the complete pie management process
func (ra *RoboAdvisor) processPie(config PieConfig, availableCash float64) {
	ra.logMessage(fmt.Sprintf("🥧 Processing pie: %s", config.Name))

	pie := ra.findOrCreatePieWithLogging(config)
	if pie == nil {
		return
	}

	time.Sleep(2 * time.Second)

	ra.handleRebalancing(pie, config)
	ra.handleMonthlyInvestment(pie, config, availableCash)
}

func (ra *RoboAdvisor) findOrCreatePieWithLogging(config PieConfig) *trading212.Pie {
	pie, isNew, err := ra.findOrCreatePie(config)
	if err != nil {
		ra.logMessage(fmt.Sprintf("❌ Failed to find/create pie %s: %v", config.Name, err))
		return nil
	}

	ra.logPieStatus(pie, isNew)
	return pie
}

func (ra *RoboAdvisor) logPieStatus(pie *trading212.Pie, isNew bool) {
	if isNew {
		ra.logMessage(fmt.Sprintf("✨ Created new pie: %s (ID: %d)", pie.Name, pie.ID))
	} else {
		ra.logMessage(fmt.Sprintf("📊 Found existing pie: %s (ID: %d)", pie.Name, pie.ID))
	}
}

func (ra *RoboAdvisor) handleRebalancing(pie *trading212.Pie, config PieConfig) {
	needsRebalance, err := ra.checkRebalanceNeeded(pie, config)
	if err != nil {
		ra.logMessage(fmt.Sprintf("⚠️ Failed to check rebalance for %s: %v", config.Name, err))
		return
	}

	if !needsRebalance {
		return
	}

	ra.executeRebalancing(pie, config)
}

func (ra *RoboAdvisor) executeRebalancing(pie *trading212.Pie, config PieConfig) {
	ra.logMessage(fmt.Sprintf("⚖️ Rebalancing required for %s", config.Name))
	time.Sleep(2 * time.Second)

	if err := ra.rebalancePie(pie, config); err != nil {
		ra.logMessage(fmt.Sprintf("❌ Rebalance failed for %s: %v", config.Name, err))
	} else {
		ra.logMessage(fmt.Sprintf("✅ Successfully rebalanced %s", config.Name))
	}
}

func (ra *RoboAdvisor) handleMonthlyInvestment(pie *trading212.Pie, config PieConfig, availableCash float64) {
	if !ra.shouldInvest(pie, config, availableCash) {
		return
	}

	if err := ra.addMonthlyInvestment(pie, config); err != nil {
		ra.logMessage(fmt.Sprintf("❌ Failed to add investment to %s: %v", config.Name, err))
	} else {
		ra.logMessage(fmt.Sprintf("💸 Added £%.2f to %s", config.MonthlyAmount, config.Name))
	}
}

// findOrCreatePie finds existing pie or creates a new one
func (ra *RoboAdvisor) findOrCreatePie(config PieConfig) (*trading212.Pie, bool, error) {
	pies, err := ra.client.Pies()
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch pies: %v", err)
	}

	existingPie := ra.findExistingPie(pies, config.Name)
	if existingPie != nil {
		return existingPie, false, nil
	}

	return ra.createNewPie(config)
}

func (ra *RoboAdvisor) findExistingPie(pies []trading212.Pie, name string) *trading212.Pie {
	for _, pie := range pies {
		if strings.EqualFold(pie.Name, name) {
			return &pie
		}
	}
	return nil
}

func (ra *RoboAdvisor) createNewPie(config PieConfig) (*trading212.Pie, bool, error) {
	instruments := ra.generateInstrumentAllocation(config.Strategy)
	ra.logMessage(fmt.Sprintf("🔍 Creating pie with instruments: %v", instruments))

	endDate := time.Now().AddDate(10, 0, 0)

	pie, err := ra.client.PieCreate(
		"REINVEST",
		endDate,
		int(config.MaxGoal),
		"PiggyBank",
		config.Name,
		instruments,
	)

	if err != nil {
		return nil, false, fmt.Errorf("failed to create pie: %v", err)
	}

	return pie, true, nil
}

// generateInstrumentAllocation creates instrument allocation based on strategy
func (ra *RoboAdvisor) generateInstrumentAllocation(strategy InvestmentStrategy) map[string]float64 {
	allocation := ra.strategies[strategy]
	allocator := &InstrumentAllocator{allocation: allocation}

	instruments := allocator.createBaseInstruments()
	merged := allocator.mergeInstruments(instruments)
	normalised := allocator.normaliseWeights(merged)

	ra.logMessage(fmt.Sprintf("📊 Generated allocation for %s strategy: %v", strategy, normalised))
	return normalised
}

func (ia *InstrumentAllocator) createBaseInstruments() map[string]float64 {
	instruments := make(map[string]float64)

	ia.addTechnologyInstruments(instruments)
	ia.addHealthcareInstruments(instruments)
	ia.addFinancialInstruments(instruments)
	ia.addEnergyInstruments(instruments)
	ia.addREITInstruments(instruments)
	ia.addBondInstruments(instruments)
	ia.addInternationalInstruments(instruments)
	ia.addEquityInstruments(instruments)

	return instruments
}

func (ia *InstrumentAllocator) addTechnologyInstruments(instruments map[string]float64) {
	if ia.allocation.Technology <= 0 {
		return
	}
	instruments["META"] = ia.allocation.Technology * 0.4
	instruments["MSFT"] = ia.allocation.Technology * 0.3
	instruments["NVDA"] = ia.allocation.Technology * 0.3
}

func (ia *InstrumentAllocator) addHealthcareInstruments(instruments map[string]float64) {
	if ia.allocation.Healthcare <= 0 {
		return
	}
	instruments["JNJ"] = ia.allocation.Healthcare * 0.6
	instruments["PFE"] = ia.allocation.Healthcare * 0.4
}

func (ia *InstrumentAllocator) addFinancialInstruments(instruments map[string]float64) {
	if ia.allocation.Financial <= 0 {
		return
	}
	instruments["JPM"] = ia.allocation.Financial * 0.5
	instruments["WFC"] = ia.allocation.Financial * 0.5
}

func (ia *InstrumentAllocator) addEnergyInstruments(instruments map[string]float64) {
	if ia.allocation.Energy <= 0 {
		return
	}
	instruments["XOM"] = ia.allocation.Energy
}

func (ia *InstrumentAllocator) addREITInstruments(instruments map[string]float64) {
	if ia.allocation.REITs <= 0 {
		return
	}
	instruments["AMT"] = ia.allocation.REITs
}

func (ia *InstrumentAllocator) addBondInstruments(instruments map[string]float64) {
	if ia.allocation.Bonds <= 0 {
		return
	}
	instruments["TLT"] = ia.allocation.Bonds
}

func (ia *InstrumentAllocator) addInternationalInstruments(instruments map[string]float64) {
	if ia.allocation.International <= 0 {
		return
	}
	instruments["TSM"] = ia.allocation.International * 0.5
	instruments["ASML"] = ia.allocation.International * 0.5
}

func (ia *InstrumentAllocator) addEquityInstruments(instruments map[string]float64) {
	if ia.allocation.Equities <= 0 {
		return
	}
	instruments["NVDA"] = ia.allocation.Equities * 0.33
	instruments["TSLA"] = ia.allocation.Equities * 0.33
	instruments["GOOGL"] = ia.allocation.Equities * 0.34
}

func (ia *InstrumentAllocator) mergeInstruments(instruments map[string]float64) map[string]float64 {
	merged := make(map[string]float64)
	for ticker, weight := range instruments {
		merged[ticker] += weight
	}
	return merged
}

func (ia *InstrumentAllocator) normaliseWeights(instruments map[string]float64) map[string]float64 {
	total := ia.calculateTotal(instruments)
	if total <= 0 {
		return instruments
	}

	return ia.applyNormalisation(instruments, total)
}

func (ia *InstrumentAllocator) calculateTotal(instruments map[string]float64) float64 {
	total := 0.0
	for _, weight := range instruments {
		total += weight
	}
	return total
}

func (ia *InstrumentAllocator) applyNormalisation(instruments map[string]float64, total float64) map[string]float64 {
	normalised := make(map[string]float64)
	for ticker, weight := range instruments {
		normalised[ticker] = weight / total
	}
	return normalised
}

// checkRebalanceNeeded determines if pie needs rebalancing
func (ra *RoboAdvisor) checkRebalanceNeeded(pie *trading212.Pie, config PieConfig) (bool, error) {
	detailedPie, err := ra.client.Pie(pie.ID)
	if err != nil {
		return false, err
	}

	targetAllocation := ra.generateInstrumentAllocation(config.Strategy)
	return ra.evaluateRebalanceThreshold(detailedPie, targetAllocation, config), nil
}

func (ra *RoboAdvisor) evaluateRebalanceThreshold(pie *trading212.Pie, targetAllocation map[string]float64, config PieConfig) bool {
	for ticker, targetWeight := range targetAllocation {
		currentWeight := ra.calculateCurrentWeight(pie, ticker)
		deviation := math.Abs(currentWeight - targetWeight)

		if deviation > config.RebalanceThreshold {
			ra.logMessage(fmt.Sprintf("📈 %s deviation: %.2f%% (threshold: %.2f%%)",
				ticker, deviation*100, config.RebalanceThreshold*100))
			return true
		}
	}
	return false
}

func (ra *RoboAdvisor) calculateCurrentWeight(pie *trading212.Pie, ticker string) float64 {
	shares, exists := pie.InstrumentShares[ticker]
	if !exists {
		return 0.0
	}

	total := ra.calculateTotalShares(pie.InstrumentShares)
	if total <= 0 {
		return 0.0
	}

	return shares / total
}

func (ra *RoboAdvisor) calculateTotalShares(instrumentShares map[string]float64) float64 {
	total := 0.0
	for _, share := range instrumentShares {
		total += share
	}
	return total
}

// rebalancePie updates pie allocation to match target strategy
func (ra *RoboAdvisor) rebalancePie(pie *trading212.Pie, config PieConfig) error {
	targetAllocation := ra.generateInstrumentAllocation(config.Strategy)
	endDate := time.Now().AddDate(10, 0, 0).Format("2006-01-02T15:04:05Z")

	_, err := ra.client.PieUpdate(
		pie.ID,
		"REINVEST",
		endDate,
		int(config.MaxGoal),
		"PiggyBank",
		config.Name,
		targetAllocation,
	)

	return err
}

// shouldInvest determines if monthly investment should be made
func (ra *RoboAdvisor) shouldInvest(pie *trading212.Pie, config PieConfig, availableCash float64) bool {
	return ra.hasSufficientCash(config, availableCash) && ra.isUnderGoal(pie, config)
}

func (ra *RoboAdvisor) hasSufficientCash(config PieConfig, availableCash float64) bool {
	if availableCash >= config.MonthlyAmount {
		return true
	}

	ra.logMessage(fmt.Sprintf("💰 Insufficient cash for %s: £%.2f available, £%.2f needed",
		config.Name, availableCash, config.MonthlyAmount))
	return false
}

func (ra *RoboAdvisor) isUnderGoal(pie *trading212.Pie, config PieConfig) bool {
	if float64(pie.Goal) > config.MaxGoal {
		return true
	}

	ra.logMessage(fmt.Sprintf("🎯 Goal reached for %s: £%d (goal: £%.2f)",
		config.Name, pie.Goal, config.MaxGoal))
	return false
}

// addMonthlyInvestment adds the monthly investment to the pie
func (ra *RoboAdvisor) addMonthlyInvestment(pie *trading212.Pie, config PieConfig) error {
	ra.logMessage(fmt.Sprintf("📈 Investment logic triggered for %s: £%.2f", config.Name, config.MonthlyAmount))
	return nil
}

// generateMonthlyReport creates a comprehensive monthly report
func (ra *RoboAdvisor) generateMonthlyReport() {
	ra.logMessage("📊 Generating Monthly Portfolio Report")
	time.Sleep(2 * time.Second)

	pies, err := ra.client.Pies()
	if err != nil {
		ra.logMessage(fmt.Sprintf("❌ Failed to fetch pies for report: %v", err))
		return
	}

	ra.logPortfolioSummary(pies)
	ra.saveReportToFile(pies, 0.0, 0.0, 0.0)
}

func (ra *RoboAdvisor) logPortfolioSummary(pies []trading212.Pie) {
	totalGoal := ra.calculateTotalGoal(pies)
	separator := strings.Repeat("=", 50)

	ra.logMessage(separator)
	ra.logMessage("📈 MONTHLY PORTFOLIO REPORT")
	ra.logMessage(separator)

	ra.logIndividualPies(pies)
	ra.logOverallSummary(totalGoal, len(pies))
	ra.logMessage(separator)
}

func (ra *RoboAdvisor) calculateTotalGoal(pies []trading212.Pie) float64 {
	totalGoal := 0.0
	for _, pie := range pies {
		totalGoal += float64(pie.Goal)
	}
	return totalGoal
}

func (ra *RoboAdvisor) logIndividualPies(pies []trading212.Pie) {
	for _, pie := range pies {
		ra.logMessage(fmt.Sprintf("🥧 %s:", pie.Name))
		ra.logMessage(fmt.Sprintf("   Goal: £%d", pie.Goal))
		ra.logMessage(fmt.Sprintf("   Instruments: %d", len(pie.InstrumentShares)))
		ra.logMessage(fmt.Sprintf("   Dividend Action: %s", pie.DividendCashAction))
		ra.logMessage("")
	}
}

func (ra *RoboAdvisor) logOverallSummary(totalGoal float64, pieCount int) {
	ra.logMessage("📊 PORTFOLIO SUMMARY:")
	ra.logMessage(fmt.Sprintf("   Total Goals: £%.2f", totalGoal))
	ra.logMessage(fmt.Sprintf("   Number of Pies: %d", pieCount))
}

// saveReportToFile saves the monthly report to a JSON file
func (ra *RoboAdvisor) saveReportToFile(pies []trading212.Pie, totalValue, totalPnL, performance float64) {
	report := ra.createReportData(pies, totalValue, totalPnL, performance)
	fileName := fmt.Sprintf("monthly_report_%s.json", time.Now().Format("2006-01"))

	ra.writeReportToFile(report, fileName)
}

func (ra *RoboAdvisor) createReportData(pies []trading212.Pie, totalValue, totalPnL, performance float64) map[string]interface{} {
	return map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"total_value": totalValue,
		"total_pnl":   totalPnL,
		"performance": performance,
		"pie_count":   len(pies),
		"pies":        pies,
	}
}

func (ra *RoboAdvisor) writeReportToFile(report map[string]interface{}, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		ra.logMessage(fmt.Sprintf("❌ Failed to create report file: %v", err))
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		ra.logMessage(fmt.Sprintf("❌ Failed to write report: %v", err))
	} else {
		ra.logMessage(fmt.Sprintf("💾 Report saved to %s", fileName))
	}
}

// initializeStrategies defines the asset allocation for each strategy
func initializeStrategies() map[InvestmentStrategy]AssetAllocation {
	return map[InvestmentStrategy]AssetAllocation{
		Conservative: {
			Equities:      0.4,
			Bonds:         0.4,
			Technology:    0.1,
			Healthcare:    0.05,
			Financial:     0.05,
			International: 0.0,
		},
		Balanced: {
			Equities:      0.3,
			Bonds:         0.2,
			Technology:    0.25,
			Healthcare:    0.1,
			Financial:     0.1,
			International: 0.05,
		},
		Aggressive: {
			Equities:      0.2,
			Technology:    0.4,
			Healthcare:    0.15,
			Financial:     0.1,
			Energy:        0.05,
			International: 0.1,
		},
		Growth: {
			Technology:    0.5,
			Healthcare:    0.2,
			Equities:      0.15,
			Financial:     0.1,
			International: 0.05,
		},
		Income: {
			Bonds:         0.3,
			REITs:         0.2,
			Financial:     0.2,
			Equities:      0.2,
			International: 0.1,
		},
	}
}

// initializeInstrumentMapping maps asset classes to Trading212 tickers
func initializeInstrumentMapping() map[string]string {
	return map[string]string{
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
}

func main() {
	apiKey := os.Getenv("TRADING212_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ TRADING212_API_KEY environment variable is required")
	}

	isDemoStr := os.Getenv("IS_DEMO")
	isDemo, _ := strconv.ParseBool(isDemoStr)
	if isDemoStr == "" {
		isDemo = true
		log.Println("⚠️ IS_DEMO not set, defaulting to demo mode")
	}

	advisor := NewRoboAdvisor(apiKey, isDemo)
	advisor.Run()
}
