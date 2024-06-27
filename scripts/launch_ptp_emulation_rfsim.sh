echo "Bringing up 5g system" &&
docker compose --profile cn --profile ran-rfsim up -d &&

echo "Waiting 15s for data session establishment and gtp tunnel creation" &&
sleep 15 &&

echo "Launching ptp containers" &&
docker compose --profile ptpsim up -d &&

echo "Setting up routing for PTP packets" &&
# For network graph see docs/structure.drawio
# This can only run after the data session has been established by the 5gs (oaitun_ue1 created)

# UPF and UE use tunnel interfaces to communicate
docker exec oai-nr-ue ip r add 10.100.200.137 dev oaitun_ue1
docker exec upf ip r add 10.60.0.1 dev upfgtp
