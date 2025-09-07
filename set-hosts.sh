#!/bin/bash

BIND_IP=172.21.22.200
if [ -n "$BIND_IP" ]; then
  grep -q "oidc.local" /etc/hosts || echo "$BIND_IP    oidc.local" | sudo tee -a /etc/hosts
else
  echo "Could not determine oidc.local bind IP. Skipping /etc/hosts update for oidc.local."
fi
