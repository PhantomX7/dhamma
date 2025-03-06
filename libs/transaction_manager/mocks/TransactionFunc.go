// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// TransactionFunc is an autogenerated mock type for the TransactionFunc type
type TransactionFunc struct {
	mock.Mock
}

type TransactionFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *TransactionFunc) EXPECT() *TransactionFunc_Expecter {
	return &TransactionFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: tx
func (_m *TransactionFunc) Execute(tx *gorm.DB) error {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransactionFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type TransactionFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - tx *gorm.DB
func (_e *TransactionFunc_Expecter) Execute(tx interface{}) *TransactionFunc_Execute_Call {
	return &TransactionFunc_Execute_Call{Call: _e.mock.On("Execute", tx)}
}

func (_c *TransactionFunc_Execute_Call) Run(run func(tx *gorm.DB)) *TransactionFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gorm.DB))
	})
	return _c
}

func (_c *TransactionFunc_Execute_Call) Return(_a0 error) *TransactionFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransactionFunc_Execute_Call) RunAndReturn(run func(*gorm.DB) error) *TransactionFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransactionFunc creates a new instance of TransactionFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionFunc {
	mock := &TransactionFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
