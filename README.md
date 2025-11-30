<!-- @format -->

# Treasure Hunt Backend API

A Go-based backend API for managing treasure hunt QR codes with batch creation, secure URLs, and automatic QR code image generation.

## 🚀 Features

- ✅ **Batch QR Code Creation** - Create multiple QR codes in a single API call
- ✅ **Source Uniqueness Validation** - Ensures no duplicate sources
- ✅ **Secure & Non-Secure Modes** - Control URL parameter inclusion
- ✅ **Auto-Generated QR Codes** - Automatic UUID and QR match generation
- ✅ **Image Generation** - PNG QR code images saved to disk
- ✅ **PostgreSQL Database** - Persistent storage with transactions
- ✅ **Clean Architecture** - Layered design (handlers → service → repository)

## 📋 Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Make

## 🛠️ Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd treasure-hunt-backend
```

### 2. Create Environment File

Copy the example environment file:

```bash
cp env.example .env
```

Edit `.env` with your configuration:

```bash
# Server Configuration
HOST=0.0.0.0
PORT=8008

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=dev_pg_user
DB_PASSWORD=dev_pg_password
DB_NAME=treasure_dev_db
DB_SSLMODE=disable

# Application Configuration
APP_ENV=development
LOG_LEVEL=info
ALLOWED_ORIGINS=*

# QR Code Configuration
QR_CODE_SIZE=256
QR_CODE_BASE_URL=https://svc.mihir.page
QR_PATH=/treasure
QR_STORAGE_DIR=qr_codes

# Geolocation Configuration
MAX_CLAIM_DISTANCE_METERS=100
```

⚠️ **Important:** All environment variables are **required**. The application will panic if any are missing.

### 3. Start Database

```bash
make dev_infra    # Start PostgreSQL in Docker
make tables       # Create database tables
```

### 4. Install Dependencies

```bash
make install
```

### 5. Run Server

```bash
make run
```

Server will start on `http://localhost:8008`

## 📡 API Endpoints

### Health Check

```bash
GET http://localhost:8008/status
```

**Response:**

```json
{
  "status": "healthy",
  "message": "Treasure Hunt API is running",
  "version": "1.0.0"
}
```

### Create QR Codes (Batch)

```bash
POST http://localhost:8008/api/v1/qr-codes
Content-Type: application/json

{
  "count": 10,
  "source": "Park Bench #42",
  "is_secure": true,
  "created_by": "alice@example.com"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Successfully created 10 treasure chests with QR codes",
  "source": "Park Bench #42",
  "count": 10,
  "qr_codes": [
    {
      "chest_id": "550e8400-e29b-41d4-a716-446655440000",
      "qr_match": "1-X7K9P2M4",
      "url": "https://svc.mihir.page/treasure?id=550e8400-e29b-41d4-a716-446655440000&source=Park+Bench+%2342&qr_match=1-X7K9P2M4",
      "file_path": "qr_codes/Park_Bench_42_1-X7K9P2M4.png"
    }
    // ... 9 more entries
  ]
}
```

### Claim Chest

**For Non-Secure Chests:**

```bash
POST http://localhost:8008/treasure/claim
Content-Type: application/json

{
  "chest_id": "550e8400-e29b-41d4-a716-446655440000",
  "source": "Park Bench #42",
  "name": "Bob Smith",
  "email": "bob@example.com",
  "phone": "1234567890"
}
```

**For Secure Chests (includes qr_match):**

```bash
POST http://localhost:8008/treasure/claim
Content-Type: application/json

{
  "chest_id": "550e8400-e29b-41d4-a716-446655440000",
  "source": "Park Bench #42",
  "qr_match": "1-X7K9P2M4",
  "name": "Bob Smith",
  "email": "bob@example.com",
  "phone": "1234567890"
}
```

**Success Response (200):**

```json
{
  "success": true,
  "message": "Chest claimed successfully!",
  "chest": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "source": "Park Bench #42",
    "qr_match": "1-X7K9P2M4",
    "is_secure": true,
    "status": "CLAIMED",
    "created_by": "alice@example.com",
    "created_at": "2025-11-30T10:30:00Z",
    "updated_at": "2025-11-30T11:00:00Z"
  }
}
```

**Duplicate Owner Response (409 Conflict):**

```json
{
  "success": false,
  "message": "This chest already has an owner",
  "chest": {...}
}
```

## 🏗️ Project Structure

```
treasure-hunt-backend/
├── cmd/api/                    # Application entry point
│   └── main.go
├── internal/
│   ├── models/                 # Data models and DTOs
│   │   ├── treasure.go
│   │   ├── requests.go
│   │   └── responses.go
│   ├── repository/             # Data access layer
│   │   ├── treasure_repository.go
│   │   └── postgres/
│   │       ├── db.go
│   │       └── treasure_chest_repo.go
│   ├── service/                # Business logic layer
│   │   └── treasure_service.go
│   ├── qrcode/                 # QR code generation
│   │   └── generator.go
│   └── handlers/               # HTTP handlers
│       ├── treasure_handler.go
│       ├── middleware.go
│       └── router.go
├── config/                     # Configuration management
│   └── config.go
├── db-scripts/                 # Database schemas
│   ├── treasure_chest.sql
│   ├── treasure_owner.sql
│   └── treasure_explorer.sql
├── qr_codes/                   # Generated QR code images
├── infra/                      # Infrastructure configs
│   └── docker.compose.dev.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🗄️ Database Schema

### treasure_chest Table

```sql
CREATE TABLE treasure_chest (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source VARCHAR(50) NOT NULL,
    qr_match VARCHAR(20) NOT NULL UNIQUE DEFAULT generate_qr_match(),
    is_secure BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(20) NOT NULL DEFAULT 'UNCLAIMED',
    created_by VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Auto-Generated Fields:**

- `id`: UUID (e.g., `550e8400-e29b-41d4-a716-446655440000`)
- `qr_match`: Format `{sequence}-{8_random_chars}` (e.g., `1-X7K9P2M4`)

## 🔒 Security Modes

### Secure Mode (`is_secure: true`)

- **URL Format:** `{base_url}/treasure?id={uuid}&source={source}&qr_match={qr_match}`
- **Parameters:** 3 (id, source, qr_match)
- **Security:** High - requires exact QR scan
- **Use Case:** High-value treasures, competitive hunts

### Non-Secure Mode (`is_secure: false`)

- **URL Format:** `{base_url}/treasure?id={uuid}&source={source}`
- **Parameters:** 2 (id, source)
- **Security:** Medium - can be manually entered
- **Use Case:** Promotional events, easy access

## 🧪 Testing

### Run Tests

```bash
make test
```

### Test Coverage

```bash
make test-coverage
```

### API Testing with Curl Scripts

Pre-configured curl scripts are available in the `api-curls/` directory:

```bash
# Make scripts executable
chmod +x api-curls/*.sh

# Test health check
./api-curls/01-health-check.sh

# Create QR codes
./api-curls/02-create-qr-codes.sh

# Claim chest
./api-curls/03-claim-chest.sh
```

See [api-curls/README.md](api-curls/README.md) for detailed usage.

### Manual Testing

```bash
# 1. Start server
make run

# 2. Test health endpoint
curl http://localhost:8008/status

# 3. Create QR codes
curl -X POST http://localhost:8008/api/v1/qr-codes \
  -H "Content-Type: application/json" \
  -d '{
    "count": 3,
    "source": "Test Location",
    "is_secure": true,
    "created_by": "tester@example.com"
  }'

# 4. Verify QR images
ls -l qr_codes/

# 5. Claim a chest (secure)
curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -d '{
    "chest_id": "550e8400-e29b-41d4-a716-446655440000",
    "source": "Test Location",
    "qr_match": "1-X7K9P2M4",
    "name": "Explorer Name",
    "email": "explorer@example.com",
    "phone": "1234567890"
  }'
```

## 🛠️ Development

### Available Make Commands

```bash
make help          # Show all available commands
make install       # Install dependencies
make build         # Build the application
make run           # Run the server
make test          # Run tests
make clean         # Clean build cache
make dev_infra     # Start database
make tables        # Create database tables
make formatting    # Setup code formatting tools
```

### Code Formatting

```bash
make formatting    # Setup pre-commit hooks
```

## 📦 Docker

### Build Docker Image

```bash
make docker-build
```

### Run in Docker

```bash
make docker-run
```

## 🔧 Configuration

All configuration is managed through environment variables. **No default values** are provided - all variables must be set in `.env`.

Required environment variables:

- `HOST`, `PORT`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `APP_ENV`, `LOG_LEVEL`, `ALLOWED_ORIGINS`
- `QR_CODE_SIZE`, `QR_CODE_BASE_URL`, `QR_PATH`, `QR_STORAGE_DIR`
- `MAX_CLAIM_DISTANCE_METERS`

## 🚦 Status Codes

| Code | Description                            |
| ---- | -------------------------------------- |
| 200  | Success                                |
| 201  | Created                                |
| 400  | Bad Request                            |
| 404  | Not Found                              |
| 409  | Conflict (e.g., source already exists) |
| 500  | Internal Server Error                  |

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Check if database is running
docker ps

# Restart database
make dev_infra
```

### Missing Environment Variables

```bash
# Error: "Environment variable PORT is required but not set"
# Solution: Copy env.example to .env and fill in all values
cp env.example .env
```

### QR Code Directory Issues

```bash
# Create directory manually if needed
mkdir -p qr_codes
chmod 755 qr_codes
```

## 📝 License

[Add your license here]

## 👥 Contributors

[Add contributors here]

---

**Status:** ✅ Production Ready

**Server:** `http://localhost:8008`

**Deployed:** `https://svc.mihir.page`
