package server

import (
	"log"
	"net"
	"testing"
	"time"

	client "github.com/JointFaaS/Storage-Center/client"
	inter "github.com/JointFaaS/Storage-Center/inter"
	state "github.com/JointFaaS/Storage-Center/state"
	pb "github.com/JointFaaS/Storage-Center/status"
	"google.golang.org/grpc"
)

func Test_Init(t *testing.T) {
	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMaintainerServer(s, &RPCServer{
		state: state.NewState(),
		hosts: state.NewHost(),
	})
	log.Println("rpc server start")
	defer s.GracefulStop()
	go func() {
		s.Serve(lis)
	}()

	time.Sleep(time.Duration(2) * time.Second)

	var c inter.UserClient
	c = client.NewUserClientImpl("test", "127.0.0.1", ":50001", "127.0.0.1:50000", nil)
	defer c.Close()
	err = c.Start()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_Single_Client(t *testing.T) {
	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMaintainerServer(s, &RPCServer{
		state: &state.StateImpl{},
		hosts: &state.HostImpl{},
	})
	defer s.GracefulStop()
	go func() {
		log.Println("rpc server start")
		s.Serve(lis)
	}()

	time.Sleep(time.Duration(2) * time.Second)

	var c inter.UserClient
	c = client.NewUserClientImpl("test", "127.0.0.1", ":50001", "127.0.0.1:50000", nil)
	defer c.Close()
	err = c.Start()
	if err != nil {
		t.Errorf(err.Error())
	}

}
