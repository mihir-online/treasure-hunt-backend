<!-- @format -->

# Claim Chest API - Implementation Guide

## 🎯 Overview

The Claim Chest API allows users to claim treasure chests by providing chest details and user information. The API supports both secure and non-secure chests with different validation requirements.

## 📡 API Endpoint

```
POST http://localhost:8008/treasure/claim
```

Production: `https://svc.mihir.page/treasure/claim`

## 📋 Request Format

### For Non-Secure Chests

```json
{
  "chest_id": "550e8400-e29b-41d4-a716-446655440000",
  "source": "Park Bench #42",
  "name": "Bob Smith",
  "email": "bob@example.com",
  "phone": "1234567890"
}
```

### For Secure Chests (includes qr_match)

```json
{
  "chest_id": "550e8400-e29b-41d4-a716-446655440000",
  "source": "Park Bench #42",
  "qr_match": "1-X7K9P2M4",
  "name": "Bob Smith",
  "email": "bob@example.com",
  "phone": "1234567890"
}
```

## 📝 Request Parameters

| Field      | Type          | Required       | Description                                                     |
| ---------- | ------------- | -------------- | --------------------------------------------------------------- |
| `chest_id` | string (UUID) | ✅ Yes         | UUID of the chest to claim                                      |
| `source`   | string        | ✅ Yes         | Source identifier for verification                              |
| `qr_match` | string        | ⚠️ Secure Only | QR match code (e.g., "1-X7K9P2M4") - required for secure chests |
| `name`     | string        | ✅ Yes         | Claimer's full name                                             |
| `email`    | string        | ✅ Yes         | Claimer's email address                                         |
| `phone`    | string        | ✅ Yes         | Claimer's phone number                                          |

## 🔄 Business Logic Flow

### Step 1: Find Chest by ID

```
GET chest WHERE id = chest_id
↓
If not found → Return 400: "Chest not found"
```

### Step 2: Verify Source

```
IF chest.source != request.source
↓
Return 400: "Source mismatch"
```

### Step 3: Security Check

```
IF chest.is_secure == true
  ↓
  IF qr_match is empty
    → Return 400: "QR match is required for secure chests"
  ↓
  IF chest.qr_match != request.qr_match
    → Return 400: "Invalid QR match code"
```

### Step 4: Check Claim Status

```
IF chest.status == "CLAIMED"
↓
Return 400: "This chest has already been claimed"
```

### Step 5: Mark as Claimed

```
UPDATE treasure_chest
SET status = 'CLAIMED'
WHERE id = chest_id
```

### Step 6: Insert Explorer Entry

```
INSERT INTO treasure_explorer
  (chest_id, name, email, phone_number)
VALUES
  (...)
```

**Note:** Multiple explorers can be associated with the same chest.

### Step 7: Insert Owner Entry

```
INSERT INTO treasure_owner
  (chest_id, name, email, phone_number)
VALUES
  (...)
```

**Note:** Only ONE owner per chest (unique constraint on chest_id).

If duplicate → Return **409 Conflict**: "This chest already has an owner"

## 📊 Response Formats

### ✅ Success Response (200 OK)

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

### ❌ Error Responses

#### Chest Not Found (400 Bad Request)

```json
{
  "success": false,
  "message": "Chest with ID 'xxx' not found"
}
```

#### Source Mismatch (400 Bad Request)

```json
{
  "success": false,
  "message": "Source mismatch"
}
```

#### Invalid QR Match (400 Bad Request)

```json
{
  "success": false,
  "message": "Invalid QR match code"
}
```

#### Already Claimed (400 Bad Request)

```json
{
  "success": false,
  "message": "This chest has already been claimed",
  "chest": {...}
}
```

#### Duplicate Owner (409 Conflict)

```json
{
  "success": false,
  "message": "This chest already has an owner",
  "chest": {...}
}
```

## 🗄️ Database Tables

### treasure_chest

- **Status Update:** `status` changed from `UNCLAIMED` to `CLAIMED`

### treasure_explorer

- **Allows:** Multiple entries per chest_id
- **Tracks:** Everyone who attempted to claim

### treasure_owner

- **Constraint:** UNIQUE on chest_id
- **Purpose:** First successful claimer becomes the owner
- **Duplicate Error:** Returns 409 Conflict

## 🧪 Example Usage

### Example 1: Claim Non-Secure Chest

```bash
curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -d '{
    "chest_id": "abc123",
    "source": "Coffee Shop",
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "phone": "5551234567"
  }'
```

### Example 2: Claim Secure Chest

```bash
curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -d '{
    "chest_id": "550e8400-e29b-41d4-a716-446655440000",
    "source": "Museum Exhibit A",
    "qr_match": "1-X7K9P2M4",
    "name": "Bob Smith",
    "email": "bob@example.com",
    "phone": "5559876543"
  }'
```

### Example 3: Production (svc.mihir.page)

```bash
curl 'https://svc.mihir.page/treasure/claim' \
  -H 'Content-Type: application/json' \
  -H 'Referer: https://treasure.mihir.page/' \
  --data-raw '{
    "chest_id": "abc",
    "source": "ABC",
    "name": "Mihir",
    "email": "mihir@example.com",
    "phone": "7428872048"
  }'
```

## 🔒 Security Modes

### Secure Chests (`is_secure: true`)

**Requirements:**

- ✅ chest_id
- ✅ source
- ✅ **qr_match** (validated against chest.qr_match)
- ✅ name, email, phone

**Use Case:** High-value treasures requiring exact QR scan

### Non-Secure Chests (`is_secure: false`)

**Requirements:**

- ✅ chest_id
- ✅ source
- ❌ qr_match (not required)
- ✅ name, email, phone

**Use Case:** Promotional events, easy-access treasures

## 📊 Status Codes

| Code | Meaning               | Description                                             |
| ---- | --------------------- | ------------------------------------------------------- |
| 200  | OK                    | Chest claimed successfully                              |
| 400  | Bad Request           | Invalid request, chest not found, already claimed, etc. |
| 409  | Conflict              | Chest already has an owner (duplicate constraint)       |
| 500  | Internal Server Error | Database or server error                                |

## ⚠️ Important Notes

1. **Explorer vs Owner:**

   - **Explorer:** Multiple entries allowed per chest (tracks all claim attempts)
   - **Owner:** Only one entry per chest (first successful claimer)

2. **409 Conflict:**

   - Occurs when treasure_owner insert fails due to UNIQUE constraint on chest_id
   - Means someone else already claimed ownership
   - Explorer entry is still created

3. **Source Verification:**

   - Always validates source matches chest.source
   - Prevents claims from wrong URLs or QR codes

4. **QR Match Validation:**
   - Only enforced for secure chests
   - Must exactly match chest.qr_match
   - Prevents manual URL construction attacks

## 🧪 Testing

### Test Scenario 1: First Claim (Success)

```bash
# Create chest
curl -X POST http://localhost:8008/api/v1/qr-codes -d '{...}'

# Claim chest (should succeed)
curl -X POST http://localhost:8008/treasure/claim -d '{...}'
# Expected: 200 OK
```

### Test Scenario 2: Duplicate Claim (409 Conflict)

```bash
# Same chest, different user
curl -X POST http://localhost:8008/treasure/claim -d '{...}'
# Expected: 409 Conflict - "This chest already has an owner"
```

### Test Scenario 3: Invalid QR Match (Secure Chest)

```bash
curl -X POST http://localhost:8008/treasure/claim -d '{
  "chest_id": "...",
  "source": "...",
  "qr_match": "WRONG",
  ...
}'
# Expected: 400 Bad Request - "Invalid QR match code"
```

### Test Scenario 4: Already Claimed

```bash
# Claim same chest with same user
curl -X POST http://localhost:8008/treasure/claim -d '{...}'
# Expected: 400 Bad Request - "This chest has already been claimed"
```

## 🔧 Troubleshooting

### "Chest not found"

- Verify chest_id is correct UUID
- Check if chest exists in database

### "Source mismatch"

- Ensure source matches exactly (case-sensitive)
- Check for typos or URL encoding issues

### "Invalid QR match code"

- Only for secure chests
- Verify qr_match from QR code scan
- Check for typos or incorrect parsing

### "This chest already has an owner" (409)

- Normal behavior for duplicate claims
- First claimer wins ownership
- Subsequent claims are tracked as explorers only

## 📚 Database Schema

```sql
-- treasure_chest: Tracks chest status
UPDATE treasure_chest SET status = 'CLAIMED' WHERE id = ?

-- treasure_explorer: Multiple entries allowed
INSERT INTO treasure_explorer (chest_id, name, email, phone_number) VALUES (...)

-- treasure_owner: One entry per chest (UNIQUE constraint)
INSERT INTO treasure_owner (chest_id, name, email, phone_number) VALUES (...)
-- Duplicate: violates unique constraint -> 409 Conflict
```

---

**Status:** ✅ Fully Implemented

**Endpoint:** `POST /treasure/claim`

**Updated:** 2025-11-30
