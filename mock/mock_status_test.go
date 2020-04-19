package mock_status

import (
	"fmt"
	"sync"
	// "context"
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

	var wg sync.WaitGroup
  	defer wg.Wait()
	mockMaintainerClient := NewMockMaintainerClient(ctrl)
	stream := NewMockMaintainer_InvalidClient(ctrl)

	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(&pb.InvalidReply{Token: ""}, nil)
	// stream.EXPECT().CloseSend().DoAndReturn(func() {
	// 	return
	// })
	mockMaintainerClient.EXPECT().Register(
		gomock.Any(),
		// &rpcMsg{msg: req},
		gomock.Any(),
	).Return(&pb.RegisterReply{Code: 1, Msg:"OK"}, nil)
	mockMaintainerClient.EXPECT().Invalid(
		gomock.Any(),
	).Return(stream, nil)

	// // for missing call
	// invalidStream, err := mockMaintainerClient.Invalid(context.Background())
	// if err != nil {
	// 	t.Errorf("openn stream error %v", err)
	// }
	
	// if err := invalidStream.Send(&pb.InvalidRequest{Name: "test"}); err != nil {
	// 	t.Errorf("can not send %v", err)
	// }
	// _, err = invalidStream.Recv()


	var c inter.UserClient
	tmp := client.NewUserClientImpl("test", "127.0.0.1:50001", "127.0.0.1:50000", mockMaintainerClient)
	c = &tmp
	err := c.Start()
	if (err != nil) {
		t.Errorf(err.Error())
	}
	c.Close()
}

// func Test_Simple_Create_Pair(t *testing.T) {
// 	key := "casecloud"
// 	storageValue := "nmsl"
// 	clientHost := "127.0.0.1:50001"
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockMaintainerClient := NewMockMaintainerClient(ctrl)
// 	stream := NewMockMaintainer_InvalidClient(ctrl)

// 	stream.EXPECT().Send(
// 		gomock.Any(),
// 	).Return(nil)
// 	// Set expectation on receiving.
// 	stream.EXPECT().Recv().DoAndReturn(func() {
// 		time.Sleep(time.Duration(1)*time.Second)
// 		return pb.InvalidReply{Token:""}
// 	})
// 	stream.EXPECT().CloseSend().Return(nil)

// 	mockMaintainerClient.EXPECT().Register(
// 		gomock.Any(),
// 		gomock.Any(),
// 	).Return(&pb.RegisterReply{Code: 1, Msg:"OK"}, nil)
// 	mockMaintainerClient.EXPECT().ChangeStatus(
// 		gomock.Any(),
// 		gomock.Any(),
// 	).Return(&pb.StatusReply{Token: key, Host: clientHost}, nil)
// 	mockMaintainerClient.EXPECT().Invalid(
// 		gomock.Any(),
// 	).Return(stream, nil)
// 	var c inter.UserClient
// 	tmp := client.NewUserClientImpl("test", clientHost, "127.0.0.1:50000", mockMaintainerClient)
// 	c = &tmp
// 	err := c.Start()
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}
// 	err = c.Set(key, storageValue)
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}
// 	value, err := c.Get(key)
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}
// 	if value != storageValue {
// 		t.Errorf("value %v should be %v", value, storageValue)
// 	}
// }

