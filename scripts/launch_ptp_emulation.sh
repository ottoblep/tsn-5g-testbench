echo "Bringing up 5g system" &&
docker compose --profile cn --profile ran-rfsim up -d &&

echo "Waiting 15s for data session establishment and gtp tunnel creation" &&
sleep 15 &&

echo "Launching ptp containers" &&
docker compose --profile ptpsim up -d &&

echo "Setting up routing for PTP packets" &&
$(dirname "$0")/setup_ptp_routing.sh
