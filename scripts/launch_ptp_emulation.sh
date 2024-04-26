echo "Bringing up 5g system" &&
docker compose --profile 5gs up -d &&

echo "Waiting 20s for data session establishment and gtp tunnel creation" &&
sleep 20 &&

echo "Launching ptp containers" &&
docker compose --profile ptpsim up -d &&

echo "Setting up routing for PTP packets" &&
$(dirname "$0")/setup_ptp_routing.sh &&

echo "Running ptp slave" &&
docker exec -d ptp-slave /bestvibes-ptp/slave/slave.out &&
sleep 1 &&

echo "Running ptp master" &&
docker exec ptp-master /bestvibes-ptp/master/master.out