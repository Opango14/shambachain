# Implementation Plan: Blockchain Produce Traceability System

## Overview

This implementation plan converts the blockchain produce traceability design into discrete coding tasks. The system will be built incrementally, starting with core data models and blockchain logic, then adding QR code generation, API handlers, and finally integration. Each task builds on previous work to ensure a cohesive, working system at each checkpoint.

## Tasks

- [x] 1. Set up project structure and core models
  - [x] 1.1 Create blockchain models (Block, Batch)
    - Create `models/block.go` with Block struct and all required fields
    - Create `models/batch.go` with Batch struct and all required fields
    - Add GORM tags for database mapping
    - Add JSON tags for API serialization
    - _Requirements: 1.1, 1.2, 2.1, 2.3_
  
  - [x] 1.2 Create event data structures
    - Create `models/events.go` with RegistrationEvent, TransportEvent, QualityCheckEvent, TransferEvent structs
    - Add JSON tags for all event structures
    - _Requirements: 7.1, 7.2, 8.1, 8.2, 9.1, 9.2_
  
  - [x] 1.3 Create API request/response structures
    - Create `models/requests.go` with RegisterBatchRequest, AddEventRequest structs
    - Create `models/responses.go` with RegisterBatchResponse, TraceabilityResponse structs
    - Add validation tags (binding:"required", binding:"gt=0")
    - _Requirements: 1.4, 1.5, 11.1, 11.2, 11.3, 11.4, 11.5_
  
  - [x] 1.4 Update database migrations
    - Update `database/database.go` to include Block and Batch models in AutoMigrate
    - Add database indices for BatchID and FarmerID
    - _Requirements: 1.1, 2.1_

- [-] 2. Implement core blockchain logic
  - [x] 2.1 Implement block hash computation
    - Create `blockchain/hash.go` with ComputeBlockHash function
    - Concatenate block fields in deterministic order
    - Use SHA-256 algorithm and return 64-character hex string
    - _Requirements: 2.6, 2.7, 2.8, 12.1_
  
  - [ ]* 2.2 Write property test for hash computation
    - **Property 1: Hash Determinism**
    - **Validates: Requirements 2.8, 12.1**
    - Test that identical block data produces identical hashes
  
  - [x] 2.3 Implement block creation function
    - Create `blockchain/block.go` with CreateBlock function
    - Set all block fields including timestamp
    - Compute and assign block hash
    - Validate event type and actor role
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_
  
  - [ ]* 2.4 Write property test for block creation
    - **Property 2: Block Hash Immutability**
    - **Validates: Requirements 4.2, 12.3**
    - Test that stored hash equals recomputed hash
  
  - [x] 2.5 Implement blockchain validation
    - Create `blockchain/validate.go` with ValidateChain function
    - Validate genesis block (index 0, prevHash "0")
    - Validate hash links between blocks
    - Validate sequential indices
    - Validate timestamp ordering
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 4.8_
  
  - [ ]* 2.6 Write property tests for blockchain validation
    - **Property 4: Chain Link Integrity**
    - **Validates: Requirements 2.2, 3.2, 4.3, 4.4**
    - **Property 5: Temporal Ordering**
    - **Validates: Requirements 4.5**
    - **Property 20: Chain Validation Correctness**
    - **Validates: Requirements 4.7**
    - **Property 21: Tamper Detection**
    - **Validates: Requirements 4.6, 12.4**

- [ ] 3. Checkpoint - Verify blockchain core logic
  - Ensure all tests pass, ask the user if questions arise.

- [~] 4. Implement QR code generation
  - [x] 4.1 Add QR code library dependency
    - Add `github.com/skip2/go-qrcode` to go.mod
    - Run `go mod tidy`
    - _Requirements: 5.1, 5.2_
  
  - [x] 4.2 Implement QR code generation function
    - Create `utils/qrcode.go` with GenerateQRCode function
    - Encode URL with batch ID: `https://domain.com/trace/{batchID}`
    - Save as PNG format with size validation (< 100KB)
    - Create qrcodes directory if it doesn't exist
    - Handle errors and cleanup partial files
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6_
  
  - [ ]* 4.3 Write property test for QR code generation
    - **Property 15: QR Code Generation**
    - **Validates: Requirements 1.6, 5.1, 5.2, 5.3, 5.4**
    - Test QR code format, size, and uniqueness

- [~] 5. Implement batch registration service
  - [x] 5.1 Create batch registration service function
    - Create `services/batch_service.go` with RegisterBatch function
    - Begin database transaction
    - Generate unique batch ID using UUID
    - Create batch record with status "registered"
    - Create genesis block with index 0 and prevHash "0"
    - Generate QR code and store path
    - Update batch with genesis hash and current hash
    - Commit transaction or rollback on error
    - _Requirements: 1.1, 1.2, 1.3, 1.6, 1.7, 1.8, 10.1, 10.2, 10.3_
  
  - [ ]* 5.2 Write property tests for batch registration
    - **Property 3: Genesis Block Properties**
    - **Validates: Requirements 1.2, 1.3, 4.1, 12.5**
    - **Property 16: Initial Batch Status**
    - **Validates: Requirements 1.1, 1.8**
    - **Property 28: Transaction Atomicity on Failure**
    - **Validates: Requirements 1.7, 3.8, 5.5, 10.2, 10.4**
  
  - [ ]* 5.3 Write unit tests for batch registration
    - Test successful registration flow
    - Test validation errors (negative quantity, future date, empty fields)
    - Test transaction rollback on QR generation failure
    - _Requirements: 1.4, 1.5, 11.1, 11.2, 11.3, 11.4, 11.5_

- [x] 6. Implement event recording service
  - [x] 6.1 Create add event service function
    - Create `services/event_service.go` with AddEvent function
    - Begin database transaction with row locking
    - Fetch batch and latest block
    - Create new block with incremented index
    - Set prevHash to latest block's hash
    - Save new block to database
    - Update batch current hash and status based on event type
    - Commit transaction or rollback on error
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 10.1, 10.2, 10.3, 10.5_
  
  - [ ]* 6.2 Write property tests for event recording
    - **Property 7: Sequential Block Indexing**
    - **Validates: Requirements 3.1**
    - **Property 6: Batch-Chain Consistency**
    - **Validates: Requirements 3.3, 6.4**
    - **Property 17: Status Transition on Transport**
    - **Validates: Requirements 3.4, 7.4**
    - **Property 18: Status Transition on Transfer**
    - **Validates: Requirements 3.5, 9.5**
    - **Property 19: Status Transition on Sale**
    - **Validates: Requirements 3.6**
  
  - [ ]* 6.3 Write unit tests for event recording
    - Test transport event recording
    - Test quality check event recording
    - Test transfer event recording
    - Test error handling for non-existent batch
    - Test concurrent event handling with locking
    - _Requirements: 7.1, 7.2, 7.3, 7.5, 8.1, 8.2, 8.3, 8.4, 8.5, 9.1, 9.2, 9.3, 9.4_

- [~] 7. Implement traceability retrieval service
  - [x] 7.1 Create traceability retrieval function
    - Create `services/traceability_service.go` with GetTraceability function
    - Fetch batch record by ID
    - Fetch all blocks for batch ordered by index
    - Validate blockchain using ValidateChain
    - Verify batch current hash matches last block hash
    - Return TraceabilityResponse with validation results
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_
  
  - [ ]* 7.2 Write property tests for traceability retrieval
    - **Property 22: Traceability Retrieval Completeness**
    - **Validates: Requirements 6.1, 6.2**
    - **Property 24: Read-Only Validation**
    - **Validates: Requirements 4.8, 6.6**
    - **Property 23: Non-Existent Batch Error**
    - **Validates: Requirements 3.7, 6.5**
  
  - [ ]* 7.3 Write unit tests for traceability retrieval
    - Test successful retrieval with valid chain
    - Test retrieval with tampered chain
    - Test error for non-existent batch
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

- [ ] 8. Checkpoint - Verify service layer
  - Ensure all tests pass, ask the user if questions arise.

- [-] 9. Implement API handlers
  - [x] 9.1 Create batch registration handler
    - Create `handlers/batch_handler.go` with RegisterBatchHandler
    - Extract farmer ID from authentication context
    - Bind and validate RegisterBatchRequest
    - Call RegisterBatch service function
    - Return RegisterBatchResponse with 201 status
    - Handle errors with appropriate HTTP status codes
    - _Requirements: 1.1, 1.4, 1.5, 11.1, 11.2, 11.3, 11.4, 11.5_
  
  - [ ] 9.2 Create event recording handler
    - Create `handlers/event_handler.go` with AddEventHandler
    - Extract batch ID from URL parameter
    - Extract actor ID and role from authentication context
    - Bind and validate AddEventRequest
    - Validate event type and event data structure
    - Call AddEvent service function
    - Return success response with 200 status
    - Handle errors with appropriate HTTP status codes
    - _Requirements: 2.4, 2.5, 3.7, 7.5, 8.5, 11.6, 11.7, 11.8_
  
  - [~] 9.3 Create traceability retrieval handler
    - Create `handlers/traceability_handler.go` with GetTraceabilityHandler
    - Extract batch ID from URL parameter or query string
    - Call GetTraceability service function
    - Return TraceabilityResponse with 200 status
    - Handle errors with appropriate HTTP status codes
    - _Requirements: 6.1, 6.5_
  
  - [~] 9.4 Create QR code serving handler
    - Create handler to serve QR code images from filesystem
    - Validate batch ID and file path
    - Return PNG image with appropriate content-type header
    - Handle file not found errors
    - _Requirements: 5.1, 5.6_
  
  - [ ]* 9.5 Write integration tests for API handlers
    - Test batch registration endpoint
    - Test event recording endpoint
    - Test traceability retrieval endpoint
    - Test QR code serving endpoint
    - Test error responses and validation

- [~] 10. Set up API routes
  - [~] 10.1 Create route configuration
    - Create `routes/routes.go` with SetupRoutes function
    - Define POST /api/batches for batch registration
    - Define POST /api/batches/:id/events for event recording
    - Define GET /api/trace/:id for traceability retrieval
    - Define GET /qrcodes/:filename for QR code serving
    - Apply authentication middleware to protected routes
    - _Requirements: 1.1, 3.1, 6.1_
  
  - [~] 10.2 Update main.go to use new routes
    - Import routes package
    - Initialize database with new models
    - Call SetupRoutes with Gin router
    - Remove old blockchain demo code
    - _Requirements: 1.1, 2.1_

- [~] 11. Integration and wiring
  - [~] 11.1 Wire all components together
    - Ensure database initialization includes all models
    - Ensure all handlers are registered with routes
    - Ensure QR code directory is created on startup
    - Add logging for key operations
    - _Requirements: 1.1, 2.1, 5.6_
  
  - [~] 11.2 Add input validation helpers
    - Create `utils/validation.go` with validation helper functions
    - Add ValidateEventType function
    - Add ValidateActorRole function
    - Add ValidateEventData function for each event type
    - _Requirements: 2.4, 2.5, 11.6, 11.7, 11.8_
  
  - [ ]* 11.3 Write end-to-end integration tests
    - Test complete workflow: register batch → add transport → add quality check → add transfer → retrieve traceability
    - Test blockchain validation with tampered data
    - Test concurrent event additions
    - _Requirements: All requirements_

- [ ] 12. Final checkpoint - Complete system verification
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at logical breakpoints
- Property tests validate universal correctness properties from the design document
- Unit tests validate specific examples and edge cases
- Integration tests validate end-to-end workflows
- The implementation uses Go with Gin framework and GORM as specified in the design
- QR code generation requires the `github.com/skip2/go-qrcode` library
- All database operations use transactions to ensure atomicity
