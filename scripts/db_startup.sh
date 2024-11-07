docker cp ./db/migrations/init.cql cassandra1:/init.cql
docker exec -it cassandra1 cqlsh -f /init.cql