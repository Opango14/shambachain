# Requirements Document: Blockchain Produce Traceability System

## Introduction

This document specifies the requirements for a blockchain-powered farm produce traceability system. The system enables farmers to register produce batches with immutable blockchain records, generate QR codes for product identification, and provide transparent supply chain visibility from farm to market. All stakeholders (farmers, transporters, inspectors, buyers) can record events that are cryptographically linked in an immutable chain, ensuring data integrity and preventing tampering.

## Glossary

- **System**: The blockchain produce traceability backend service
- **Batch**: A registered unit of farm produce with unique identifier and blockchain
- **Block**: An immutable record in the blockchain representing a supply chain event
- **Genesis_Block**: The first block (index 0) in a batch's blockchain, created during registration
- **Chain**: The complete ordered sequence of blocks for a batch
- **Hash**: A SHA-256 cryptographic fingerprint of block data
- **Actor**: A user performing an action (farmer, transporter, inspector, buyer)
- **Event**: A supply chain occurrence recorded as a blockchain block
- **QR_Code**: A scannable code encoding the batch's traceability URL
- **Traceability**: The complete verifiable history of a produce batch

## Requirements

### Requirement 1: Batch Registration

**User Story:** As a farmer, I want to register my produce batches with harvest details, so that I can create an immutable record and generate a QR code for buyers to verify authenticity.

#### Acceptance Criteria

1. WHEN a farmer submits valid batch registration data THEN THE System SHALL create a unique batch record with status "registered"
2. WHEN a batch is registered THEN THE System SHALL create a genesis block with index 0 and previous hash "0"
3. WHEN a genesis block is created THEN THE System SHALL compute its hash from batch data and store it as the batch's genesis hash
4. WHEN batch registration includes quantity THEN THE System SHALL validate that quantity is positive (> 0)
5. WHEN batch registration includes harvest date THEN THE System SHALL validate that the date is not in the future
6. WHEN a batch is successfully registered THEN THE System SHALL generate a unique QR code image and store its file path
7. WHEN batch registration fails at any step THEN THE System SHALL rollback all changes and return an error
8. WHEN a batch is registered THEN THE System SHALL set the batch's current hash equal to the genesis hash

### Requirement 2: Blockchain Block Creation

**User Story:** As a system component, I want to create immutable blockchain blocks for supply chain events, so that all produce history is cryptographically secured and tamper-evident.

#### Acceptance Criteria

1. WHEN creating a block THEN THE System SHALL assign a sequential index starting from 0 for genesis blocks
2. WHEN creating a non-genesis block THEN THE System SHALL set the previous hash to the hash of the preceding block
3. WHEN creating any block THEN THE System SHALL record the current UTC timestamp
4. WHEN creating a block THEN THE System SHALL validate that the event type is one of: registration, transport, quality_check, transfer, sale
5. WHEN creating a block THEN THE System SHALL validate that the actor role is one of: farmer, transporter, inspector, buyer
6. WHEN creating a block THEN THE System SHALL compute the block hash from all block fields in deterministic order
7. WHEN computing a block hash THEN THE System SHALL use SHA-256 algorithm and return a 64-character hexadecimal string
8. WHEN the same block data is hashed multiple times THEN THE System SHALL produce identical hash values (determinism)

### Requirement 3: Supply Chain Event Recording

**User Story:** As a supply chain participant (farmer, transporter, inspector, buyer), I want to record events in the blockchain, so that the produce history is transparent and verifiable.

#### Acceptance Criteria

1. WHEN an actor adds an event to an existing batch THEN THE System SHALL create a new block with index equal to the previous maximum index plus 1
2. WHEN a new block is added THEN THE System SHALL set its previous hash to the current block's hash
3. WHEN a new block is added THEN THE System SHALL update the batch's current hash to the new block's hash
4. WHEN a transport event is added THEN THE System SHALL update the batch status to "in_transit"
5. WHEN a transfer event is added THEN THE System SHALL update the batch status to "delivered"
6. WHEN a sale event is added THEN THE System SHALL update the batch status to "sold"
7. WHEN adding an event to a non-existent batch THEN THE System SHALL return an error
8. WHEN event addition fails at any step THEN THE System SHALL rollback all changes atomically

### Requirement 4: Blockchain Validation

**User Story:** As a buyer, I want the system to validate the blockchain integrity, so that I can trust the produce history has not been tampered with.

#### Acceptance Criteria

1. WHEN validating a chain THEN THE System SHALL verify that the genesis block has index 0 and previous hash "0"
2. WHEN validating a chain THEN THE System SHALL verify that each block's hash matches the recomputed hash from its data
3. WHEN validating a chain THEN THE System SHALL verify that each block's previous hash matches the preceding block's hash
4. WHEN validating a chain THEN THE System SHALL verify that block indices are sequential (0, 1, 2, ...)
5. WHEN validating a chain THEN THE System SHALL verify that timestamps are monotonically increasing
6. WHEN any validation check fails THEN THE System SHALL return false and indicate the chain is invalid
7. WHEN all validation checks pass THEN THE System SHALL return true indicating the chain is valid
8. WHEN validating a chain THEN THE System SHALL not modify any block data

### Requirement 5: QR Code Generation

**User Story:** As a farmer, I want a unique QR code for each batch, so that buyers can easily scan and access the complete traceability information.

#### Acceptance Criteria

1. WHEN generating a QR code THEN THE System SHALL encode a URL containing the batch ID
2. WHEN generating a QR code THEN THE System SHALL save the image as PNG format
3. WHEN generating a QR code THEN THE System SHALL ensure the file size is reasonable (< 100KB)
4. WHEN a QR code is generated THEN THE System SHALL return the file path to the saved image
5. WHEN QR code generation fails THEN THE System SHALL return an error without leaving partial files
6. WHEN multiple batches are registered THEN THE System SHALL ensure each QR code file path is unique

### Requirement 6: Traceability Retrieval

**User Story:** As a buyer, I want to scan a QR code and view the complete produce history, so that I can verify authenticity and make informed purchasing decisions.

#### Acceptance Criteria

1. WHEN a buyer requests traceability for a batch ID THEN THE System SHALL retrieve the batch record and all associated blocks
2. WHEN retrieving blockchain blocks THEN THE System SHALL order them by index in ascending order
3. WHEN traceability is retrieved THEN THE System SHALL validate the entire blockchain and include the validation result
4. WHEN traceability is retrieved THEN THE System SHALL verify that the batch's current hash matches the last block's hash
5. WHEN the batch ID does not exist THEN THE System SHALL return an error
6. WHEN retrieving traceability THEN THE System SHALL not modify any batch or block data

### Requirement 7: Transport Event Recording

**User Story:** As a transporter, I want to record pickup and delivery details, so that the produce movement is tracked in the blockchain.

#### Acceptance Criteria

1. WHEN a transporter records a transport event THEN THE System SHALL store from_location, to_location, transport_id, and vehicle_info
2. WHEN a transport event is recorded THEN THE System SHALL store departure_time and estimated_arrival timestamps
3. WHEN a transport event is added THEN THE System SHALL create a block with event type "transport"
4. WHEN a transport event is added THEN THE System SHALL update the batch status to "in_transit"
5. WHEN transport event data is invalid or incomplete THEN THE System SHALL return an error

### Requirement 8: Quality Inspection Recording

**User Story:** As an inspector, I want to record quality checks and grades, so that produce quality is documented in the immutable blockchain.

#### Acceptance Criteria

1. WHEN an inspector records a quality check THEN THE System SHALL store inspector_id, inspector_name, grade, and notes
2. WHEN a quality check is recorded THEN THE System SHALL store temperature reading for perishable items
3. WHEN a quality check is recorded THEN THE System SHALL store a boolean passed/failed status
4. WHEN a quality check event is added THEN THE System SHALL create a block with event type "quality_check"
5. WHEN quality check data is invalid THEN THE System SHALL return an error

### Requirement 9: Ownership Transfer Recording

**User Story:** As a farmer or current owner, I want to record ownership transfers, so that the chain of custody is transparent and verifiable.

#### Acceptance Criteria

1. WHEN recording an ownership transfer THEN THE System SHALL store from_owner_id and to_owner_id
2. WHEN recording an ownership transfer THEN THE System SHALL store transfer_type (sale, donation, return)
3. WHEN recording a sale transfer THEN THE System SHALL store price and currency
4. WHEN a transfer event is added THEN THE System SHALL create a block with event type "transfer"
5. WHEN a transfer event is added THEN THE System SHALL update the batch status to "delivered"

### Requirement 10: Data Integrity and Atomicity

**User Story:** As a system administrator, I want all database operations to be atomic, so that the system maintains consistency even during failures.

#### Acceptance Criteria

1. WHEN any multi-step operation begins THEN THE System SHALL use a database transaction
2. WHEN any step in a transaction fails THEN THE System SHALL rollback all changes
3. WHEN all steps in a transaction succeed THEN THE System SHALL commit all changes atomically
4. WHEN a rollback occurs during batch registration THEN THE System SHALL delete any generated QR code files
5. WHEN concurrent events are added to the same batch THEN THE System SHALL use database locking to prevent race conditions

### Requirement 11: Input Validation

**User Story:** As a system component, I want to validate all inputs, so that invalid data cannot corrupt the blockchain or database.

#### Acceptance Criteria

1. WHEN receiving batch registration data THEN THE System SHALL validate that produce_type is non-empty
2. WHEN receiving batch registration data THEN THE System SHALL validate that quantity is positive
3. WHEN receiving batch registration data THEN THE System SHALL validate that unit is non-empty
4. WHEN receiving batch registration data THEN THE System SHALL validate that harvest_date is not in the future
5. WHEN receiving batch registration data THEN THE System SHALL validate that location is non-empty
6. WHEN receiving event data THEN THE System SHALL validate that event_type is one of the allowed types
7. WHEN receiving event data THEN THE System SHALL validate that actor_id is non-empty
8. WHEN receiving event data THEN THE System SHALL validate that event_data can be marshaled to valid JSON

### Requirement 12: Blockchain Immutability Guarantees

**User Story:** As a system stakeholder, I want mathematical guarantees that blockchain data cannot be altered, so that the system provides trustworthy traceability.

#### Acceptance Criteria

1. WHEN a block is created THEN THE System SHALL ensure its hash is computed from immutable block fields
2. WHEN a block is stored THEN THE System SHALL never allow modification of its hash, previous hash, or index
3. WHEN a block's data is read THEN THE System SHALL verify the stored hash matches the recomputed hash
4. WHEN any block in a chain is altered THEN THE System SHALL detect the tampering during validation
5. WHEN the genesis block is created THEN THE System SHALL ensure it has no valid previous block (previous hash "0")

