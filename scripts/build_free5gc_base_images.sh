SCRIPTPATH=`dirname $0`

cd $SCRIPTPATH/../free5gc-compose

# Patches for free5gc-compose
# - Clone the free5gc source instead of using a local copy
git apply ../patches/free5gc-compose/*.patch

docker build -t free5gc/base:latest ./base
docker build --build-arg F5GC_MODULE=upf -t free5gc/upf-base:latest -f ./base/Dockerfile.nf ./base