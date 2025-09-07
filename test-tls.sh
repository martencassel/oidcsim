#!/bin/bash

openssl s_client -connect oidc.local:443 -cert $(pwd)/server.pem -key $(pwd)/server.key -CAfile /tmp/ca.pem
