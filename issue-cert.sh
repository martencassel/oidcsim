#!/bin/bash
k get secret  root-ca-keypair  -n=cert-manager -o yaml

# Save the certificate
kubectl get secret root-ca-keypair -n cert-manager -o jsonpath='{.data.tls\.crt}' \
  | base64 -d > /tmp/ca.pem

# Save the private key
kubectl get secret root-ca-keypair -n cert-manager -o jsonpath='{.data.tls\.key}' \
  | base64 -d > /tmp/ca.key

./gen-cert.sh /tmp/ca.pem /tmp/ca.key out.pem out.key

openssl x509 -in out.pem -text -noout

