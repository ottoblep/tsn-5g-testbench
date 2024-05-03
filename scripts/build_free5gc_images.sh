SCRIPTPATH=`dirname $0`

# We build with the project root as context to have access to the submodules
cd $SCRIPTPATH/..

# Patches for free5gc-compose
# - Use submodule as source 
git apply ./patches/free5gc-compose/*.patch --directory free5gc-compose

docker build -t free5gc/base:latest -f ./free5gc-compose/base .
docker build --build-arg F5GC_MODULE=upf -t free5gc/upf-base:latest -f ./free5gc-compose/base/Dockerfile.nf .