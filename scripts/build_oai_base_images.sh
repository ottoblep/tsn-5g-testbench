SCRIPTPATH=`dirname $0`

cd $SCRIPTPATH/../openairinterface5g

# Patches for OAI
# - Sending the 5g mobility management capability field during UE registration, which free5gc requires
# - A dirty fix for the UE to skip unknown fields in the PDU establishment accept message, otherwise the UE aborts parsing the message
# - Add iptables to UE for routing purposes
# - Speed up the OAI build process by removing 4G components
git apply ../patches/openairinterface5g/*.patch

docker build --target ran-base --tag ran-base:latest --file docker/Dockerfile.base.ubuntu20 .
docker build --target ran-build --tag ran-build:latest --file docker/Dockerfile.build.ubuntu20 .
docker build --target oai-nr-ue --tag oai-nr-ue:develop --file docker/Dockerfile.nrUE.ubuntu20 .