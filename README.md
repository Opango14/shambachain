# Shambachain

Hackathon project on Blockchain.

## API Endpoints

### Authentication Routes (Public)

- **`POST /api/auth/register`**
  Register a new user.
  **Expected Payload:**
  ```json
  {
    "username": "string (required)",
    "email": "string (required, valid email)",
    "password": "string (required, min 6 chars)"
  }
  ```

- **`POST /api/auth/login`**
  Login and receive an authentication token.
  **Expected Payload:**
  ```json
  {
    "email": "string (required, valid email)",
    "password": "string (required)"
  }
  ```

### Protected Routes (Requires Authentication)

*These endpoints require a valid JWT token in the `Authorization` header (`Bearer <token>`).*

- **`GET /api/user/profile`**
  Retrieve the profile details of the currently authenticated user.
  **Expected Payload:** None (uses token from header).

- **`POST /api/batches`**
  Register a new batch.
  **Expected Payload:**
  ```json
  {
    "produce_type": "string (required)",
    "quantity": "number (required, > 0)",
    "unit": "string (required)",
    "harvest_date": "string (ISO 8601 Date, required)",
    "location": "string (required)",
    "farm_name": "string (optional)"
  }
  ```

- **`POST /api/batches/:batchID/events`**
  Record a new event for a specific batch.
  **Expected Payload:**
  ```json
  {
    "event_type": "string (required, e.g. registration, transport, quality_check, transfer, sale)",
    "event_data": { 
      "key": "value (dynamic JSON object, required)" 
    }
  }
  ```

### Public Routes

- **`GET /api/marketplace`**
  Fetch all products currently available on the market (e.g. status is `registered` or `in_transit`).
  **Expected Payload:** None.

- **`GET /api/trace/:batchID`**
  Retrieve traceability information for a specific batch. Publicly accessible so buyers can scan QR codes.

- **`GET /qrcodes/:filename`**
  Serve a QR code image file.
