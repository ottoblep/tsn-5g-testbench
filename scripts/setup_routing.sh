# Setup route from ptp-master to ptp-slave via the 5gs
docker exec ptp-slave ip r add 10.100.201.0/24 via 10.100.202.100
docker exec oai-nr-ue ./bin/iptables.sh
docker exec oai-nr-ue ip r add 10.100.201.0/24 dev oaitun_ue1
docker exec upf ./iptables.sh
docker exec upf ip r add 10.100.202.0/24 dev upfgtp # This doesnt work the tunnel only accepts one destination IP 
docker exec ptp-master ip r add 10.100.202.0/24 via 10.100.201.100

### Debug tools
# docker exec upf apt update 
# docker exec upf apt install -y iputils-ping watch net-tools tcpdump
# docker exec oai-nr-ue apt update 
# docker exec oai-nr-ue apt install -y iputils-ping watch net-tools tcpdump