# Setup routing for PTP bridge via the 5gs
# For network graph see docs/structure.drawio
# This can only run after the data session has been established by the 5gs (oaitun_ue1 created)

# UPF and UE use tunnel interfaces to communicate
docker exec oai-nr-ue ip r add 10.100.200.137 dev oaitun_ue1
docker exec upf ip r add 10.60.0.1 dev upfgtp