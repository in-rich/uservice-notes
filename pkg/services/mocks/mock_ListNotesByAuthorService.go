// Code generated by mockery v2.46.0. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	models "github.com/in-rich/uservice-notes/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// MockListNotesByAuthorService is an autogenerated mock type for the ListNotesByAuthorService type
type MockListNotesByAuthorService struct {
	mock.Mock
}

type MockListNotesByAuthorService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockListNotesByAuthorService) EXPECT() *MockListNotesByAuthorService_Expecter {
	return &MockListNotesByAuthorService_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, selector
func (_m *MockListNotesByAuthorService) Exec(ctx context.Context, selector *models.ListNotesByAuthor) ([]*models.Note, error) {
	ret := _m.Called(ctx, selector)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 []*models.Note
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.ListNotesByAuthor) ([]*models.Note, error)); ok {
		return rf(ctx, selector)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.ListNotesByAuthor) []*models.Note); ok {
		r0 = rf(ctx, selector)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Note)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.ListNotesByAuthor) error); ok {
		r1 = rf(ctx, selector)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockListNotesByAuthorService_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockListNotesByAuthorService_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - selector *models.ListNotesByAuthor
func (_e *MockListNotesByAuthorService_Expecter) Exec(ctx interface{}, selector interface{}) *MockListNotesByAuthorService_Exec_Call {
	return &MockListNotesByAuthorService_Exec_Call{Call: _e.mock.On("Exec", ctx, selector)}
}

func (_c *MockListNotesByAuthorService_Exec_Call) Run(run func(ctx context.Context, selector *models.ListNotesByAuthor)) *MockListNotesByAuthorService_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.ListNotesByAuthor))
	})
	return _c
}

func (_c *MockListNotesByAuthorService_Exec_Call) Return(_a0 []*models.Note, _a1 error) *MockListNotesByAuthorService_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockListNotesByAuthorService_Exec_Call) RunAndReturn(run func(context.Context, *models.ListNotesByAuthor) ([]*models.Note, error)) *MockListNotesByAuthorService_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockListNotesByAuthorService creates a new instance of MockListNotesByAuthorService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockListNotesByAuthorService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockListNotesByAuthorService {
	mock := &MockListNotesByAuthorService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
