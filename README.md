## Trading212

[![Test](https://github.com/0xnu/trading212/actions/workflows/test.yaml/badge.svg)](https://github.com/0xnu/trading212/actions/workflows/test.yaml)
[![Release](https://img.shields.io/github/release/0xnu/trading212.svg)](https://github.com/0xnu/trading212/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xnu/trading212)](https://goreportcard.com/report/github.com/0xnu/trading212)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xnu/trading212.svg)](https://pkg.go.dev/github.com/0xnu/trading212)
[![License](https://img.shields.io/github/license/0xnu/trading212)](/LICENSE)

An unofficial Go library for interacting with the [Trading212](https://trading212.com) API.

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

### Trading Example

+ Update the API Key with your own inside [here](./cmd/demo/nvidia.go)
+ Execute this command: `make trade`

### Using the Trading212 API

You can read the [API documentation](https://t212public-api-docs.redoc.ly/) to understand what's possible with the Trading212 API.

### License

This project is licensed under the [MIT License](./LICENSE).

### Copyright

(c) 2025 [Finbarrs Oketunji](https://finbarrs.eu). All Rights Reserved.