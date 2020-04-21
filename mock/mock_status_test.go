package mock_status

import (
	"testing"

	client "github.com/JointFaaS/Storage-Center/client"
	inter "github.com/JointFaaS/Storage-Center/inter"
	pb "github.com/JointFaaS/Storage-Center/status"
	"github.com/golang/mock/gomock"
)

func Test_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMaintainerClient := NewMockMaintainerClient(ctrl)
	stream := NewMockMaintainer_InvalidClient(ctrl)

	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil).AnyTimes()
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(&pb.InvalidReply{Token: ""}, nil).AnyTimes()
	stream.EXPECT().CloseSend().Return(nil).AnyTimes()
	mockMaintainerClient.EXPECT().Register(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.RegisterReply{Code: 1, Msg: "OK"}, nil)
	mockMaintainerClient.EXPECT().Invalid(
		gomock.Any(),
	).Return(stream, nil).AnyTimes()

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
	c = client.NewUserClientImpl("test", "127.0.0.1", ":50001", "127.0.0.1:50000", mockMaintainerClient)
	err := c.Start()
	if err != nil {
		t.Errorf(err.Error())
	}
	c.Close()
}

func Test_Simple_Create_Pair(t *testing.T) {
	key := "animal"
	storageValue := "pig"
	clientHost := "127.0.0.1"
	clientPort := ":50001"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMaintainerClient := NewMockMaintainerClient(ctrl)
	stream := NewMockMaintainer_InvalidClient(ctrl)

	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil).AnyTimes()
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(&pb.InvalidReply{Token: ""}, nil).AnyTimes()
	stream.EXPECT().CloseSend().Return(nil).AnyTimes()

	mockMaintainerClient.EXPECT().Register(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.RegisterReply{Code: 1, Msg: "OK"}, nil)
	mockMaintainerClient.EXPECT().ChangeStatus(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.StatusReply{Token: key, Host: clientHost}, nil)
	mockMaintainerClient.EXPECT().Invalid(
		gomock.Any(),
	).Return(stream, nil)
	var c inter.UserClient
	c = client.NewUserClientImpl("test", clientHost, clientPort, "127.0.0.1:50000", mockMaintainerClient)

	err := c.Start()
	defer c.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = c.Set(key, storageValue)
	if err != nil {
		t.Errorf(err.Error())
	}
	value, err := c.Get(key)
	if err != nil {
		t.Errorf(err.Error())
	}
	if value != storageValue {
		t.Errorf("value %v should be %v", value, storageValue)
	}
}

func Test_Get_Value_From_Others(t *testing.T) {
	key := "animal"
	storageValue := "pig"
	clientHost1 := "127.0.0.1"
	clientPort1 := ":50001"
	clientHost2 := "127.0.0.1"
	clientPort2 := ":50002"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMaintainerClient := NewMockMaintainerClient(ctrl)
	stream := NewMockMaintainer_InvalidClient(ctrl)

	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil).AnyTimes()
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(&pb.InvalidReply{Token: ""}, nil).AnyTimes()
	stream.EXPECT().CloseSend().Return(nil).AnyTimes()

	mockMaintainerClient.EXPECT().Register(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.RegisterReply{Code: 1, Msg: "OK"}, nil).AnyTimes()
	mockMaintainerClient.EXPECT().Query(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.QueryReply{Token: key, Host: clientHost1 + clientPort1}, nil).AnyTimes()
	mockMaintainerClient.EXPECT().ChangeStatus(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.StatusReply{Token: key, Host: clientHost1}, nil).AnyTimes()
	mockMaintainerClient.EXPECT().Invalid(
		gomock.Any(),
	).Return(stream, nil).AnyTimes()
	var c1 inter.UserClient
	var c2 inter.UserClient
	c1 = client.NewUserClientImpl("test1", clientHost1, clientPort1, "127.0.0.1:50000", mockMaintainerClient)
	defer c1.Close()
	c2 = client.NewUserClientImpl("test2", clientHost2, clientPort2, "127.0.0.1:50000", mockMaintainerClient)
	defer c2.Close()
	err := c1.Start()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = c1.Set(key, storageValue)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = c2.Start()
	if err != nil {
		t.Errorf(err.Error())
	}
	value, err := c2.Get(key)
	if err != nil {
		t.Errorf(err.Error())
	}
	if value != storageValue {
		t.Errorf("value %v should be %v", value, storageValue)
	}
}
