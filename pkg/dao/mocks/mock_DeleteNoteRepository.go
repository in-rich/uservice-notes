// Code generated by mockery v2.46.0. DO NOT EDIT.

package daomocks

import (
	context "context"

	entities "github.com/in-rich/uservice-notes/pkg/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockDeleteNoteRepository is an autogenerated mock type for the DeleteNoteRepository type
type MockDeleteNoteRepository struct {
	mock.Mock
}

type MockDeleteNoteRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDeleteNoteRepository) EXPECT() *MockDeleteNoteRepository_Expecter {
	return &MockDeleteNoteRepository_Expecter{mock: &_m.Mock}
}

// DeleteNote provides a mock function with given fields: ctx, author, target, publicIdentifier
func (_m *MockDeleteNoteRepository) DeleteNote(ctx context.Context, author string, target entities.Target, publicIdentifier string) (*entities.Note, error) {
	ret := _m.Called(ctx, author, target, publicIdentifier)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNote")
	}

	var r0 *entities.Note
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, entities.Target, string) (*entities.Note, error)); ok {
		return rf(ctx, author, target, publicIdentifier)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, entities.Target, string) *entities.Note); ok {
		r0 = rf(ctx, author, target, publicIdentifier)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Note)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, entities.Target, string) error); ok {
		r1 = rf(ctx, author, target, publicIdentifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDeleteNoteRepository_DeleteNote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteNote'
type MockDeleteNoteRepository_DeleteNote_Call struct {
	*mock.Call
}

// DeleteNote is a helper method to define mock.On call
//   - ctx context.Context
//   - author string
//   - target entities.Target
//   - publicIdentifier string
func (_e *MockDeleteNoteRepository_Expecter) DeleteNote(ctx interface{}, author interface{}, target interface{}, publicIdentifier interface{}) *MockDeleteNoteRepository_DeleteNote_Call {
	return &MockDeleteNoteRepository_DeleteNote_Call{Call: _e.mock.On("DeleteNote", ctx, author, target, publicIdentifier)}
}

func (_c *MockDeleteNoteRepository_DeleteNote_Call) Run(run func(ctx context.Context, author string, target entities.Target, publicIdentifier string)) *MockDeleteNoteRepository_DeleteNote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(entities.Target), args[3].(string))
	})
	return _c
}

func (_c *MockDeleteNoteRepository_DeleteNote_Call) Return(_a0 *entities.Note, _a1 error) *MockDeleteNoteRepository_DeleteNote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDeleteNoteRepository_DeleteNote_Call) RunAndReturn(run func(context.Context, string, entities.Target, string) (*entities.Note, error)) *MockDeleteNoteRepository_DeleteNote_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDeleteNoteRepository creates a new instance of MockDeleteNoteRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeleteNoteRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeleteNoteRepository {
	mock := &MockDeleteNoteRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
