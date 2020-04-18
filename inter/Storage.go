package inter

// Storage to storage value
type Storage interface {
	Set(token string, value string)
	Get(token string) string
	Delete(token string) error
}