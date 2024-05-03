SCRIPTPATH=`dirname $0`

# We build with the project root as context to have access to the submodules
cd $SCRIPTPATH/..

docker build -t free5gc/base:latest -f ./docker/free5gc/base/Dockerfile .
docker build --build-arg F5GC_MODULE=upf -t free5gc/upf-base:latest -f ./docker/free5gc/base/Dockerfile.nf .