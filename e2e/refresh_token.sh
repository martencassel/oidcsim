#!/usr/bin/env bash
set -euo pipefail

# ==== CONFIGURATION ====
AUTH_SERVER="https://auth.example.com"
TOKEN_ENDPOINT="$AUTH_SERVER/token"
USERINFO_ENDPOINT="$AUTH_SERVER/userinfo"

CLIENT_ID="your_client_id"
CLIENT_SECRET="your_client_secret"   # omit if public client
REFRESH_TOKEN="your_refresh_token_here"   # obtained from a previous login flow

# ==== STEP 1: Request New Access Token ====
echo "[*] Requesting new access token via refresh_token..."
TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_ENDPOINT" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=refresh_token" \
  -d "refresh_token=$REFRESH_TOKEN" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET")

echo "[*] Token response: $TOKEN_RESPONSE"

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
NEW_REFRESH_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.refresh_token // empty')

if [[ -z "$ACCESS_TOKEN" || "$ACCESS_TOKEN" == "null" ]]; then
  echo "[!] Failed to obtain new access token."
  exit 1
fi

echo "[*] New Access Token: $ACCESS_TOKEN"
if [[ -n "$NEW_REFRESH_TOKEN" ]]; then
  echo "[*] New Refresh Token: $NEW_REFRESH_TOKEN"
fi

# ==== STEP 2: Call UserInfo Endpoint (Optional) ====
echo "[*] Calling UserInfo endpoint..."
curl -s -X GET "$USERINFO_ENDPOINT" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
