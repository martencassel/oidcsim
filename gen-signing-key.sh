#!/bin/bash
# Generate a new RSA private key for signing (PEM format)
# Usage: ./gen-signing-key.sh signing-key.pem

set -e

KEYFILE="$1"
if [ -z "$KEYFILE" ]; then
  echo "Usage: $0 <output-key.pem>"
  exit 1
fi

openssl genpkey -algorithm RSA -out "$KEYFILE" -pkeyopt rsa_keygen_bits:2048
chmod 600 "$KEYFILE"
echo "Signing key generated: $KEYFILE"
