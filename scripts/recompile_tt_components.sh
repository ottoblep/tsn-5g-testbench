SCRIPTPATH=`dirname $0`

cd $SCRIPTPATH/..

docker stop $(docker ps -a -q)
docker container prune -f 
docker build --target oai-nr-ue-dstt --tag oai-nr-ue-dstt:latest --file openairinterface5g/docker/Dockerfile.nrUEdsTT.ubuntu20 . --no-cache
docker build -t free5gc-upf-nwtt -f ./docker/free5gc/upfnwtt/Dockerfile . --no-cache
