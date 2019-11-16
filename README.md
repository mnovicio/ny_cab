# Database Setup
Database setup involves creating docker containers for
* MySQL server
* Database management tool [optional]

And creation of DB to be used and importing data set

## Docker Setup
### Pre-requisites: 
* **Docker** - v19.03.5 or later https://docs.docker.com/install/
* **Docker Compose** - v1.24.1 or later https://docs.docker.com/compose/install/ 

MySQL and Adminer services defined in stack.yml. Initiale and start services using docker-compose:
```
cd setup
docker-compose -f stack.yml up
```
This will create and start MySQL Server named 'mysql_server' and Adminer named 'mysql_adminer'

Verify containers are running:
```
docker ps -a
```

## Importing Dataset
### Recreate 'ny_cab_data' Database
```
docker exec -i mysql_server mysql -v -uroot -padmin123 < recreate_ny_cab_data_db.sql
```
### Import Dataset
```
docker exec -i mysql_server mysql -v -uroot -padmin123  < ny_cab_data_cab_trip_data_full.sql
```
**Note:** If you've been thinking about that coffee, this will be the great time to prepare that as this **import step may take a while**.

When that is done, verify that you have the 'ny_cab_data' database created with the 'cab_trip_data' table imported:
```
docker exec -i mysql_server mysql -v -uroot -padmin123 -e "select count(*) from ny_cab_data.cab_trip_data;"
docker exec -i mysql_server mysql -v -uroot -padmin123 -e "select * from ny_cab_data.cab_trip_data limit 10;"
```
