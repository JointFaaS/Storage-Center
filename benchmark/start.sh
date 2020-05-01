# !/bin/bash

docker-compose down
docker-compose up -d server

echo 'wait 5s before serve ready'
sleep 5

docker-compose up -d

echo 'apply linux tc'

docker-compose exec client1 tc qdisc add dev eth0 root netem delay 30ms
docker-compose exec client2 tc qdisc add dev eth0 root netem delay 30ms
docker-compose exec server tc qdisc add dev eth0 root netem delay 30ms

docker-compose exec benchmark /go/src/app/build/tester random