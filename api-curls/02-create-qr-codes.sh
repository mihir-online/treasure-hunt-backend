#!/bin/bash

# Create QR Codes Endpoint
# POST /api/v1/qr-codes
#
# Description: Create multiple treasure chests with QR codes in batch
# Expected Response: 201 Created with QR code details

echo "=== Create QR Codes - Non-Secure (5 codes) ==="
echo ""

curl -X POST http://localhost:8008/api/v1/qr-codes \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "count": 5,
    "source": "Park Bench #42",
    "is_secure": false,
    "first_score": 100,
    "subsequent_score": 50,
    "created_by": "alice@example.com"
  }'

echo ""
echo ""
echo "=== Create QR Codes - Secure (3 codes) ==="
echo ""

curl -X POST http://localhost:8008/api/v1/qr-codes \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "count": 3,
    "source": "Museum Exhibit A",
    "is_secure": true,
    "first_score": 100,
    "subsequent_score": 50,
    "created_by": "curator@museum.com"
  }'

echo ""
echo ""
echo "=== Create QR Codes - Large Batch (50 codes) ==="
echo ""

curl -X POST http://localhost:8008/api/v1/qr-codes \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "count": 50,
    "source": "City Scavenger Hunt 2025",
    "is_secure": true,
    "first_score": 100,
    "subsequent_score": 50,
    "created_by": "event@cityhunt.com"
  }'

echo ""
echo ""
echo "=== Error Case - Duplicate Source ==="
echo ""

curl -X POST http://localhost:8008/api/v1/qr-codes \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "count": 3,
    "source": "Park Bench #42",
    "is_secure": false,
    "first_score": 100,
    "subsequent_score": 50,
    "created_by": "alice@example.com"
  }'

echo ""
echo ""
echo "=== Production Example ==="
echo ""

curl -X POST https://svc.mihir.page/api/v1/qr-codes \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "count": 10,
    "source": "Production Test Event",
    "is_secure": true,
    "first_score": 100,
    "subsequent_score": 50,
    "created_by": "admin@mihir.page"
  }'

