## Trading212

[![Test](https://github.com/0xnu/trading212/actions/workflows/test.yaml/badge.svg)](https://github.com/0xnu/trading212/actions/workflows/test.yaml)
[![Release](https://img.shields.io/github/release/0xnu/trading212.svg)](https://github.com/0xnu/trading212/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xnu/trading212)](https://goreportcard.com/report/github.com/0xnu/trading212)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xnu/trading212.svg)](https://pkg.go.dev/github.com/0xnu/trading212)
[![License](https://img.shields.io/github/license/0xnu/trading212)](/LICENSE)

An unofficial Go library for interacting with the [Trading212](https://trading212.com) API.

> [!WARNING]
> Disclaimer: Using the Trading212 Golang library to create trading bots involves risks, including potential losses from market volatility and reliance on historical price patterns that may not predict future movements. External factors like news events and economic changes can also affect the bot's performance, so only invest money you can afford to lose.

### Generating API Key

Create an API key for Trading212 by following these steps:

+ Click the hamburger menu after logging into your Trading212 account
+ Scroll to the bottom
+ Click the green button "Switch to Practice"
+ Click Settings
+ Click "API (Beta)"
+ Click "Generate API key"

### Tests

Execute this command: `make test`

### Demo

+ Update the API Key with your own inside [here](./cmd/demo/main.go)
+ Execute this command: `make run`

### Single Stock Trading

+ Update the API Key with your own inside [here](./cmd/demo/nvidia.go)
+ Execute this command: `make trade`

### Multi-Stock Trading

+ Update the API Key with your own inside [here](./cmd/demo/multistock.go)
+ Execute this command: `make multistock`

### Using the Trading212 API

You can read the [API documentation](https://t212public-api-docs.redoc.ly/) to understand what's possible with the Trading212 API.

### License

This project is licensed under the [MIT License](./LICENSE).

### Copyright

(c) 2025 [Finbarrs Oketunji](https://finbarrs.eu). All Rights Reserved.