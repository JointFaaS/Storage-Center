package client
import (
	"context"
	"io"
	"log"
	"errors"
	"google.golang.org/grpc"
	pb "github.com/JointFaaS/Storage-Center/status"
	inter "github.com/JointFaaS/Storage-Center/inter"
	mapset "github.com/deckarep/golang-set"
)
type ClientImpl struct {
	name string
	clientHost string
	serverHost string
	state inter.ClientState
	storage inter.Storage
	client pb.MaintainerClient
}

// NewClientImpl is a tool to get a client
func NewClientImpl(name string, clientHost string, serverHost string, c pb.MaintainerClient) (ClientImpl) {
	return ClientImpl{
		name: name,
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

func (c *ClientImpl) Start() error {
	// dail server
	if (c.client == nil) {
		conn, err := grpc.Dial(c.serverHost, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("can not connect with server %v", err)
			return err
		}
		// create stream
		c.client = pb.NewMaintainerClient(conn)
	}
	resp, err := c.client.Register(context.Background(), &pb.RegisterRequest{Name: c.name, Host: c.clientHost})
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}
	if resp.Code > 0 {
		log.Printf("register %s over", c.name)
	} else {
		log.Printf("register error : %s\n", resp.Msg)
		return errors.New("register error");
	}

	
	go func() {
		invalidStream, err := c.client.Invalid(context.Background())
		ctx := invalidStream.Context()
		done := make(chan bool)
		if err != nil {
			log.Fatalf("openn stream error %v", err)
		}
		
		if err := invalidStream.Send(&pb.InvalidRequest{Name: c.name}); err != nil {
			log.Fatalf("can not send %v", err)
		}

		// nested goroutine closes done channel
		// if context is done
		go func() {
			<-ctx.Done()
			if err := ctx.Err(); err != nil {
				log.Println(err)
			}
			close(done)
		}()

		for {
			resp, err := invalidStream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			token := resp.Token
			log.Printf("new invalid token %s received", token)
			// delete line by token
			c.state.Delete(token)
			err = c.storage.Delete(token)
			if err != nil {
				log.Fatalf("can not delete %v in storage", err)
			}
		}
	}()
	return nil
}

func (c *ClientImpl) ChangeStatus(token string) error {
	if c.client == nil {
		return errors.New("client not init, should call Start first")
	}

	// TODO check local first
	holded := c.state.Query(token)
	if (holded) {
		return nil
	}

	// apply for permission
	resp, err := c.client.ChangeStatus(context.Background(), &pb.StatusRequest{Token: token, Name: c.name})
	if err != nil {
		log.Fatalf("could not changeStatus: %v", err)
		return err
	}
	// TODO update state use resp.Token and resp.Host
	if (resp.Host == c.clientHost) {
		// we hold the 
	}
	return nil
}

func (c *ClientImpl) Query(token string) (string, error) {
	// query local first holded means can read/write
	holded := c.state.Query(token)
	if (holded) {
		return c.storage.Get(token), nil
	}
	// local does not contains token value
	if c.client == nil {
		return "", errors.New("client not init, should call Start first")
	}
	resp, err := c.client.Query(context.Background(), &pb.QueryRequest{Token: token})
	if err != nil {
		log.Fatalf("could not changeStatus: %v", err)
		return "", err
	}
	// TODO call other client to get the data
	return resp.Host, nil
}

func (c *ClientImpl) Set(token string, value string) {
	c.storage.Set(token, value)
}
