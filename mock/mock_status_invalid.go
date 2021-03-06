// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/JointFaaS/Storage-Center/status (interfaces: Maintainer_InvalidClient)

// Package mock_status is a generated GoMock package.
package mock_status

import (
	context "context"
	status "github.com/JointFaaS/Storage-Center/status"
	gomock "github.com/golang/mock/gomock"
	metadata "google.golang.org/grpc/metadata"
	reflect "reflect"
)

// MockMaintainer_InvalidClient is a mock of Maintainer_InvalidClient interface
type MockMaintainer_InvalidClient struct {
	ctrl     *gomock.Controller
	recorder *MockMaintainer_InvalidClientMockRecorder
}

// MockMaintainer_InvalidClientMockRecorder is the mock recorder for MockMaintainer_InvalidClient
type MockMaintainer_InvalidClientMockRecorder struct {
	mock *MockMaintainer_InvalidClient
}

// NewMockMaintainer_InvalidClient creates a new mock instance
func NewMockMaintainer_InvalidClient(ctrl *gomock.Controller) *MockMaintainer_InvalidClient {
	mock := &MockMaintainer_InvalidClient{ctrl: ctrl}
	mock.recorder = &MockMaintainer_InvalidClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMaintainer_InvalidClient) EXPECT() *MockMaintainer_InvalidClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockMaintainer_InvalidClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockMaintainer_InvalidClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockMaintainer_InvalidClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockMaintainer_InvalidClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).Context))
}

// Header mocks base method
func (m *MockMaintainer_InvalidClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockMaintainer_InvalidClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).Header))
}

// Recv mocks base method
func (m *MockMaintainer_InvalidClient) Recv() (*status.InvalidReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*status.InvalidReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockMaintainer_InvalidClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockMaintainer_InvalidClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockMaintainer_InvalidClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockMaintainer_InvalidClient) Send(arg0 *status.InvalidRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockMaintainer_InvalidClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).Send), arg0)
}

// SendMsg mocks base method
func (m *MockMaintainer_InvalidClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockMaintainer_InvalidClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockMaintainer_InvalidClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockMaintainer_InvalidClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockMaintainer_InvalidClient)(nil).Trailer))
}
