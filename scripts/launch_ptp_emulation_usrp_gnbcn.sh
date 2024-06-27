echo "Launching software ptp server" &&
docker compose up -d ptp-server &&

echo "Setting up routing for PTP packets" &&
# For network graph see docs/structure.drawio
# This can only run after the data session has been established by the 5gs (oaitun_ue1 created)

# UPF and UE use tunnel interfaces to communicate
docker exec upf ip r add 10.60.0.1 dev upfgtp
