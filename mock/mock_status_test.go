package mock_status

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	pb "github.com/JointFaaS/Storage-Center/status"
	client "github.com/JointFaaS/Storage-Center/client"
	inter "github.com/JointFaaS/Storage-Center/inter"
)




// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func Test_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMaintainerClient := NewMockMaintainerClient(ctrl)
	// req := &pb.RegisterRequest{Name: "test", Host: "127.0.0.1:50001"}
	mockMaintainerClient.EXPECT().Register(
		gomock.Any(),
		// &rpcMsg{msg: req},
		gomock.Any(),
	).Return(&pb.RegisterReply{Code: 1, Msg:"OK"}, nil)
	var c inter.UserClient
	tmp := client.NewUserClientImpl("test", "127.0.0.1:50001", "127.0.0.1:50000", mockMaintainerClient)
	c = &tmp
	err := c.Start()
	if (err != nil) {
		t.Errorf(err.Error())
	}
}

// func Test_Init(t *testing.T) {
// 	var c inter.UserClient
// 	tmp := client.NewUserClientImpl("test", "127.0.0.1:50001", "127.0.0.1:50000", mockInit(t))
// 	c = &tmp
// 	err := c.Start()
// 	if (err != nil) {
// 		t.Errorf(err.Error())
// 	}
// }

