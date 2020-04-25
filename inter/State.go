package inter

// State for maintain StateInterface
type State interface {
	ChangeStatus(token string, host string) (string, string, uint64, uint64, error)
	Query(token string) (string, uint64, error)
}