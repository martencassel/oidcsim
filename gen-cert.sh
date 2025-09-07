#!/bin/bash
# Usage: ./gen-cert.sh <ca-cert.pem> <ca-key.pem> <cert.pem> <key.pem> [SANs...]
# Generates a server certificate and key signed by the provided CA, with optional SANs

set -e

CA_CERT="$1"
CA_KEY="$2"
CERT_PEM="$3"
KEY_PEM="$4"
shift 4
SANS=("$@")

if [ -z "$CA_CERT" ] || [ -z "$CA_KEY" ] || [ -z "$CERT_PEM" ] || [ -z "$KEY_PEM" ]; then
  echo "Usage: $0 <ca-cert.pem> <ca-key.pem> <cert.pem> <key.pem> [SANs...]"
  echo "Example: $0 ca-cert.pem ca-key.pem out.pem out.key DNS:localhost,IP:127.0.0.1"
  exit 1
fi

# Prepare SANs config if provided
EXTFILE=""
if [ ${#SANS[@]} -gt 0 ]; then
  EXTFILE="san.ext"
  echo "subjectAltName=$(IFS=,; echo "${SANS[*]}")" > "$EXTFILE"
fi

openssl genrsa -out "$KEY_PEM" 2048
openssl req -new -key "$KEY_PEM" -out server.csr -subj "/CN=localhost"
if [ -n "$EXTFILE" ]; then
  openssl x509 -req -in server.csr -CA "$CA_CERT" -CAkey "$CA_KEY" -CAcreateserial -out "$CERT_PEM" -days 365 -sha256 -extfile "$EXTFILE"
  rm "$EXTFILE"
else
  openssl x509 -req -in server.csr -CA "$CA_CERT" -CAkey "$CA_KEY" -CAcreateserial -out "$CERT_PEM" -days 365 -sha256
fi
rm server.csr

echo "Certificate and key generated: $CERT_PEM $KEY_PEM"
