package inter

// ClientState is an interface about client side state maintainer
type ClientState interface {
	Delete(token string)
	Query(token string) bool
	GetVersion(token string) uint64
	Add(token string, version uint64) bool
}
