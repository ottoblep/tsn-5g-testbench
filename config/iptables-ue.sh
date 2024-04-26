#!/bin/bash
#
# Configure iptables in UE
#
# Forward everything
iptables-nft -I FORWARD -j ACCEPT

# Since we can only use the UE as a destination we port forward traffic from the UE to the PTP-Slave
iptables-nft -t nat -I PREROUTING -p tcp --dport 2468 -j DNAT --to-destination 10.100.202.200:2468
iptables-nft -t nat -I PREROUTING -p udp --dport 2468 -j DNAT --to-destination 10.100.202.200:2468

# Allow return traffic
iptables-nft -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT