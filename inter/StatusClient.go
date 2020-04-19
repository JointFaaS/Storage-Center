package inter

import (
	"sync"
	"context"
)

// StatusClient for client basic action by rpc
type StatusClient interface {
	Start(ctx context.Context, wg *sync.WaitGroup) error
	ChangeStatus(token string) error
	Query(token string) (string, error)
	Set(token string , value string)
}