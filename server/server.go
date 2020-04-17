package main

import (
    "log"
    "net"
	
	pb "github.com/JointFaaS/Storage-Center/status"
	state "github.com/JointFaaS/Storage-Center/state"
	inter "github.com/JointFaaS/Storage-Center/inter"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

const (
    // PORT defined for listen port
    PORT = ":50001"
)

type server struct {
	state inter.State
	hosts inter.Host
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	log.Println("request: Name", in.Name)
	log.Println("request: Host", in.Host)
	err := s.hosts.Insert(in.Host, in.Name)
	if (err != nil) {
		return &pb.RegisterReply{Code: -1, Msg: err.Error()}, err
	}
	return &pb.RegisterReply{Code: 1, Msg: "OK"}, nil
}

func (s *server) ChangeStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusReply, error) {
	_, oldname, err := s.state.ChangeStatus(in.Token, in.Name)
	if (err != nil) {
		panic(err)
	}
	host, err := s.hosts.Query(oldname) 
	if (err != nil) {
		panic(err)
	}

	// TODO announce to host => invalid

    return &pb.StatusReply{Token: in.Token, Host: host}, nil
}

func (s *server) Query(ctx context.Context, in *pb.QueryRequest) (*pb.QueryReply, error) {
	name, err := s.state.Query(in.Token);
	if (err != nil) {
		panic(err)
	}
	host, err := s.hosts.Query(name);
	if (err != nil) {
		panic(err)
	}
    return &pb.QueryReply{Token: in.Token, Host: host}, nil
}

func main() {
    lis, err := net.Listen("tcp", PORT)

    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterMaintainerServer(s, &server{
		state: &state.StateImpl{},
		hosts: &state.HostImpl{},
	})
    log.Println("rpc服务已经开启")
    s.Serve(lis)
}