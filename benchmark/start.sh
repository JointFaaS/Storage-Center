# !/bin/bash

docker-compose down
docker-compose up -d server

echo 'wait 5s before server ready'
sleep 5

docker-compose up -d

echo 'apply linux tc'

# export client1=`docker-compose exec client1 tail -n 1 /etc/hosts | sed 's/\t/\n/g' | head -n 1`
# export client2=`docker-compose exec client2 tail -n 1 /etc/hosts | sed 's/\t/\n/g' | head -n 1`
# export server=`docker-compose exec server tail -n 1 /etc/hosts | sed 's/\t/\n/g' | head -n 1`

# docker-compose exec client1 tc qdisc add dev eth0 root handle 1: prio bands 4
# docker-compose exec client1 tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 30ms
# docker-compose exec client1 tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $client2 flowid 1:4

# docker-compose exec client2 tc qdisc add dev eth0 root handle 1: prio bands 4
# docker-compose exec client2 tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 30ms
# docker-compose exec client2 tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $client1 flowid 1:4

docker-compose exec server tc qdisc add dev eth0 root netem delay 30ms

echo 'wait 5s before client ready'
sleep 5

echo 'random test'

docker-compose exec benchmark /go/src/app/build/tester random >> result.txt
