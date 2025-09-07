#!/usr/bin/env bash
set -euo pipefail

# ==== CONFIGURATION ====
AUTH_SERVER="https://auth.example.com"
AUTH_ENDPOINT="$AUTH_SERVER/authorize"
TOKEN_ENDPOINT="$AUTH_SERVER/token"
USERINFO_ENDPOINT="$AUTH_SERVER/userinfo"

CLIENT_ID="your_client_id"
CLIENT_SECRET="your_client_secret"   # omit if public client
REDIRECT_URI="https://client.example.com/callback"
SCOPES="openid profile email"
STATE="xyz123"

# ==== PKCE GENERATION ====
CODE_VERIFIER=$(openssl rand -base64 32 | tr -d '=+/')
CODE_CHALLENGE=$(echo -n "$CODE_VERIFIER" | openssl dgst -sha256 -binary | openssl base64 | tr -d '=+/')

echo "[*] PKCE code_verifier: $CODE_VERIFIER"
echo "[*] PKCE code_challenge: $CODE_CHALLENGE"

# ==== STEP 1: Get Authorization Code ====
# This assumes your server allows non-interactive login or has a test user session.
# If login is required, you may need to add `-d "username=...&password=..."` or cookies.

AUTH_URL="$AUTH_ENDPOINT?response_type=code&client_id=$CLIENT_ID&redirect_uri=$(printf %s "$REDIRECT_URI" | jq -s -R -r @uri)&scope=$(printf %s "$SCOPES" | jq -s -R -r @uri)&state=$STATE&code_challenge=$CODE_CHALLENGE&code_challenge_method=S256"

echo "[*] Requesting authorization code..."
REDIRECT_LOCATION=$(curl -s -i -L "$AUTH_URL" | grep -i "^location:" || true)

if [[ -z "$REDIRECT_LOCATION" ]]; then
  echo "[!] Could not capture redirect with code. You may need to log in manually."
  exit 1
fi

AUTH_CODE=$(echo "$REDIRECT_LOCATION" | grep -oP 'code=\K[^&]+')
echo "[*] Authorization code: $AUTH_CODE"

# ==== STEP 2: Exchange Code for Tokens ====
echo "[*] Exchanging code for tokens..."
TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_ENDPOINT" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=$AUTH_CODE" \
  -d "redirect_uri=$REDIRECT_URI" \
  -d "client_id=$CLIENT_ID" \
  -d "code_verifier=$CODE_VERIFIER" \
  -u "$CLIENT_ID:$CLIENT_SECRET")

echo "[*] Token response: $TOKEN_RESPONSE"

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
ID_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.id_token')

# ==== STEP 3: Call UserInfo Endpoint ====
if [[ "$ACCESS_TOKEN" != "null" && -n "$ACCESS_TOKEN" ]]; then
  echo "[*] Calling UserInfo endpoint..."
  curl -s -X GET "$USERINFO_ENDPOINT" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
else
  echo "[!] No access token received."
fi
