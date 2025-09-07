#!/bin/bash
# Finds a free IP in the eth0 /20 subnet and assigns it to eth0

set -e

# Get the base IP and netmask
SUBNET=$(ip -o -f inet addr show eth0 | awk '{print $4}')
NETMASK=$(echo $SUBNET | cut -d'/' -f2)
MYIP=$(hostname -I | awk '{print $1}')

# For /20: 172.21.16.1 to 172.21.31.254
OCTET1=$(echo $MYIP | cut -d. -f1)
OCTET2=$(echo $MYIP | cut -d. -f2)

echo "Searching for free IP in subnet..."
echo "Current IP: $MYIP"
echo "Subnet: $SUBNET"
echo "Netmask: $NETMASK"

for i in $(seq 16 31); do
  for j in $(seq 1 254); do
    IP="$OCTET1.$OCTET2.$i.$j"
    IP_NO_QUOTES=$(echo $IP | tr -d '"')
    if [[ "$IP_NO_QUOTES" == "$MYIP" ]]; then
      continue
    fi
    # Check if IP is already assigned to any interface
    ip addr show | grep -qw "$IP_NO_QUOTES" && continue
    # Check if IP responds to ARP (faster than ping for local net)
    arping -c 1 -w 1 "$IP_NO_QUOTES" &>/dev/null
    if [ $? -ne 0 ]; then
      echo "Assigning free IP: $IP_NO_QUOTES"
      sudo ip addr add $IP_NO_QUOTES/$NETMASK dev eth0
      exit 0
    fi
  done
done

echo "No free IP found in subnet."
