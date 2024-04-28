# Setup route from ptp-master to ptp-slave via the 5gs
# For network graph see docs/structure.drawio
# This can only run after the data session has been established by the 5gs (oaitun_ue1 created)

# PTP Containers are connected to UE and UPF
docker exec ptp-slave ip r add 10.100.201.0/24 via 10.100.202.100
docker exec ptp-master ip r add 10.60.0.0/24 via 10.100.201.100

# Forward everything
docker exec upf iptables -I FORWARD -j ACCEPT
docker exec oai-nr-ue iptables-nft -I FORWARD -j ACCEPT

# Allow return traffic
docker exec upf iptables -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT
docker exec oai-nr-ue iptables-nft -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT

# UE shall use the tunnel for upstream traffic
docker exec oai-nr-ue ip r add 10.100.201.0/24 dev oaitun_ue1

# We would like to do the same on the UPF side.
# However, the tunnel only accepts one destination IP (10.60.0.1) in the downstream direction
# docker exec upf ip r add 10.100.202.0/24 dev upfgtp # This doesnt work
# Instead we use the UE ip as the destination and port forward to the connected services
docker exec oai-nr-ue iptables-nft -t nat -I PREROUTING -p tcp --dport 2468 -j DNAT --to-destination 10.100.202.200:2468
docker exec oai-nr-ue iptables-nft -t nat -I PREROUTING -p udp --dport 2468 -j DNAT --to-destination 10.100.202.200:2468