package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JointFaaS/Storage-Center/client"
	"github.com/JointFaaS/Storage-Center/inter"
)

type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func get(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	err := request.ParseForm()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	key := request.Form["key"]
	value, err := rpcClient.Get(key[0])
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	res := GetResponse{
		Key:   key[0],
		Value: value,
	}
	json.NewEncoder(response).Encode(res)
}

func set(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	err := request.ParseForm()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}
	var req SetRequest
	err = json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	err = rpcClient.Set(req.Key, req.Value)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	res := SetResponse{
		Key:   req.Key,
		Value: req.Value,
	}
	json.NewEncoder(response).Encode(res)
}

var (
	help      bool
	host      string
	name      string
	port      string
	syncPort  string
	server    string
	rpcClient inter.UserClient
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&host, "H", "127.0.0.1", "set own host")
	flag.StringVar(&name, "n", "test", "set own name")
	flag.StringVar(&syncPort, "sp", ":55389", "set sync port")
	flag.StringVar(&port, "p", ":9091", "set server port")
	flag.StringVar(&server, "r", "127.0.0.1:50000", "set remote master addr")
}

func main() {
	flag.Parse()
	if help {
		fmt.Fprintf(os.Stderr, `Tx-Server version: 1.0.0
Usage: server [-h] [-H host] [-n name] [-sp sync port] [-p port] [-r remote addr]

Options:
`)
		flag.Usage()
		return
	}
	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
	rpcClient = client.NewUserClientImpl(name, host, syncPort, server, nil)
	err := rpcClient.Start()
	if err != nil {
		log.Fatal("RpcClient error", err)
	}
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
