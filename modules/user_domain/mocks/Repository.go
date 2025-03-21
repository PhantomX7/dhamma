// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/PhantomX7/dhamma/entity"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// AssignDomain provides a mock function with given fields: ctx, userID, domainID, tx
func (_m *Repository) AssignDomain(ctx context.Context, userID uint64, domainID uint64, tx *gorm.DB) error {
	ret := _m.Called(ctx, userID, domainID, tx)

	if len(ret) == 0 {
		panic("no return value specified for AssignDomain")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, *gorm.DB) error); ok {
		r0 = rf(ctx, userID, domainID, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_AssignDomain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AssignDomain'
type Repository_AssignDomain_Call struct {
	*mock.Call
}

// AssignDomain is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - domainID uint64
//   - tx *gorm.DB
func (_e *Repository_Expecter) AssignDomain(ctx interface{}, userID interface{}, domainID interface{}, tx interface{}) *Repository_AssignDomain_Call {
	return &Repository_AssignDomain_Call{Call: _e.mock.On("AssignDomain", ctx, userID, domainID, tx)}
}

func (_c *Repository_AssignDomain_Call) Run(run func(ctx context.Context, userID uint64, domainID uint64, tx *gorm.DB)) *Repository_AssignDomain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(uint64), args[3].(*gorm.DB))
	})
	return _c
}

func (_c *Repository_AssignDomain_Call) Return(_a0 error) *Repository_AssignDomain_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_AssignDomain_Call) RunAndReturn(run func(context.Context, uint64, uint64, *gorm.DB) error) *Repository_AssignDomain_Call {
	_c.Call.Return(run)
	return _c
}

// FindByUserID provides a mock function with given fields: ctx, userID, preloadRelations
func (_m *Repository) FindByUserID(ctx context.Context, userID uint64, preloadRelations bool) ([]entity.UserDomain, error) {
	ret := _m.Called(ctx, userID, preloadRelations)

	if len(ret) == 0 {
		panic("no return value specified for FindByUserID")
	}

	var r0 []entity.UserDomain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, bool) ([]entity.UserDomain, error)); ok {
		return rf(ctx, userID, preloadRelations)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, bool) []entity.UserDomain); ok {
		r0 = rf(ctx, userID, preloadRelations)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.UserDomain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, bool) error); ok {
		r1 = rf(ctx, userID, preloadRelations)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_FindByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByUserID'
type Repository_FindByUserID_Call struct {
	*mock.Call
}

// FindByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - preloadRelations bool
func (_e *Repository_Expecter) FindByUserID(ctx interface{}, userID interface{}, preloadRelations interface{}) *Repository_FindByUserID_Call {
	return &Repository_FindByUserID_Call{Call: _e.mock.On("FindByUserID", ctx, userID, preloadRelations)}
}

func (_c *Repository_FindByUserID_Call) Run(run func(ctx context.Context, userID uint64, preloadRelations bool)) *Repository_FindByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(bool))
	})
	return _c
}

func (_c *Repository_FindByUserID_Call) Return(_a0 []entity.UserDomain, _a1 error) *Repository_FindByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_FindByUserID_Call) RunAndReturn(run func(context.Context, uint64, bool) ([]entity.UserDomain, error)) *Repository_FindByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// HasDomain provides a mock function with given fields: ctx, userID, domainID
func (_m *Repository) HasDomain(ctx context.Context, userID uint64, domainID uint64) (bool, error) {
	ret := _m.Called(ctx, userID, domainID)

	if len(ret) == 0 {
		panic("no return value specified for HasDomain")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) (bool, error)); ok {
		return rf(ctx, userID, domainID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) bool); ok {
		r0 = rf(ctx, userID, domainID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64) error); ok {
		r1 = rf(ctx, userID, domainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_HasDomain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasDomain'
type Repository_HasDomain_Call struct {
	*mock.Call
}

// HasDomain is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - domainID uint64
func (_e *Repository_Expecter) HasDomain(ctx interface{}, userID interface{}, domainID interface{}) *Repository_HasDomain_Call {
	return &Repository_HasDomain_Call{Call: _e.mock.On("HasDomain", ctx, userID, domainID)}
}

func (_c *Repository_HasDomain_Call) Run(run func(ctx context.Context, userID uint64, domainID uint64)) *Repository_HasDomain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(uint64))
	})
	return _c
}

func (_c *Repository_HasDomain_Call) Return(_a0 bool, _a1 error) *Repository_HasDomain_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_HasDomain_Call) RunAndReturn(run func(context.Context, uint64, uint64) (bool, error)) *Repository_HasDomain_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
