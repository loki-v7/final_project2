# UDP Message Dashboard

A lightweight real-time dashboard that listens for UDP messages and displays them via a web interface. Built in Go with a dynamic frontend using JavaScript.

## 🔧 Features

- Listens for incoming UDP messages on port `9002`
- Stores and displays messages via a web interface on port `9001`
- Live updates with JavaScript (no page refresh required)
- JSON API for external access (`/live`)
- Graceful handling of incoming data
- Detects a `ServerStop` message to shut down the UDP listener

## 🧪 How to Run

### Prerequisites

- [Go (1.20+ recommended)](https://golang.org/dl/)
- `go.mod` for module dependencies

### Run the UDP Client

Navigate to the `client` directory and run the UDP client with:

```bash
cd client
go run .
```
### Run the server

```bash
go run main.go
