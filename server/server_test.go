package server

import (
	"log"
	"net"
	"testing"
	"time"

	client "github.com/JointFaaS/Storage-Center/client"
	inter "github.com/JointFaaS/Storage-Center/inter"
	"github.com/JointFaaS/Storage-Center/state"
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
	defer s.GracefulStop()
	go func() {
		log.Println("rpc server start")
		s.Serve(lis)
	}()
	var c inter.UserClient
	c = client.NewUserClientImpl("test0", "127.0.0.1", ":50001", "127.0.0.1:50000", nil)
	defer c.Close()
	err = c.Start()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_Single_Client(t *testing.T) {
	key := "casecloud"
	storageValue := "nmsl"
	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMaintainerServer(s, &RPCServer{
		state: state.NewState(),
		hosts: state.NewHost(),
	})
	defer s.GracefulStop()
	go func() {
		log.Println("rpc server start")
		s.Serve(lis)
	}()

	time.Sleep(time.Duration(2) * time.Second)

	var c inter.UserClient
	c = client.NewUserClientImpl("test0", "127.0.0.1", ":50001", "127.0.0.1:50000", nil)
	defer c.Close()
	err = c.Start()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = c.Set(key, storageValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	val, err := c.Get(key)
	if err != nil {
		t.Errorf(err.Error())
	}

	if val != storageValue {
		t.Errorf("value %v should be %v\n", val, storageValue)
	}

}

func Test_Two_Client_With_One_Set_And_One_Get(t *testing.T) {
	key := "casecloud"
	storageValue := "nmsl"
	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMaintainerServer(s, &RPCServer{
		state: state.NewState(),
		hosts: state.NewHost(),
	})
	defer s.GracefulStop()
	go func() {
		log.Println("rpc server start")
		s.Serve(lis)
	}()

	time.Sleep(time.Duration(2) * time.Second)

	var c1 inter.UserClient
	c1 = client.NewUserClientImpl("test0", "127.0.0.1", ":50001", "127.0.0.1:50000", nil)
	defer c1.Close()
	err = c1.Start()
	if err != nil {
		t.Errorf(err.Error())
	}

	var c2 inter.UserClient
	c2 = client.NewUserClientImpl("test1", "127.0.0.1", ":50002", "127.0.0.1:50000", nil)
	defer c2.Close()
	err = c2.Start()
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c1.Set(key, storageValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	val, err := c2.Get(key)
	if err != nil {
		t.Errorf(err.Error())
	}

	if val != storageValue {
		t.Errorf("value %v should be %v\n", val, storageValue)
	}

}

func Test_Two_Client_Double_Set(t *testing.T) {
	key := "casecloud"
	storageValue := "nmsl"
	reverseValue := "lsmn"
	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMaintainerServer(s, &RPCServer{
		state: state.NewState(),
		hosts: state.NewHost(),
	})
	defer s.GracefulStop()
	go func() {
		log.Println("rpc server start")
		s.Serve(lis)
	}()

	time.Sleep(time.Duration(2) * time.Second)

	var c1 inter.UserClient
	c1 = client.NewUserClientImpl("test0", "127.0.0.1", ":50001", "127.0.0.1:50000", nil)
	defer c1.Close()
	err = c1.Start()
	if err != nil {
		t.Errorf(err.Error())
	}

	var c2 inter.UserClient
	c2 = client.NewUserClientImpl("test1", "127.0.0.1", ":50002", "127.0.0.1:50000", nil)
	defer c2.Close()
	err = c2.Start()
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c1.Set(key, storageValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c2.Set(key, reverseValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	val, err := c1.Get(key)

	if err != nil {
		t.Errorf(err.Error())
	}

	if val != reverseValue {
		t.Errorf("value %v should be %v\n", val, reverseValue)
	}

	val, err = c2.Get(key)

	if err != nil {
		t.Errorf(err.Error())
	}

	if val != reverseValue {
		t.Errorf("value %v should be %v\n", val, reverseValue)
	}

}
