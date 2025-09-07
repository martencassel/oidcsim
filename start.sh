#!/bin/bash
make allocate-ip
sudo make run ADDR=172.21.22.200 PORT=443 CERT=$(pwd)/server.pem KEY=$(pwd)/server.key
