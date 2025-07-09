## Trading212 Robo-Advisor

A simple automated investment system that manages your Trading212 pies each month.

### What It Does

The program automatically:
- Creates investment pies with different strategies
- Adds money to your pies each month
- Rebalances your investments when needed
- Creates monthly reports

### Investment Strategies

The system offers five different strategies:

**Conservative** - Safe investments (40% stocks, 40% bonds)
- Monthly investment: £500
- Goal: £50,000

**Balanced** - Mix of safe and growth investments
- Monthly investment: £750
- Goal: £75,000

**Tech Growth** - Focus on technology companies
- Monthly investment: £1,000
- Goal: £100,000

**Aggressive** - High growth potential
**Income** - Focus on dividends and income

### Setup

#### 1. Get Your API Key
- Log into Trading212
- Go to settings and create an API key
- Copy the key (keep it secret! 🤫)

#### 2. Set Environment Variables

You need to set two important settings:

```bash
# Your Trading212 API key (required)
export TRADING212_API_KEY="your_api_key_here"

# Use demo mode (recommended for testing)
export IS_DEMO="true"
```

**Important**: Start with `IS_DEMO="true"` to test safely with fake money.

#### 3. Run the Program

```bash
make roboadvisor
```

### What Happens When You Run It

1. **Checks your cash** - Sees how much money you have
2. **Creates pies** - Makes investment pies if they don't exist
3. **Checks balance** - Sees if your investments need rebalancing
4. **Adds money** - Puts your monthly amount into each pie
5. **Makes report** - Creates a file showing your portfolio

### Files Created

- `robo_advisor_YYYY-MM-DD_HH-MM-SS.log` - Activity log
- `monthly_report_YYYY-MM.json` - Monthly portfolio report

### Safety Features

- Defaults to demo mode if not specified
- Logs everything it does
- Won't invest if you don't have enough cash
- Stops investing when goals are reached

### Custom Configuration

You can set your own pie configurations using the `PIE_CONFIGURATIONS` environment variable with JSON format. If not set, it uses the default strategies above.

### Requirements

- [Go](https://go.dev) programming language
- [Trading212](https://trading212.com) account with API access

### Warning

The program can spend real money. Always test with demo mode first (`IS_DEMO="true"`) before using with real money.

---

**Happy Trading! 🚀📊😎💰**