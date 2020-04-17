package main

import (
	"log"
	"io"
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
	if err != nil {
		return &pb.RegisterReply{Code: -1, Msg: err.Error()}, err
	}
	return &pb.RegisterReply{Code: 1, Msg: "OK"}, nil
}

func (s *server) ChangeStatus(srv pb.Maintainer_ChangeStatusServer) error {
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		//receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("ChangeStatus receive error %v", err)
			continue
		}
		newname, oldname, err := s.state.ChangeStatus(req.Token, req.Name)
		if err != nil {
			panic(err)
		}
		host, err := s.hosts.Query(newname)
		if err != nil {
			panic(err)
		}
		if oldname != "" {
			channel, err := s.hosts.GetChan(oldname)
			channel <- req.Token
			if err != nil {
				log.Printf("GetChan in ChaneStatus error %v", err)
				panic(err)
			}

		}
		resp := pb.StatusReply {
			Token: req.Token,
			Host: host,
		}
		if err := srv.Send(&resp); err != nil {
			log.Printf("ChangeStatus send error %v", err)
		}
	}

	// TODO announce to host => invalid
}

func (s *server) Query(srv pb.Maintainer_QueryServer)  error {
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		//receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Query receive error %v", err)
			continue
		}
		name, err := s.state.Query(req.Token);
		if err != nil {
			panic(err)
		}
		host, err := s.hosts.Query(name);
		if err != nil {
			panic(err)
		}
		resp := pb.QueryReply{Token: req.Token, Host: host}
		if err := srv.Send(&resp); err != nil {
			log.Printf("Query send error %v", err)
		}
	}
}

func (s *server) Invalid(srv pb.Maintainer_InvalidServer)  error {
	ctx := srv.Context()
	//receive data from stream
	req, err := srv.Recv()	
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Printf("Invalid receive error %v", err)
		panic(err)
	}
	invalidChannel, err := s.hosts.GetChan(req.Name)
	if err != nil {
		log.Printf("Invalid GetChan error %v", err)
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case token := <- invalidChannel: {	
			resp := pb.InvalidReply{Token: token}
			if err := srv.Send(&resp); err != nil {
				log.Printf("Query send error %v", err)
			}
			break
		}
		default:
		}
	}
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