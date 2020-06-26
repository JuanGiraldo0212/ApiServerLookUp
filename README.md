# ApiServerLookUp
This API works as a endpoint to get SSL and servers information about an specific web domain. We provide two endpoints to do requests, the first one is /servers/[domain] to get the information, and /archive to get the previous domains that have been searched.
In order to run the server properly you will need initialize a CockroachDB node to handle the server requests. If you want to replicate this database you will have to create two tables, "items" and "servers" with the following columns respectively.

items(id STRING,grade STRING, time INT)

servers(id STRING, address STRING, grade STRING, country STRING, owner STRING, item STRING), the column item in the servers table is a foreign key that links a server to its item.

After setting up the DB you must change the data base path in the database.go file.
