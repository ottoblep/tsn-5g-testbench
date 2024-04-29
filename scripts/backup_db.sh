# Archives the database from the mongodb container and writes it to /config/database_backup.tar.gz
SCRIPTPATH=`dirname $0`
docker compose up db -d
docker exec mongodb mongodump -d free5gc -o database_backup
docker exec mongodb tar zcvf /database_backup.tar.gz /database_backup
docker cp mongodb:/database_backup.tar.gz $SCRIPTPATH/../config/database_backup.tar.gz
docker exec mongodb rm -rf /database_backup database_backup.tar.gz
docker compose stop db