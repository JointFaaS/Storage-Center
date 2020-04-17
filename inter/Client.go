package inter

import (
	enum "github.com/JointFaaS/Storage-Center/enum"
)

// ClientInterface for client basic action
type ClientInterface interface {
	Register(name string ,host string) error
	ChangeStatus(token string, status enum.Status) error
	Query(token string) (string, error)
}