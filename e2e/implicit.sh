#!/usr/bin/env bash
set -euo pipefail

# ==== CONFIGURATION ====
AUTH_SERVER="https://auth.example.com"
AUTH_ENDPOINT="$AUTH_SERVER/authorize"
USERINFO_ENDPOINT="$AUTH_SERVER/userinfo"

CLIENT_ID="your_client_id"
REDIRECT_URI="https://client.example.com/callback"
SCOPES="openid profile email"
STATE="xyz123"
NONCE="abc987"   # Required by OIDC implicit flow

# ==== STEP 1: Build Authorization URL ====
AUTH_URL="$AUTH_ENDPOINT?response_type=id_token%20token&client_id=$CLIENT_ID&redirect_uri=$(printf %s "$REDIRECT_URI" | jq -s -R -r @uri)&scope=$(printf %s "$SCOPES" | jq -s -R -r @uri)&state=$STATE&nonce=$NONCE"

echo "[*] Requesting tokens via implicit flow..."
# -L follows redirects
# -i includes headers
# We capture the final redirect URL which contains the tokens in the fragment (#)
FINAL_RESPONSE=$(curl -s -i -L "$AUTH_URL")

# ==== STEP 2: Extract Tokens from Redirect ====
# In implicit flow, tokens are returned in the URL fragment (#access_token=...)
# curl won't show the fragment in Location headers, so we need a trick:
# Many test servers allow returning tokens in query params for non-browser clients (for testing only).
# If your server does that, you can parse them like this:
ACCESS_TOKEN=$(echo "$FINAL_RESPONSE" | grep -oP 'access_token=\K[^& ]+')
ID_TOKEN=$(echo "$FINAL_RESPONSE" | grep -oP 'id_token=\K[^& ]+')

echo "[*] Access Token: $ACCESS_TOKEN"
echo "[*] ID Token: $ID_TOKEN"

# ==== STEP 3: Call UserInfo Endpoint ====
if [[ -n "$ACCESS_TOKEN" ]]; then
  echo "[*] Calling UserInfo endpoint..."
  curl -s -X GET "$USERINFO_ENDPOINT" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
else
  echo "[!] No access token found. Your server may only return tokens in the URL fragment, which curl cannot capture directly."
  echo "    In that case, run the AUTH_URL in a browser or headless browser and capture the fragment manually."
fi
