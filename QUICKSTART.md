# Quick Start Guide

## Run the Application in 3 Steps

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Run the Server
```bash
go run main.go
```

You should see:
```
Initializing database...
Database connnected
Starting server on :8080...
```

### 3. Test the API

**Register a Batch:**
```bash
curl -X POST http://localhost:8080/api/batches \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "produce_type": "potatoes",
    "quantity": 500.0,
    "unit": "kg",
    "harvest_date": "2024-01-15T00:00:00Z",
    "location": "Nakuru Farm",
    "farm_name": "Green Valley Farm"
  }'
```

**Add Transport Event:**
```bash
curl -X POST http://localhost:8080/api/batches/{BATCH_ID}/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "event_type": "transport",
    "event_data": {
      "from_location": "Nakuru Farm",
      "to_location": "Nairobi Market",
      "transport_id": "TRK-001",
      "vehicle_info": "Truck KBZ 123A"
    }
  }'
```

**Get Traceability:**
```bash
curl http://localhost:8080/api/trace/{BATCH_ID}
```

**View QR Code:**
Open in browser: `http://localhost:8080/qrcodes/{BATCH_ID}.png`

## Note on Authentication

The current auth middleware is basic. For testing, it accepts any Bearer token. 

To modify for testing without tokens, edit `auth/middleware.go`:

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set test user for development
		ctx.Set("farmer_id", "test-farmer-123")
		ctx.Set("role", "farmer")
		ctx.Next()
	}
}
```

## API Endpoints

- `POST /api/batches` - Register batch (auth required)
- `POST /api/batches/:batchID/events` - Add event (auth required)
- `GET /api/trace/:batchID` - Get traceability (public)
- `GET /qrcodes/:filename` - Get QR code (public)

## Files Created

The system automatically creates:
- `shambachain.db` - SQLite database
- `qrcodes/` - Directory for QR code images

## Full Documentation

See `RUN_INSTRUCTIONS.md` for complete documentation.
