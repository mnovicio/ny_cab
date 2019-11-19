
Cab Data Researcher is a company that provides insights on the open data about NY cab trips

Cab trips in NY are publicly available as CSV downloadable files. In order to make it more useful, we want to wrap the data in a public API.

Backend is implemented in Go using grpc and grpc-gateway to provide REST endpoints.

Frontend CLI is also implemented in GO using Cobra.

Enjoy!

---------------------------------------
* [Database Setup](#database-setup)
  * [Docker Setup](#docker-setup)
    * [Prerequisite](#prerequisites)
  * [Importing Dataset](#importing-dataset)
    * [Recreate DB](#recreate-'ny_cab_data'-database)
    * [Import Dataset](#import-dataset)
* [Backend Service](#backend-service)
  * [Prerequisite](#prerequisites-1)
  * [Build backend service](#build-backend-service)
  * [Run service](#run-service)
  * [Protobuf GO code generation](#protobuf-go-code-generation)
  * [Backend REST Enpoints](#backend-rest-enpoints)
    * [/v1/cabtrips](#/v1/cabtrips)
    * [/v1/cabtrips/bypickupdate](#/v1/cabtrips/bypickupdate)
    * [/v1/cabtrips/clearcache](#/v1/cabtrips/clearcache)
* [Command Line Client - REST](#command-line-client---rest)
  * [Build](#build)
  * [Usage](#usage)
     * [Show help](#show-help)
     * [Show command help](#show-command-help)
* [Command Line Client - GRPC](#command-line-client---grpc)
  * [Build](#build-1)
  * [Usage](#usage-1)
    * [Show help](#show-help-1)
    * [Show command help](#show-command-help-1)
---------------------------------------

# Database Setup
Database setup involves creating docker containers for
* MySQL server
* Database management tool [optional]

And creation of DB to be used and importing data set

## Docker Setup
### Prerequisites: 
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
# Backend Service

The backend is implemented using gRPC with gRPC-gateway to provide both gRPC interface as well as HTTP REST endpoints.
For more information, https://github.com/grpc-ecosystem/grpc-gateway#grpc-gateway.

* gRPC server is hosted at port 10001
* HTTP REST gateway is host at port 10002

## Prerequisites
* **Go** v1.13 or later https://golang.org/doc/install
* **GNU Make** (optional)
## Build backend service
Using Make

**Note:** This generates the 'ny_cab_server' binary inside src/mnovicio.com/nycab/bin folder
```
cd src/mnovicio.com/nycab
make server
```

Using manual step
```
cd src/mnovicio.com/nycab
mkdir -p ./bin
cd server/cmd && go build -o ../../bin/ny_cab_server -v
```
## Run service
```
cd src/mnovicio.com/nycab/bin
./ny_cab_server
```

## Protobuf GO code generation
The server/client GO code has already been generated from corresponding proto files inside:
* src/mnovicio.com/protocol/objects
* src/mnovicio.com/protocol/rpc

For more information about protobuf: https://developers.google.com/protocol-buffers/docs/gotutorial


## Backend REST endpoints
Host: http://localhost:10002
### **/v1/cabtrips**

    Method: POST
    Description: Returns all cab trips per day on record
    Body Content type: application/json
    Body (example):
    {
        ignore_cache: true
    }
    Parameters:
        ignore_cache: true - ignores cached data and fetch fresh data from DB, false - use cached data
    Returns:
    

### **/v1/cabtrips/bypickupdate**

    Method: POST
    Description: Returns number of trips a particular cab has made given a particular pickup date, time ignored
    Body Content type: application/json
    Body (example):
    {
        "cab_ids": [
            "D7D598CD99978BD012A87A76A7C891B7",
            "42D815590CE3A33F3A23DBF145EE66E3",
            "NONEXISTENTMEDALION"
            ],
        "pickup_date": "2013-12-01",
        "ignore_cache": false
    }    
    Parameters:
        cab_ids: list of cab IDs to fetch
        pickup_date: specified pickup date
        ignore_cache:
            true - ignores cached data and fetch fresh data from DB
            false - use cached data if available, fetches the DB for any cab ID with pickup date not found in cache
    Returns (example):
    {
        "cab_trips_per_day": {
            "cab_trips": {
                "42D815590CE3A33F3A23DBF145EE66E3": {
                    "trips_per_day": {
                        "2013-12-01": 1
                    }
                },
                "D7D598CD99978BD012A87A76A7C891B7": {
                    "trips_per_day": {
                        "2013-12-01": 3
                    }
                },
                "NONEXISTENTMEDALION": {
                    "trips_per_day": {
                        "2013-12-01": 0
                    }
                }
            }
        }
    }

### **/v1/cabtrips/clearcache**

    Method: GET
    Description: Clears all cached data


# Command Line Client - REST
## Build
Using Make
```
cd src/mnovicio.com/nycab
make client/rest
```

Using manual step
```
cd src/mnovicio.com/nycab
mkdir -p ./bin
cd client/rest && go build -o ../../bin/ny_cab_client_rest -v
```

## Usage
### **show help**
$ ./ny_cab_client_rest -h
```
$ ./ny_cab_client_rest -h
connects to localhost:10002 to use NY CAB REST endpoints

Usage:
  ny_cab_client_rest [command]

Available Commands:
  clear-cache             Clears cached data on the server
  get-all-cab-trip-count  Prints all cab trips on record
  get-trip-counts-for-cab Prints cab trip count on given pickup date
  help                    Help about any command

Flags:
  -h, --help            help for ny_cab_client_rest
  -s, --server string   NY CAB service host (default "http://localhost:10002")

Use "ny_cab_client_rest [command] --help" for more information about a command.
```
### **show command help**
$ ./ny_cab_client_rest [command] -h
```
$ ./ny_cab_client_rest get-trip-counts-for-cab -h
Prints cab trip count on given pickup date
Example: ./ny_cab_client_rest get-trip-counts-for-cab --cab-ids="cab1,cab2" --pickup-date="2013-12-01" --ignore-cache=true

Usage:
  ny_cab_client_rest get-trip-counts-for-cab [flags]

Flags:
      --cab-ids strings      list of cab IDs to fetch (default [D7D598CD99978BD012A87A76A7C891B7,42D815590CE3A33F3A23DBF145EE66E3])
  -h, --help                 help for get-trip-counts-for-cab
      --ignore-cache         Ignore cached data and force fetch DB
      --pickup-date string   pickup date (default "2013-12-01")

Global Flags:
  -s, --server string   NY CAB service host (default "http://localhost:10002")
```

# Command Line Client - GRPC
## Build
Using Make
```
cd src/mnovicio.com/nycab
make client/grpc
```

Using manual step
```
cd src/mnovicio.com/nycab
mkdir -p ./bin
cd client/grpc && go build -o ../../bin/ny_cab_client_grpc -v
```

## Usage
### **show help**
$ ./ny_cab_client_grpc -h
```
$ ./ny_cab_client_grpc -h
connects to localhost:10001 to use NY CAB gRPC APIs

Usage:
  ny_cab_client_grpc [command]

Available Commands:
  clear-cache             Clears cached data on the server
  get-all-cab-trip-count  Prints all cab trips on record
  get-trip-counts-for-cab Prints cab trip count on given pickup date
  help                    Help about any command

Flags:
  -h, --help            help for ny_cab_client_grpc
  -s, --server string   NY CAB gRPC host (default "localhost:10001")

Use "ny_cab_client_grpc [command] --help" for more information about a command.
```
### **show command help**
$ ./ny_cab_client_grpc [command] -h
```
$ ./ny_cab_client_grpc get-trip-counts-for-cab -h
Prints cab trip count on given pickup date
Example: ./ny_cab_client_grpc get-trip-counts-for-cab --cab-ids="cab1,cab2" --pickup-date="2013-12-01" --ignore-cache=true

Usage:
  ny_cab_client_grpc get-trip-counts-for-cab [flags]

Flags:
      --cab-ids strings      list of cab IDs to fetch (default [D7D598CD99978BD012A87A76A7C891B7,42D815590CE3A33F3A23DBF145EE66E3])
  -h, --help                 help for get-trip-counts-for-cab
      --ignore-cache         Ignore cached data and force fetch DB
      --pickup-date string   pickup date (default "2013-12-01")

Global Flags:
  -s, --server string   NY CAB gRPC host (default "localhost:10001")
```