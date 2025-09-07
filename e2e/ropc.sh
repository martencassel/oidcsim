#!/usr/bin/env bash
set -euo pipefail

# ==== CONFIGURATION ====
AUTH_SERVER="https://auth.example.com"
TOKEN_ENDPOINT="$AUTH_SERVER/token"
USERINFO_ENDPOINT="$AUTH_SERVER/userinfo"

CLIENT_ID="your_client_id"
CLIENT_SECRET="your_client_secret"   # omit if public client
USERNAME="testuser@example.com"
PASSWORD="SuperSecretPassword"
SCOPE="openid profile email"

# ==== STEP 1: Request Access Token ====
echo "[*] Requesting access token via ROPC..."
TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_ENDPOINT" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "username=$USERNAME" \
  -d "password=$PASSWORD" \
  -d "scope=$SCOPE" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET")

echo "[*] Token response: $TOKEN_RESPONSE"

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
ID_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.id_token
