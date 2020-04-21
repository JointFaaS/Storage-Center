package client

import (
	"context"
	"sync"

	inter "github.com/JointFaaS/Storage-Center/inter"
	pb "github.com/JointFaaS/Storage-Center/status"
)

// UserClientImpl use context and RPC to implement UserClient
type UserClientImpl struct {
	statusClient inter.StatusClient
	rootContext  context.Context
	finalizeFunc context.CancelFunc
	wg           sync.WaitGroup
}

// NewUserClientImpl is an support function for other package
func NewUserClientImpl(name string, clientHost string, clientPort string, serverHost string, c pb.MaintainerClient) *UserClientImpl {
	statusClient := NewClientImpl(name, clientHost, clientPort, serverHost, c)
	return &UserClientImpl{
		statusClient: &statusClient,
		wg:           sync.WaitGroup{},
	}
}

// Start will init context and start statusClient
func (c *UserClientImpl) Start() error {
	c.rootContext, c.finalizeFunc = context.WithCancel(context.Background())
	err := c.statusClient.Start(c.rootContext, &c.wg)
	if err != nil {
		return err
	}
	return nil
}

// Close use cancel function and wait all routine over
func (c *UserClientImpl) Close() {
	c.finalizeFunc()

	c.wg.Wait()
}

// Get value by key
func (c *UserClientImpl) Get(token string) (string, error) {
	return c.statusClient.Query(c.rootContext, token)
}

// Set value
func (c *UserClientImpl) Set(token string, value string) error {
	err := c.statusClient.ChangeStatus(c.rootContext, token)
	if err != nil {
		return err
	}
	c.statusClient.Set(token, value)
	return nil
}
