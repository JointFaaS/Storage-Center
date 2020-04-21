package inter

// UserClient for client basic action by rpc
type UserClient interface {
	Start() error
	Close()
	Get(token string) (string, error)
	Set(token string, value string) error
}