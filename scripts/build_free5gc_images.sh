SCRIPTPATH=`dirname $0`

# We build with the project root as context to have access to the submodules
# Docker builds can not copy files from a parent directory
cd $SCRIPTPATH/..

docker build -t free5gc/base:latest -f ./docker/free5gc/base/Dockerfile .