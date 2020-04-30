# !/bin/bash

docker-compose up -d

docker-compose exec client1 tc qdisc add dev eth0 root netem delay 30ms
docker-compose exec client2 tc qdisc add dev eth0 root netem delay 30ms
docker-compose exec server tc qdisc add dev eth0 root netem delay 30ms

docker run --rm -i 