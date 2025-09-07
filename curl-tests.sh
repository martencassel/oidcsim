#!/bin/bash

set -e

curl -sk https://idp.local/.well-known/openid-configuration | jq

# Define variables
AUTH_URL="https://idp.local/oauth2/v1/authorize"
CLIENT_ID="web-client"
REDIRECT_URI="https://app.local/callback"
SCOPE="openid email"
STATE="xyz"
NONCE="abc"

# PKCE variables
CODE_VERIFIER=$(openssl rand -base64 32 | tr -d '=+/ ' | cut -c1-43)
CODE_CHALLENGE=$(printf "%s" "$CODE_VERIFIER" | openssl dgst -sha256 -binary | openssl base64 | tr -d '=+/ ' | tr 'A-Z' 'a-z' | cut -c1-43)

# Encode redirect URI and scope
ENCODED_REDIRECT_URI=$(printf "%s" "$REDIRECT_URI" | jq -s -R -r @uri)
ENCODED_SCOPE=$(printf "%s" "$SCOPE" | jq -s -R -r @uri)

# Make the request with PKCE
curl -sk "${AUTH_URL}?response_type=code&client_id=${CLIENT_ID}&redirect_uri=${ENCODED_REDIRECT_URI}&scope=${ENCODED_SCOPE}&state=${STATE}&nonce=${NONCE}&code_challenge=${CODE_CHALLENGE}&code_challenge_method=S256"

# /POST token
AUTH_CODE="code1234"
CLIENT_SECRET="secret123"
curl -sk -X POST https://idp.local/oauth2/v1/token \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "grant_type=authorization_code&code=${AUTH_CODE}&redirect_uri=${ENCODED_REDIRECT_URI}&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}"
