# Each container network interface is accessible from the host for monitoring.
# The interface names on the host are structured as vethXXXXXXX
# This script identifies the host veth ID of the interface eth0 for each docker container
for container in $(docker ps -q); do
    iflink=`docker exec -it $container bash -c 'cat /sys/class/net/eth0/iflink'`
    iflink=`echo $iflink|tr -d '\r'`
    veth=`grep -l "$iflink" /sys/class/net/veth*/ifindex`
    veth=`echo $veth|sed -e 's;^.*net/\(.*\)/ifindex$;\1;'`
    echo $container:$veth
done