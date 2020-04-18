package client

import(
	inter "github.com/JointFaaS/Storage-Center/inter"
	"testing"
)


func Test_Init(t *testing.T) {
	var c inter.UserClient
	tmp := NewClientImpl("test", "127.0.0.1:50001", "127.0.0.1:50000", nil)
	c = &tmp
	if (c == nil) {
		t.Error("client not init")
	}
}