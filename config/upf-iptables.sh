#!/bin/bash
#
# Configure iptables in UPF
#
# This forwards all traffic from the UE to the eth1 interface
iptables -t nat -A POSTROUTING -o eth1  -j MASQUERADE
iptables -I FORWARD 1 -j ACCEPT

