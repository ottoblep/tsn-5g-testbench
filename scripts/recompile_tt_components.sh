SCRIPTPATH=`dirname $0`

cd $SCRIPTPATH/..

docker stop $(docker ps -a -q)
docker container prune -f 
# docker image rm 5g-tsn-testbed-free5gc-upf
docker build --target oai-nr-ue --tag oai-nr-ue:develop --file openairinterface5g/docker/Dockerfile.nrUE.ubuntu20 . --no-cache
docker compose up -d free5gc-upf
docker stop $(docker ps -a -q)
