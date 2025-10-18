# Go Test Wallet

A REST API service for managing user wallets and balances, built with Go. The service supports concurrent operations (up to 1000 RPS per wallet) and provides endpoints for creating wallets, checking balances, and performing deposit/withdrawal operations.

## Features

- **Wallet Management**: Create, retrieve, and delete user wallets
- **Balance Operations**: Deposit and withdraw funds with concurrent safety
- **REST API**: OpenAPI 3.0 compliant endpoints
- **Database**: PostgreSQL with GORM ORM
- **Logging**: Structured logging with Zap
- **Docker Support**: Containerized deployment with multiple profiles (debug, test, release)
- **Testing**: Unit and integration tests included

## Tech Stack

- **Language**: Go 1.24.5
- **Framework**: Echo v4
- **Database**: PostgreSQL
- **ORM**: GORM
- **Configuration**: Viper
- **API Documentation**: OpenAPI 3.0
- **Logging**: Zap
- **Testing**: Testify

## Project Structure

```
.
├── api/
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # Custom middleware
│   └── openapi/           # Generated OpenAPI client/server code
├── build/                 # Dockerfiles for different environments
├── cmd/app/               # Application entry point
├── configs/               # Docker Compose configurations
├── internal/
│   ├── app/               # Business logic (services, repositories)
│   ├── config/            # Configuration management
│   └── models/            # Data models
├── test/                  # Test files
├── .env                   # Environment variables
├── go.mod                 # Go module definition
└── README.md
```

## Prerequisites

- Go 1.24.5 or later
- Docker and Docker Compose
- PostgreSQL (if running locally without Docker)

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/ichigo7diabol/go-test-wallet.git
   cd go-test-wallet
   ```

2. **Start the application in release mode**
   ```bash
   docker-compose --profile release up --build
   ```

   The API will be available at `http://localhost:8080/api/v1`

### Local Development

1. **Clone and setup**
   ```bash
   git clone https://github.com/ichigo7diabol/go-test-wallet.git
   cd go-test-wallet
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup PostgreSQL**
   - Install PostgreSQL locally, or
   - Use Docker: `docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=12345678 postgres:16`

4. **Configure environment**
   - Copy `.env` and adjust database connection if needed

5. **Run the application**
   ```bash
   go run cmd/app/main.go
   ```

## Development

### Debug Mode

Run with debugging enabled:
```bash
docker-compose --profile debug up --build
```

Connect debugger to port 40000.

### Testing

Run all tests:
```bash
docker-compose --profile test up --build
```

Or locally:
```bash
go test ./...
```

### Building

Build for production:
```bash
docker-compose --profile release up --build
```

## API Documentation

The API is documented using OpenAPI 3.0. The specification is available in `api/openapi.yaml`.

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Wallets

- `GET /wallets` - List all wallets
- `POST /wallets` - Create a new wallet
- `GET /wallet/{walletId}` - Get wallet information
- `DELETE /wallet/{walletId}` - Delete a wallet

#### Operations

- `POST /wallet` - Perform balance operation (DEPOSIT/WITHDRAW)

### Example Requests

**Create Wallet:**
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"initialBalance": 1000.00}'
```

**Deposit Funds:**
```bash
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "b1f04c42-2b54-4b73-996c-cc0d0579b5c0",
    "operationType": "DEPOSIT",
    "amount": 500.00
  }'
```

**Check Balance:**
```bash
curl http://localhost:8080/api/v1/wallet/b1f04c42-2b54-4b73-996c-cc0d0579b5c0
```

## Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `WALLET_APP_PORT` | Server port | 8080 |
| `WALLET_APP_DSN` | Database connection string | - |
| `WALLET_APP_DEBUG_PORT` | Debug port | 40000 |

Database environment variables (for Docker):
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `POSTGRES_PORT`

## License

This project is licensed under the MIT License.