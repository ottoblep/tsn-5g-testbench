#!/bin/bash
#
# Configure iptables in UE
#
# Forward everything
iptables-nft -I FORWARD -j ACCEPT

# Masquerade outgoing traffic
iptables-nft -t nat -I POSTROUTING -o oaitun_ue1 -j MASQUERADE
iptables-nft -t nat -I POSTROUTING -o eth1 -j MASQUERADE

# Since we can only use the UE as a destination for the tunnel we port forward traffic from the UE to the PTP-Slave
iptables-nft -t nat -A PREROUTING -p tcp --dport 1111 -j DNAT --to-destination 10.100.202.200:1111

# Allow return traffic
iptables-nft -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT