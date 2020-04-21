package server

import (
	"io"
	"log"

	inter "github.com/JointFaaS/Storage-Center/inter"
	"github.com/JointFaaS/Storage-Center/state"
	pb "github.com/JointFaaS/Storage-Center/status"
	"golang.org/x/net/context"
)

// RPCServer implement Tx-Server
type RPCServer struct {
	state inter.State
	hosts inter.Host
}

// NewRPCServer will return pointer to RPCServer
func NewRPCServer() *RPCServer {
	return &RPCServer{
		state: state.NewState(),
		hosts: state.NewHost(),
	}
}

// Register will store name and host
func (s *RPCServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	log.Println("request: Name", in.Name)
	log.Println("request: Host", in.Host)
	err := s.hosts.Insert(in.Host, in.Name)
	if err != nil {
		return &pb.RegisterReply{Code: -1, Msg: err.Error()}, nil
	}
	return &pb.RegisterReply{Code: 1, Msg: "OK"}, nil
}

// ChangeStatus will change ownership
func (s *RPCServer) ChangeStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusReply, error) {
	for {

		newname, oldname, err := s.state.ChangeStatus(in.Token, in.Name)
		if err != nil {
			panic(err)
		}
		host, err := s.hosts.Query(newname)
		if err != nil {
			panic(err)
		}
		if oldname != "" {
			channel, err := s.hosts.GetChan(oldname)
			channel <- in.Token
			if err != nil {
				log.Printf("GetChan in ChaneStatus error %v", err)
				panic(err)
			}

		}
		resp := &pb.StatusReply{
			Token: in.Token,
			Host:  host,
		}
		return resp, nil
	}

	// TODO announce to host => invalid
}

// Query state and storage
func (s *RPCServer) Query(ctx context.Context, in *pb.QueryRequest) (*pb.QueryReply, error) {
	//receive data from stream
	name, err := s.state.Query(in.Token)
	if err != nil {
		return nil, err
	}
	host, err := s.hosts.Query(name)
	if err != nil {
		return nil, err
	}
	resp := &pb.QueryReply{Token: in.Token, Host: host}
	return resp, nil
}

// Invalid stream will return the invalid key for your onw host
func (s *RPCServer) Invalid(srv pb.Maintainer_InvalidServer) error {
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
		case token := <-invalidChannel:
			{
				resp := pb.InvalidReply{Token: token}
				if err := srv.Send(&resp); err != nil {
					log.Printf("Query send error %v", err)
				}
				break
			}
		}
	}
}
