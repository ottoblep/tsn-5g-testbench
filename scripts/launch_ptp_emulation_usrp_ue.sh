echo "Bringing up UE side of the 5G system" &&
sudo docker compose --profile ue up # on the UE PC

echo "Waiting 15s for data session establishment and gtp tunnel creation" &&
sleep 15 &&

echo "Launching ptp containers" &&
docker compose up -d ptp-client &&

echo "Setting up routing for PTP packets" &&
# For network graph see docs/structure.drawio
# This can only run after the data session has been established by the 5gs (oaitun_ue1 created)

# UPF and UE use tunnel interfaces to communicate
docker exec oai-nr-ue ip r add 10.100.200.137 dev oaitun_ue1