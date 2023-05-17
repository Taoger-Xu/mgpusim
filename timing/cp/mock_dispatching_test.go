// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/akita/mgpusim/v3/timing/cp/internal/dispatching (interfaces: Dispatcher)

package cp

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	sim "gitlab.com/akita/akita/v3/sim"
	protocol "gitlab.com/akita/mgpusim/v3/protocol"
	resource "gitlab.com/akita/mgpusim/v3/timing/cp/internal/resource"
)

// MockDispatcher is a mock of Dispatcher interface.
type MockDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockDispatcherMockRecorder
}

// MockDispatcherMockRecorder is the mock recorder for MockDispatcher.
type MockDispatcherMockRecorder struct {
	mock *MockDispatcher
}

// NewMockDispatcher creates a new mock instance.
func NewMockDispatcher(ctrl *gomock.Controller) *MockDispatcher {
	mock := &MockDispatcher{ctrl: ctrl}
	mock.recorder = &MockDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDispatcher) EXPECT() *MockDispatcherMockRecorder {
	return m.recorder
}

// AcceptHook mocks base method.
func (m *MockDispatcher) AcceptHook(arg0 sim.Hook) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AcceptHook", arg0)
}

// AcceptHook indicates an expected call of AcceptHook.
func (mr *MockDispatcherMockRecorder) AcceptHook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptHook", reflect.TypeOf((*MockDispatcher)(nil).AcceptHook), arg0)
}

// Hooks mocks base method.
func (m *MockDispatcher) Hooks() []sim.Hook {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hooks")
	ret0, _ := ret[0].([]sim.Hook)
	return ret0
}

// Hooks indicates an expected call of Hooks.
func (mr *MockDispatcherMockRecorder) Hooks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hooks", reflect.TypeOf((*MockDispatcher)(nil).Hooks))
}

// InvokeHook mocks base method.
func (m *MockDispatcher) InvokeHook(arg0 sim.HookCtx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InvokeHook", arg0)
}

// InvokeHook indicates an expected call of InvokeHook.
func (mr *MockDispatcherMockRecorder) InvokeHook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeHook", reflect.TypeOf((*MockDispatcher)(nil).InvokeHook), arg0)
}

// IsDispatching mocks base method.
func (m *MockDispatcher) IsDispatching() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDispatching")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDispatching indicates an expected call of IsDispatching.
func (mr *MockDispatcherMockRecorder) IsDispatching() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDispatching", reflect.TypeOf((*MockDispatcher)(nil).IsDispatching))
}

// Name mocks base method.
func (m *MockDispatcher) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockDispatcherMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockDispatcher)(nil).Name))
}

// NumHooks mocks base method.
func (m *MockDispatcher) NumHooks() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumHooks")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumHooks indicates an expected call of NumHooks.
func (mr *MockDispatcherMockRecorder) NumHooks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumHooks", reflect.TypeOf((*MockDispatcher)(nil).NumHooks))
}

// RegisterCU mocks base method.
func (m *MockDispatcher) RegisterCU(arg0 resource.DispatchableCU) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterCU", arg0)
}

// RegisterCU indicates an expected call of RegisterCU.
func (mr *MockDispatcherMockRecorder) RegisterCU(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCU", reflect.TypeOf((*MockDispatcher)(nil).RegisterCU), arg0)
}

// StartDispatching mocks base method.
func (m *MockDispatcher) StartDispatching(arg0 *protocol.LaunchKernelReq) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartDispatching", arg0)
}

// StartDispatching indicates an expected call of StartDispatching.
func (mr *MockDispatcherMockRecorder) StartDispatching(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartDispatching", reflect.TypeOf((*MockDispatcher)(nil).StartDispatching), arg0)
}

// Tick mocks base method.
func (m *MockDispatcher) Tick(arg0 sim.VTimeInSec) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tick", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Tick indicates an expected call of Tick.
func (mr *MockDispatcherMockRecorder) Tick(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tick", reflect.TypeOf((*MockDispatcher)(nil).Tick), arg0)
}
