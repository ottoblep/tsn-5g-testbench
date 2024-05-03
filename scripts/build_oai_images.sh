SCRIPTPATH=`dirname $0`

cd $SCRIPTPATH/../openairinterface5g

docker build --target ran-base --tag ran-base:latest --file docker/Dockerfile.base.ubuntu20 .
docker build --target ran-build --tag ran-build:latest --file docker/Dockerfile.build.ubuntu20 .
docker build --target oai-nr-ue --tag oai-nr-ue:develop --file docker/Dockerfile.nrUE.ubuntu20 .