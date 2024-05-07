# Resets the containers to the state of the last compiled images 
docker stop $(docker ps -a -q)
docker container prune -f
docker network prune -f