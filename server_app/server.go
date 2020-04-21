package main

import (
	"flag"
	"log"
	"net"

	server "github.com/JointFaaS/Storage-Center/server"
	pb "github.com/JointFaaS/Storage-Center/status"
	"google.golang.org/grpc"
)

var (
	help bool
	port string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&port, "p", ":50000", "set server port")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMaintainerServer(s, server.NewRPCServer())
	log.Println("rpc server start")
	s.Serve(lis)
}
