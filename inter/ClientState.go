package inter

import (
	"github.com/JointFaaS/Storage-Center/enum"
)

// ClientState is an interface about client side state maintainer
type ClientState interface {
	Delete(token string)
	Query(token string) bool
	GetVersion(token string) uint64
	GetStatus(token string) enum.Policy
	Add(token string, version uint64, status enum.Policy) bool
}
