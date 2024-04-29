SCRIPTPATH=`dirname $0`

cd $SCRIPTPATH/../free5gc-compose

# Patches for free5gc-compose
# - Use submodule as source 
git apply ../patches/free5gc-compose/*.patch
cp -rf ../go-upf ./base # Docker copy only accepts files below it 

docker build -t free5gc/base:latest -f ./base
docker build --build-arg F5GC_MODULE=upf -t free5gc/upf-base:latest -f ./base/Dockerfile.nf ./base