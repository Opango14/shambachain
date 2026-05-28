# How to Run the Blockchain Produce Traceability System

## Prerequisites

1. **Go installed** (version 1.21 or higher)
   ```bash
   go version
   ```

2. **Git** (to clone if needed)

## Step 1: Install Dependencies

```bash
go mod tidy
```

This will download all required packages:
- Gin web framework
- GORM ORM
- SQLite driver
- QR code library
- UUID library

## Step 2: Build the Application

```bash
go build -o shambachain
```

This creates an executable named `shambachain`.

## Step 3: Run the Application

### Option A: Run the executable
```bash
./shambachain
```

### Option B: Run directly with Go
```bash
go run main.go
```

You should see:
```
Initializing database...
Database connnected
Starting server on :8080...
[GIN-debug] Listening and serving HTTP on :8080
```

## Step 4: Test the API

### 1. Register a Batch (requires authentication)

```bash
curl -X POST http://localhost:8080/api/batches \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-here" \
  -d '{
    "produce_type": "potatoes",
    "quantity": 500.0,
    "unit": "kg",
    "harvest_date": "2024-01-15T00:00:00Z",
    "location": "Nakuru Farm, GPS: -0.3031, 36.0800",
    "farm_name": "Green Valley Farm"
  }'
```

**Response:**
```json
{
  "batch_id": "550e8400-e29b-41d4-a716-446655440000",
  "qr_code_url": "/qrcodes/550e8400-e29b-41d4-a716-446655440000.png",
  "qr_code_data": "iVBORw0KGgoAAAANSUhEUgAA...",
  "genesis_hash": "a1b2c3d4e5f6..."
}
```

### 2. Add a Transport Event (requires authentication)

```bash
curl -X POST http://localhost:8080/api/batches/{batch_id}/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-here" \
  -d '{
    "event_type": "transport",
    "event_data": {
      "from_location": "Nakuru Farm",
      "to_location": "Nairobi Market",
      "transport_id": "TRK-001",
      "vehicle_info": "Truck KBZ 123A",
      "departure_time": "2024-01-15T08:00:00Z",
      "estimated_arrival": "2024-01-15T12:00:00Z"
    }
  }'
```

**Response:**
```json
{
  "message": "event added successfully",
  "batch_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "transport"
}
```

### 3. Get Traceability Information (public - no auth required)

```bash
curl http://localhost:8080/api/trace/{batch_id}
```

**Response:**
```json
{
  "batch": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "farmer_id": "farmer-123",
    "produce_type": "potatoes",
    "quantity": 500.0,
    "unit": "kg",
    "harvest_date": "2024-01-15T00:00:00Z",
    "location": "Nakuru Farm",
    "status": "in_transit",
    "qr_code_path": "./qrcodes/550e8400-e29b-41d4-a716-446655440000.png",
    "genesis_hash": "a1b2c3d4...",
    "current_hash": "b2c3d4e5...",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "blockchain": [
    {
      "id": 1,
      "batch_id": "550e8400-e29b-41d4-a716-446655440000",
      "index": 0,
      "timestamp": "2024-01-15T10:00:00Z",
      "event_type": "registration",
      "event_data": "{\"produce_type\":\"potatoes\",\"quantity\":500}",
      "actor_id": "farmer-123",
      "actor_role": "farmer",
      "prev_hash": "0",
      "hash": "a1b2c3d4...",
      "created_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "batch_id": "550e8400-e29b-41d4-a716-446655440000",
      "index": 1,
      "timestamp": "2024-01-15T10:30:00Z",
      "event_type": "transport",
      "event_data": "{\"from_location\":\"Nakuru Farm\",\"to_location\":\"Nairobi Market\"}",
      "actor_id": "transporter-001",
      "actor_role": "transporter",
      "prev_hash": "a1b2c3d4...",
      "hash": "b2c3d4e5...",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "verified": true,
  "chain_valid": true
}
```

### 4. Get QR Code Image (public)

```bash
curl http://localhost:8080/qrcodes/{batch_id}.png --output qrcode.png
```

Or open in browser:
```
http://localhost:8080/qrcodes/{batch_id}.png
```

## Authentication Setup

The current authentication middleware expects a Bearer token. You need to:

1. **For testing without real auth**, modify `auth/middleware.go` to set a test user:

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// For testing - set a test farmer ID
		ctx.Set("farmer_id", "test-farmer-123")
		ctx.Set("role", "farmer")
		ctx.Next()
	}
}
```

2. **For production**, implement proper JWT token validation in the middleware.

## Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./blockchain/...
go test ./services/...
go test ./handlers/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

## Project Structure

```
shambachain/
├── auth/              # Authentication middleware
├── blockchain/        # Blockchain core logic
├── database/          # Database initialization
├── handlers/          # HTTP request handlers
├── models/            # Data models
├── routes/            # Route configuration
├── services/          # Business logic
├── utils/             # Utility functions
├── qrcodes/           # Generated QR codes (created automatically)
├── main.go            # Application entry point
├── go.mod             # Go module dependencies
└── shambachain.db     # SQLite database (created automatically)
```

## API Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | /api/batches | Yes | Register new produce batch |
| POST | /api/batches/:batchID/events | Yes | Add event to batch |
| GET | /api/trace/:batchID | No | Get traceability info |
| GET | /qrcodes/:filename | No | Get QR code image |

## Troubleshooting

### Port already in use
```bash
# Change port in main.go
router.Run(":8081")  // Use different port
```

### Database locked
```bash
# Stop the application and delete the database
rm shambachain.db
# Restart the application
```

### QR codes not generating
```bash
# Ensure qrcodes directory exists
mkdir -p qrcodes
```

### Module errors
```bash
# Clean and reinstall dependencies
go clean -modcache
go mod tidy
```

## Production Deployment

1. **Set Gin to release mode**:
   ```go
   gin.SetMode(gin.ReleaseMode)
   ```

2. **Use PostgreSQL instead of SQLite**:
   ```go
   import "gorm.io/driver/postgres"
   
   dsn := "host=localhost user=postgres password=secret dbname=shambachain"
   db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
   ```

3. **Add proper authentication** (JWT tokens)

4. **Enable HTTPS**

5. **Add rate limiting**

6. **Set up logging**

7. **Configure CORS** if needed

## Next Steps

1. Implement proper JWT authentication
2. Add user registration and login endpoints
3. Add role-based access control
4. Implement pagination for blockchain retrieval
5. Add search and filter capabilities
6. Create a frontend application
7. Deploy to cloud (AWS, GCP, Azure)

## Support

For issues or questions, refer to:
- Design document: `.kiro/specs/blockchain-produce-traceability/design.md`
- Requirements: `.kiro/specs/blockchain-produce-traceability/requirements.md`
- Implementation summary: `IMPLEMENTATION_SUMMARY.md`
