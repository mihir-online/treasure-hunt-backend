#!/bin/bash

# Health Check Endpoint
# GET /status
# 
# Description: Check if the API is running and healthy
# Expected Response: 200 OK with status information

echo "=== Health Check ==="
echo ""

curl -X GET http://localhost:8008/status \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "=== Production ==="
echo ""

curl -X GET https://svc.mihir.page/status \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n"

