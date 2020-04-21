package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	inter "github.com/JointFaaS/Storage-Center/inter"
	pb "github.com/JointFaaS/Storage-Center/status"
	mapset "github.com/deckarep/golang-set"
	"google.golang.org/grpc"
	retry "gopkg.in/matryer/try.v1"
)

// SyncRPCServer implement method
type SyncRPCServer struct {
	storage inter.Storage
}

// Sync data to other client
func (s *SyncRPCServer) Sync(ctx context.Context, in *pb.SyncRequest) (*pb.SyncReply, error) {
	// TODO check state first
	value, err := s.storage.Get(in.Token)
	if err != nil {
		return &pb.SyncReply{Value: "", Code: -1}, err
	}
	return &pb.SyncReply{Value: value, Code: 1}, nil
}

type ClientImpl struct {
	name       string
	clientHost string
	clientPort string
	serverHost string
	state      inter.ClientState
	storage    inter.Storage
	client     pb.MaintainerClient
}

// NewClientImpl is a tool to get a client
func NewClientImpl(name string, clientHost string, clientPort string, serverHost string, c pb.MaintainerClient) ClientImpl {
	return ClientImpl{
		name:       name,
		clientPort: clientPort,
		clientHost: clientHost,
		serverHost: serverHost,
		state: &ClientStateImpl{
			holds: mapset.NewSet(),
		},
		storage: &StorageImpl{
			storage: make(map[string]string),
		},
		client: c,
	}
}

func (c *ClientImpl) Start(ctx context.Context, wg *sync.WaitGroup) error {
	// dail server
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		lis, err := net.Listen("tcp", c.clientPort)

		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			panic(err)
		}

		syncServer := grpc.NewServer()
		wg.Add(1)
		pb.RegisterSyncServer(syncServer, &SyncRPCServer{storage: c.storage})
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				{
					syncServer.GracefulStop()
					return
				}
			}
		}()
		log.Printf("%v sync server start at %v", c.name, c.clientPort)
		syncServer.Serve(lis)
	}(wg)

	if c.client == nil {
		conn, err := grpc.Dial(c.serverHost, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("can not connect with server %v", err)
			return err
		}
		// create stream
		c.client = pb.NewMaintainerClient(conn)
	}
	// retry-library refactor
	err := retry.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		resp, err := c.client.Register(ctx, &pb.RegisterRequest{Name: c.name, Host: c.clientHost + c.clientPort})
		if err != nil {
			log.Printf("open stream error %v", err)
		}
		if resp.Code > 0 {
			log.Printf("register %s over", c.name)
		} else {
			log.Printf("register error : %s\n", resp.Msg)
			err = errors.New("register error")
		}
		return
	})
	if err != nil {
		log.Printf("register error, closed")
		return errors.New("connection refused")
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var invalidStream pb.Maintainer_InvalidClient
		var err error
		err = retry.Do(func(attempt int) (retry bool, err error) {
			retry = attempt < 5
			invalidStream, err = c.client.Invalid(ctx)
			return
		})

		if err != nil {
			log.Println("error in stream open, close")
			return
		}

		if err = invalidStream.Send(&pb.InvalidRequest{Name: c.name}); err != nil {
			log.Fatalf("can not send %v", err)
		}

		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
					var token string
					err := retry.Do(func(attempt int) (retry bool, err error) {
						retry = attempt < 5
						resp, err := invalidStream.Recv()
						if err != nil {
							log.Printf("can not receive %v", err)
							return
						}
						token = resp.Token
						return
					})
					if token == "" {
						continue
					}
					if err != nil {
						log.Fatalf("can not get invalid message: %v\n", err)
					}
					// delete line by token
					c.state.Delete(token)
					err = c.storage.Delete(token)
					if err != nil {
						log.Fatalf("can not delete %v in storage\n", err)
					}
				}
			}
		}
	}(wg)
	return nil
}

// ChangeStatus will change the ownership of key
func (c *ClientImpl) ChangeStatus(ctx context.Context, token string) error {
	if c.client == nil {
		return errors.New("client not init, should call Start first")
	}

	// TODO check local first
	holded := c.state.Query(token)
	if holded {
		return nil
	}

	// apply for permission
	resp, err := c.client.ChangeStatus(ctx, &pb.StatusRequest{Token: token, Name: c.name})
	if err != nil {
		log.Fatalf("could not changeStatus: %v", err)
		return err
	}
	fmt.Printf("resp.Host %v, clientHost %v\n", resp.Host, c.clientHost)
	// TODO update state use resp.Token and resp.Host
	if resp.Host == c.clientHost {
		// we hold the status
		c.state.Add(resp.Token)
		hold := c.state.Query(resp.Token)
		fmt.Printf("after add into set %v\n", hold)
	}
	return nil
}

// Query first query from local, if it does not exist, it will sync from others
func (c *ClientImpl) Query(ctx context.Context, token string) (string, error) {
	// query local first holded means can read/write
	holded := c.state.Query(token)
	fmt.Printf("in query holded %v", holded)
	if holded {
		return c.storage.Get(token)
	}
	// local does not contains token value
	if c.client == nil {
		return "", errors.New("client not init, should call Start first")
	}
	resp, err := c.client.Query(ctx, &pb.QueryRequest{Token: token})
	if err != nil {
		log.Fatalf("could not changeStatus: %v", err)
		return "", err
	}
	// TODO call other client to get the data
	conn, err := grpc.Dial(resp.Host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
		return resp.Host, err
	}
	// create stream
	syncClient := pb.NewSyncClient(conn)
	// ugly forever loop to retry for capable server start racing
	for {
		syncResp, err := syncClient.Sync(context.Background(), &pb.SyncRequest{Token: token})
		if err != nil {
			log.Printf("sync server error: %v", err)
			continue
		}
		return syncResp.Value, nil
	}
}

func (c *ClientImpl) Set(token string, value string) {
	c.storage.Set(token, value)
	log.Printf("%v has set %v:%v\n", c.name, token, value)
}
