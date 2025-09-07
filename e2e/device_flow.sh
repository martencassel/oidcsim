#!/usr/bin/env bash
set -euo pipefail

# ==== CONFIGURATION ====
AUTH_SERVER="https://auth.example.com"
DEVICE_ENDPOINT="$AUTH_SERVER/devicecode"   # sometimes /device_authorization
TOKEN_ENDPOINT="$AUTH_SERVER/token"
USERINFO_ENDPOINT="$AUTH_SERVER/userinfo"

CLIENT_ID="your_client_id"
CLIENT_SECRET="your_client_secret"   # omit if public client
SCOPE="openid profile email"

# ==== STEP 1: Request Device & User Codes ====
echo "[*] Requesting device and user codes..."
DEVICE_RESPONSE=$(curl -s -X POST "$DEVICE_ENDPOINT" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=$CLIENT_ID" \
  -d "scope=$SCOPE")

echo "[*] Device response: $DEVICE_RESPONSE"

DEVICE_CODE=$(echo "$DEVICE_RESPONSE" | jq -r '.device_code')
USER_CODE=$(echo "$DEVICE_RESPONSE" | jq -r '.user_code')
VERIFICATION_URI=$(echo "$DEVICE_RESPONSE" | jq -r '.verification_uri')
INTERVAL=$(echo "$DEVICE_RESPONSE" | jq -r '.interval')

echo "[*] Device Code: $DEVICE_CODE"
echo "[*] User Code: $USER_CODE"
echo "[*] Verification URI: $VERIFICATION_URI"
echo "[*] Polling interval: ${INTERVAL}s"

# ==== STEP 2: (Normally) User visits VERIFICATION_URI and enters USER_CODE ====
# Since your server can skip login, we simulate that step here.
# If your server auto-approves, you can skip this entirely.

# ==== STEP 3: Poll Token Endpoint ====
echo "[*] Polling token endpoint..."
while true; do
  TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_ENDPOINT" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=urn:ietf:params:oauth:grant-type:device_code" \
    -d "device_code=$DEVICE_CODE" \
    -d "client_id=$CLIENT_ID" \
    -d "client_secret=$CLIENT_SECRET")

  ERROR=$(echo "$TOKEN_RESPONSE" | jq -r '.error // empty')

  if [[ -n "$ERROR" ]]; then
    if [[ "$ERROR" == "authorization_pending" ]]; then
      echo "[*] Waiting for user authorization..."
      sleep "$INTERVAL"
      continue
    elif [[ "$ERROR" == "slow_down" ]]; then
      echo "[*] Told to slow down..."
      INTERVAL=$((INTERVAL + 5))
      sleep "$INTERVAL"
      continue
    else
      echo "[!] Error: $ERROR"
      exit 1
    fi
  else
    echo "[*] Token response: $TOKEN_RESPONSE"
    break
  fi
done

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
ID_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.id_token')

# ==== STEP 4: Call UserInfo Endpoint ====
if [[ -n "$ACCESS_TOKEN" && "$ACCESS_TOKEN" != "null" ]]; then
  echo "[*] Calling UserInfo endpoint..."
  curl -s -X GET "$USERINFO_ENDPOINT" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
else
  echo "[!] No access token received."
fi
