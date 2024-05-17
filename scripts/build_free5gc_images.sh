SCRIPTPATH=`dirname $0`

# We build with the project root as context to have access to the submodules
# Docker builds can not copy files from a parent directory
cd $SCRIPTPATH/..

docker build -t free5gc-base --file ./docker/free5gc/base/Dockerfile .
docker build -t free5gc-upf --file ./docker/free5gc/upf/Dockerfile .
docker build -t free5gc-upf-nwtt --file ./docker/free5gc/upfnwtt/Dockerfile .