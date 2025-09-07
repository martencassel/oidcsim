#!/bin/bash
# Usage: ./add-ip.sh <ip/prefix>
# Example: ./add-ip.sh 172.21.22.200/20

set -e

IP_PREFIX="$1"

if [[ -z "$IP_PREFIX" ]]; then
  echo "Usage: $0 <ip/prefix>"
  exit 1
fi

sudo ip addr add "$IP_PREFIX" dev eth0
echo "Added $IP_PREFIX to eth0"
