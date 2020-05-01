# !/bin/bash

docker-compose down
docker-compose up -d

echo 'wait 5s before testing'
sleep 5

echo 'apply linux tc'

docker-compose exec client1 tc qdisc add dev eth0 root netem delay 30ms
docker-compose exec client2 tc qdisc add dev eth0 root netem delay 30ms
docker-compose exec server tc qdisc add dev eth0 root netem delay 30ms

docker-compose exec benchmark /go/src/app/build/tester random