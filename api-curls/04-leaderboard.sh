#!/bin/bash

# Leaderboard API Test Script
# This script tests the leaderboard endpoint

BASE_URL="http://localhost:8080"

echo "=============================================="
echo "Testing Leaderboard API"
echo "=============================================="
echo ""

# Test 1: Get leaderboard for a source with default pagination
echo "Test 1: Get leaderboard (default pagination)"
echo "----------------------------------------------"
curl -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "test",
    "limit": 10,
    "offset": 0
  }'
echo -e "\n\n"

# Test 2: Get leaderboard with different limit
echo "Test 2: Get leaderboard with limit 5"
echo "----------------------------------------------"
curl -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "test",
    "limit": 5,
    "offset": 0
  }'
echo -e "\n\n"

# Test 3: Get leaderboard with offset (pagination)
echo "Test 3: Get leaderboard with offset"
echo "----------------------------------------------"
curl -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "test",
    "limit": 10,
    "offset": 10
  }'
echo -e "\n\n"

# Test 4: Get leaderboard for different source
echo "Test 4: Get leaderboard for different source"
echo "----------------------------------------------"
curl -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "event-2024",
    "limit": 10,
    "offset": 0
  }'
echo -e "\n\n"

# Test 5: Invalid request (missing source)
echo "Test 5: Invalid request - missing source"
echo "----------------------------------------------"
curl -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d '{
    "limit": 10,
    "offset": 0
  }'
echo -e "\n\n"

# Test 6: Large limit (should be capped at 100)
echo "Test 6: Large limit (should be capped at 100)"
echo "----------------------------------------------"
curl -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "test",
    "limit": 500,
    "offset": 0
  }'
echo -e "\n\n"

echo "=============================================="
echo "Leaderboard API Tests Complete"
echo "=============================================="

