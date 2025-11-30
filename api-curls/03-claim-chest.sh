#!/bin/bash

# Claim Chest Endpoint
# POST /treasure/claim
#
# Description: Claim a treasure chest (supports secure and non-secure chests)
# Expected Response: 200 OK for successful claim, 409 Conflict for duplicate owner

# NOTE: Replace chest_id, source, and qr_match with actual values from your database
# You can get these values from the create-qr-codes endpoint response

echo "=== Claim Non-Secure Chest ==="
echo ""

curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "chest_id": "550e8400-e29b-41d4-a716-446655440000",
    "source": "Park Bench #42",
    "name": "Bob Smith",
    "email": "bob@example.com",
    "phone": "1234567890"
  }'

echo ""
echo ""
echo "=== Claim Secure Chest (with QR Match) ==="
echo ""

curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "chest_id": "660e8400-e29b-41d4-a716-446655440001",
    "source": "Museum Exhibit A",
    "qr_match": "1-X7K9P2M4",
    "name": "Alice Johnson",
    "email": "alice.j@example.com",
    "phone": "5551234567"
  }'

echo ""
echo ""
echo "=== Error Case - Invalid QR Match (Secure Chest) ==="
echo ""

curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "chest_id": "660e8400-e29b-41d4-a716-446655440001",
    "source": "Museum Exhibit A",
    "qr_match": "WRONG-CODE",
    "name": "Charlie Brown",
    "email": "charlie@example.com",
    "phone": "5559876543"
  }'

echo ""
echo ""
echo "=== Error Case - Duplicate Claim (409 Conflict) ==="
echo ""

curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "chest_id": "550e8400-e29b-41d4-a716-446655440000",
    "source": "Park Bench #42",
    "name": "Diana Prince",
    "email": "diana@example.com",
    "phone": "5551112222"
  }'

echo ""
echo ""
echo "=== Error Case - Source Mismatch ==="
echo ""

curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "chest_id": "550e8400-e29b-41d4-a716-446655440000",
    "source": "Wrong Source",
    "name": "Eve Adams",
    "email": "eve@example.com",
    "phone": "5553334444"
  }'

echo ""
echo ""
echo "=== Error Case - Chest Not Found ==="
echo ""

curl -X POST http://localhost:8008/treasure/claim \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -d '{
    "chest_id": "00000000-0000-0000-0000-000000000000",
    "source": "Nonexistent",
    "name": "Frank Castle",
    "email": "frank@example.com",
    "phone": "5555556666"
  }'

echo ""
echo ""
echo "=== Production Example ==="
echo ""

curl 'https://svc.mihir.page/treasure/claim' \
  -H 'Content-Type: application/json' \
  -H 'Referer: https://treasure.mihir.page/' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36' \
  -w "\nHTTP Status: %{http_code}\n" \
  --data-raw '{
    "chest_id": "abc",
    "source": "ABC",
    "name": "Mihir",
    "email": "mihir@example.com",
    "phone": "7428872048"
  }'

