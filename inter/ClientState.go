package inter

import (
)

// ClientState is an interface about client side state maintainer
type ClientState interface {
	Delete(token string)
	Query(token string) bool
	Add(token string) bool
}