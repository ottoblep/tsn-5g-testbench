# Recompiles the binary in the upf
docker compose up -d free5gc-upf
docker exec upf make upf
docker stop $(docker ps -a -q)