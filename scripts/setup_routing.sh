# Setup route from ptp-master to ptp-slave via the 5gs

# PTP Containers are connected to UE and UPF
docker exec ptp-slave ip r add 10.100.201.0/24 via 10.100.202.100
docker exec ptp-master ip r add 10.100.202.0/24 via 10.100.201.100

# Setup forwarding in UPF and UE
docker exec oai-nr-ue ./bin/iptables.sh
docker exec upf ./iptables.sh

# UE shall use the tunnel for upstream traffic
docker exec oai-nr-ue ip r add 10.100.201.0/24 dev oaitun_ue1
# We would like to do the same on the UPF side.
# docker exec upf ip r add 10.100.202.0/24 dev upfgtp # This doesnt work
# However, the tunnel only accepts one destination IP (10.60.0.1) in the downstream direction
# Instead we use the UE ip as the destination and port forward to the connected services (see iptables-ue.sh)

### Debug tools
# docker exec upf apt update 
# docker exec upf apt install -y iputils-ping watch net-tools tcpdump
# docker exec oai-nr-ue apt update 
# docker exec oai-nr-ue apt install -y iputils-ping watch net-tools tcpdump