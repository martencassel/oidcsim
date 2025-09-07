#!/usr/bin/env bash
set -euo pipefail

# ==== CONFIGURATION ====
AUTH_SERVER="https://auth.example.com"
TOKEN_ENDPOINT="$AUTH_SERVER/token"
RESOURCE_API="https://api.example.com/data"

CLIENT_ID="your_client_id"
CLIENT_SECRET="your_client_secret"
SCOPE="api.read"   # Optional, depends on your server config

# ==== STEP 1: Request Access Token ====
echo "[*] Requesting access token via client_credentials..."
TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_ENDPOINT" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" \
  -d "scope=$SCOPE")

echo "[*] Token response: $TOKEN_RESPONSE"

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

if [[ -z "$ACCESS_TOKEN" || "$ACCESS_TOKEN" == "null" ]]; then
  echo "[!] Failed to obtain access token."
  exit 1
fi

echo "[*] Access Token: $ACCESS_TOKEN"

# ==== STEP 2: Call Protected API ====
echo "[*] Calling protected resource..."
curl -s -X GET "$RESOURCE_API" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
