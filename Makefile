.PHONY : client server clean

client :
	go build -o build/client client_app/client.go

server :
	go build -o build/server server_app/server.go

tester :
	go build -o build/tester benchmark/main.go

clean:
	rm -f build/*

all: client server