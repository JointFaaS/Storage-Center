package inter

// StatusClient for client basic action by rpc
type StatusClient interface {
	Start() error
	ChangeStatus(token string) error
	Query(token string) (string, error)
	Set(token string , value string)
}