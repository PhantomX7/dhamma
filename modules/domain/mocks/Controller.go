// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// Controller is an autogenerated mock type for the Controller type
type Controller struct {
	mock.Mock
}

type Controller_Expecter struct {
	mock *mock.Mock
}

func (_m *Controller) EXPECT() *Controller_Expecter {
	return &Controller_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx
func (_m *Controller) Create(ctx *gin.Context) {
	_m.Called(ctx)
}

// Controller_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type Controller_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx *gin.Context
func (_e *Controller_Expecter) Create(ctx interface{}) *Controller_Create_Call {
	return &Controller_Create_Call{Call: _e.mock.On("Create", ctx)}
}

func (_c *Controller_Create_Call) Run(run func(ctx *gin.Context)) *Controller_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *Controller_Create_Call) Return() *Controller_Create_Call {
	_c.Call.Return()
	return _c
}

func (_c *Controller_Create_Call) RunAndReturn(run func(*gin.Context)) *Controller_Create_Call {
	_c.Run(run)
	return _c
}

// Index provides a mock function with given fields: ctx
func (_m *Controller) Index(ctx *gin.Context) {
	_m.Called(ctx)
}

// Controller_Index_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Index'
type Controller_Index_Call struct {
	*mock.Call
}

// Index is a helper method to define mock.On call
//   - ctx *gin.Context
func (_e *Controller_Expecter) Index(ctx interface{}) *Controller_Index_Call {
	return &Controller_Index_Call{Call: _e.mock.On("Index", ctx)}
}

func (_c *Controller_Index_Call) Run(run func(ctx *gin.Context)) *Controller_Index_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *Controller_Index_Call) Return() *Controller_Index_Call {
	_c.Call.Return()
	return _c
}

func (_c *Controller_Index_Call) RunAndReturn(run func(*gin.Context)) *Controller_Index_Call {
	_c.Run(run)
	return _c
}

// Show provides a mock function with given fields: ctx
func (_m *Controller) Show(ctx *gin.Context) {
	_m.Called(ctx)
}

// Controller_Show_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Show'
type Controller_Show_Call struct {
	*mock.Call
}

// Show is a helper method to define mock.On call
//   - ctx *gin.Context
func (_e *Controller_Expecter) Show(ctx interface{}) *Controller_Show_Call {
	return &Controller_Show_Call{Call: _e.mock.On("Show", ctx)}
}

func (_c *Controller_Show_Call) Run(run func(ctx *gin.Context)) *Controller_Show_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *Controller_Show_Call) Return() *Controller_Show_Call {
	_c.Call.Return()
	return _c
}

func (_c *Controller_Show_Call) RunAndReturn(run func(*gin.Context)) *Controller_Show_Call {
	_c.Run(run)
	return _c
}

// Update provides a mock function with given fields: ctx
func (_m *Controller) Update(ctx *gin.Context) {
	_m.Called(ctx)
}

// Controller_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type Controller_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx *gin.Context
func (_e *Controller_Expecter) Update(ctx interface{}) *Controller_Update_Call {
	return &Controller_Update_Call{Call: _e.mock.On("Update", ctx)}
}

func (_c *Controller_Update_Call) Run(run func(ctx *gin.Context)) *Controller_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *Controller_Update_Call) Return() *Controller_Update_Call {
	_c.Call.Return()
	return _c
}

func (_c *Controller_Update_Call) RunAndReturn(run func(*gin.Context)) *Controller_Update_Call {
	_c.Run(run)
	return _c
}

// NewController creates a new instance of Controller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewController(t interface {
	mock.TestingT
	Cleanup(func())
}) *Controller {
	mock := &Controller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
