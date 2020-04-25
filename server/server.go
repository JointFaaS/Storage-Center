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

		newname, oldname, newVersion, _, err := s.state.ChangeStatus(in.Token, in.Name)
		if err != nil {
			panic(err)
		}
		host, err := s.hosts.Query(newname)
		if err != nil {
			panic(err)
		}
		if oldname != "" {
			channel, err := s.hosts.GetChan(oldname)
			if err != nil {
				log.Printf("GetChan in ChaneStatus error %v", err)
				panic(err)
			}
			returnChan := make(chan int8)
			channel <- state.InvalidEntry{
				Token:   in.Token,
				Channel: returnChan,
			}
			invalidCode := <-returnChan
			if invalidCode < 0 {
				resp := &pb.StatusReply{
					Token:   in.Token,
					Host:    oldname,
					Version: 0,
				}
				return resp, nil
			}
		}
		resp := &pb.StatusReply{
			Token:   in.Token,
			Host:    host,
			Version: newVersion,
		}
		return resp, nil
	}
}

// Query state and storage
func (s *RPCServer) Query(ctx context.Context, in *pb.QueryRequest) (*pb.QueryReply, error) {
	//receive data from stream
	name, version, err := s.state.Query(in.Token)
	if err != nil {
		return nil, err
	}
	host, err := s.hosts.Query(name)
	if err != nil {
		return nil, err
	}

	resp := &pb.QueryReply{
		Token:   in.Token,
		Host:    host,
		Version: version,
	}
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
		case invalidEntry := <-invalidChannel:
			{
				resp := pb.InvalidReply{Token: invalidEntry.Token}
				if err := srv.Send(&resp); err != nil {
					log.Printf("Query send error %v", err)
					invalidEntry.Channel <- -1
				} else {
					// waiting for recv
					returnReq, err := srv.Recv()
					if returnReq.Token != invalidEntry.Token && err == nil {
						invalidEntry.Channel <- 1
					} else {
						invalidEntry.Channel <- -1
					}
				}
				break
			}
		}
	}
}
