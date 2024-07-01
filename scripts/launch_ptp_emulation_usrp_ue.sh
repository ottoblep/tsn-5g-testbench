SCRIPTPATH=`dirname $0`

echo "Initializing ptpd sync via ethernet. This is needed to timestamp consistently inside the DS-TT and NW-TT. Normally the clocks of gNB and UE will be synced via the 5GS RRC but this is not implemented yet."
# This can also be tested manually by adding --foreground
sudo ptpd -c $SCRIPTPATH/../config/eth_ptpd_client.config

echo "Launching software ptp client"
docker compose up -d ptp-client

echo "Setting up routing for PTP packets"
# For network graph see docs/structure.drawio
# This can only run after the data session has been established by the 5gs (oaitun_ue1 created)

# UPF and UE use tunnel interfaces to communicate
docker exec oai-nr-ue ip r add 10.100.200.137 dev oaitun_ue1