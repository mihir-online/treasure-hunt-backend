<!-- @format -->

# API Curl Examples

This directory contains shell scripts with curl commands for testing all API endpoints.

## 📁 Files

| File                    | Endpoint                | Description                 |
| ----------------------- | ----------------------- | --------------------------- |
| `01-health-check.sh`    | `GET /status`           | Health check and API status |
| `02-create-qr-codes.sh` | `POST /api/v1/qr-codes` | Create QR codes in batch    |
| `03-claim-chest.sh`     | `POST /treasure/claim`  | Claim treasure chests       |
| `04-leaderboard.sh`     | `POST /leaderboard`     | Get player leaderboard      |

## 🚀 Usage

### Make Scripts Executable

```bash
chmod +x api-curls/*.sh
```

### Run Individual Scripts

```bash
# Health check
./api-curls/01-health-check.sh

# Create QR codes
./api-curls/02-create-qr-codes.sh

# Claim chest
./api-curls/03-claim-chest.sh

# Get leaderboard
./api-curls/04-leaderboard.sh
```

### Run All Scripts

```bash
# Execute all scripts in order
for script in api-curls/*.sh; do
  echo "===================="
  echo "Running: $script"
  echo "===================="
  bash "$script"
  echo ""
  echo ""
done
```

## ⚙️ Configuration

### Local Development

All scripts default to `http://localhost:8008`

Make sure your server is running:

```bash
make run
```

### Production

Production examples use `https://svc.mihir.page`

## 📝 Notes

### 01-health-check.sh

- Simple GET request
- No authentication required
- Tests both local and production endpoints

### 02-create-qr-codes.sh

Includes examples for:

- ✅ Non-secure chests (5 codes)
- ✅ Secure chests (3 codes)
- ✅ Large batch (50 codes)
- ❌ Duplicate source error
- 🌐 Production example

**Important:** The `source` must be unique for each batch!

### 03-claim-chest.sh

Includes examples for:

- ✅ Non-secure chest claim
- ✅ Secure chest claim (with qr_match)
- ❌ Invalid QR match
- ❌ Duplicate claim (409 Conflict)
- ❌ Source mismatch
- ❌ Chest not found
- 🌐 Production example

**Important:** Replace `chest_id`, `source`, and `qr_match` with actual values from your database!

### 04-leaderboard.sh

Includes examples for:

- ✅ Get top 10 players
- ✅ Pagination (limit and offset)
- ✅ Different sources
- ❌ Invalid request (missing source)
- ✅ Large limit (capped at 100)

**Note:** Leaderboard shows aggregated scores from the `treasure_explorer` table grouped by player.

## 🔧 Customization

### Update Chest IDs

After creating QR codes, update the chest IDs in `03-claim-chest.sh`:

```bash
# 1. Create QR codes
./api-curls/02-create-qr-codes.sh

# 2. Copy chest_id from response
# Example: "chest_id": "550e8400-e29b-41d4-a716-446655440000"

# 3. Update 03-claim-chest.sh with actual chest_id and qr_match
```

### Update URLs

To test against a different server:

```bash
# Edit the scripts and replace:
http://localhost:8008  # with your server URL
```

## 📊 Expected Responses

### Health Check (200 OK)

```json
{
  "status": "healthy",
  "message": "Treasure Hunt API is running",
  "version": "1.0.0"
}
```

### Create QR Codes (201 Created)

```json
{
  "success": true,
  "message": "Successfully created 5 treasure chests with QR codes",
  "source": "Park Bench #42",
  "count": 5,
  "qr_codes": [...]
}
```

### Claim Chest (200 OK)

```json
{
  "success": true,
  "message": "Chest claimed successfully!",
  "chest": {...}
}
```

### Duplicate Owner (409 Conflict)

```json
{
  "success": false,
  "message": "This chest already has an owner",
  "chest": {...}
}
```

### Leaderboard (200 OK)

```json
{
  "leaderboard": [
    {
      "name": "Alice Johnson",
      "rank": 1,
      "points": 500
    }
  ],
  "total": 25,
  "source": "test"
}
```

## 🧪 Testing Workflow

### Complete Test Flow

```bash
# 1. Health check
./api-curls/01-health-check.sh

# 2. Create test chests
./api-curls/02-create-qr-codes.sh

# 3. Copy chest_id from response

# 4. Edit 03-claim-chest.sh with actual chest_id

# 5. Claim chest
./api-curls/03-claim-chest.sh
```

### Verify Database

```bash
# Connect to database
docker exec -it treasure_dev_postgres psql -U dev_pg_user -d treasure_dev_db

# Check created chests
SELECT id, source, qr_match, is_secure, status FROM treasure_chest;

# Check explorers
SELECT id, chest_id, name, email FROM treasure_explorer;

# Check owners
SELECT id, chest_id, name, email FROM treasure_owner;
```

## 🔍 Debugging

### View Full Response

Remove the `-w` flag to see only response body:

```bash
curl -X GET http://localhost:8008/status \
  -H "Content-Type: application/json"
```

### View Headers

Add `-i` flag:

```bash
curl -i -X GET http://localhost:8008/status \
  -H "Content-Type: application/json"
```

### Verbose Output

Add `-v` flag:

```bash
curl -v -X GET http://localhost:8008/status \
  -H "Content-Type: application/json"
```

### Pretty Print JSON

Pipe to `jq`:

```bash
curl -X GET http://localhost:8008/status \
  -H "Content-Type: application/json" | jq '.'
```

## 📚 Additional Resources

- [QR Code Creation Guide](../QR_CODE_GENERATION.md)

## ⚠️ Important Notes

1. **UUIDs:** Replace example UUIDs with actual values from your database
2. **Sources:** Each source must be unique across all chests
3. **QR Match:** Only required for secure chests
4. **409 Conflicts:** Normal behavior when chest already has owner
5. **Port 8008:** Default local port (update if changed in .env)

---

**Status:** ✅ Ready to Use

**Last Updated:** 2025-11-30
