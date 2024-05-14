SCRIPTPATH=`dirname $0`

cd $SCRIPTPATH/..

docker stop $(docker ps -a -q)
docker container prune -f 
docker build -t ptp-server ./docker/ptp-server --no-cache
docker build -t ptp-client ./docker/ptp-client --no-cache

