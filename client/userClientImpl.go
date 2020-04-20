package client

import (
	"context"
	"sync"

	inter "github.com/JointFaaS/Storage-Center/inter"
	pb "github.com/JointFaaS/Storage-Center/status"
)

type UserClientImpl struct {
	statusClient inter.StatusClient
	rootContext  context.Context
	finalizeFunc context.CancelFunc
	wg           sync.WaitGroup
}

func NewUserClientImpl(name string, clientHost string, clientPort string, serverHost string, c pb.MaintainerClient) *UserClientImpl {
	statusClient := NewClientImpl(name, clientHost, clientPort, serverHost, c)
	return &UserClientImpl{
		statusClient: &statusClient,
	}
}

func (c *UserClientImpl) Start() error {
	c.rootContext, c.finalizeFunc = context.WithCancel(context.Background())
	err := c.statusClient.Start(c.rootContext, &c.wg)
	if err != nil {
		return err
	}
	return nil
}

func (c *UserClientImpl) Close() {
	c.finalizeFunc()
	c.wg.Wait()
}

func (c *UserClientImpl) Get(token string) (string, error) {
	return c.statusClient.Query(token)
}

func (c *UserClientImpl) Set(token string, value string) error {
	err := c.statusClient.ChangeStatus(token)
	if err != nil {
		return err
	}
	c.statusClient.Set(token, value)
	return nil
}
