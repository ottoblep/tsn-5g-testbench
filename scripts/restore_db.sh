# Restores the database for the mongodb container from /config/database_backup.tar.gz
SCRIPTPATH=`dirname $0`
docker compose up db -d
docker cp $SCRIPTPATH/../config/database_backup.tar.gz mongodb:/database_backup.tar.gz
docker exec mongodb tar xvzf /database_backup.tar.gz
docker exec mongodb mongorestore database_backup
docker exec mongodb rm -rf /database_backup database_backup.tar.gz
docker compose stop db