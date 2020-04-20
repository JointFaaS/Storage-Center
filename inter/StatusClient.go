package inter

import (
	"context"
	"sync"
)

// StatusClient for client basic action by rpc
type StatusClient interface {
	Start(ctx context.Context, wg *sync.WaitGroup) error
	ChangeStatus(ctx context.Context, token string) error
	Query(ctx context.Context, token string) (string, error)
	Set(token string, value string)
}
