// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	dao "github.com/in-rich/uservice-notes/pkg/dao"
	entities "github.com/in-rich/uservice-notes/pkg/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockListNotesRepository is an autogenerated mock type for the ListNotesRepository type
type MockListNotesRepository struct {
	mock.Mock
}

type MockListNotesRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockListNotesRepository) EXPECT() *MockListNotesRepository_Expecter {
	return &MockListNotesRepository_Expecter{mock: &_m.Mock}
}

// ListNotes provides a mock function with given fields: ctx, authorID, filters
func (_m *MockListNotesRepository) ListNotes(ctx context.Context, authorID string, filters []dao.ListNoteFilter) ([]*entities.Note, error) {
	ret := _m.Called(ctx, authorID, filters)

	if len(ret) == 0 {
		panic("no return value specified for ListNotes")
	}

	var r0 []*entities.Note
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []dao.ListNoteFilter) ([]*entities.Note, error)); ok {
		return rf(ctx, authorID, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []dao.ListNoteFilter) []*entities.Note); ok {
		r0 = rf(ctx, authorID, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Note)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []dao.ListNoteFilter) error); ok {
		r1 = rf(ctx, authorID, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockListNotesRepository_ListNotes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListNotes'
type MockListNotesRepository_ListNotes_Call struct {
	*mock.Call
}

// ListNotes is a helper method to define mock.On call
//   - ctx context.Context
//   - authorID string
//   - filters []dao.ListNoteFilter
func (_e *MockListNotesRepository_Expecter) ListNotes(ctx interface{}, authorID interface{}, filters interface{}) *MockListNotesRepository_ListNotes_Call {
	return &MockListNotesRepository_ListNotes_Call{Call: _e.mock.On("ListNotes", ctx, authorID, filters)}
}

func (_c *MockListNotesRepository_ListNotes_Call) Run(run func(ctx context.Context, authorID string, filters []dao.ListNoteFilter)) *MockListNotesRepository_ListNotes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]dao.ListNoteFilter))
	})
	return _c
}

func (_c *MockListNotesRepository_ListNotes_Call) Return(_a0 []*entities.Note, _a1 error) *MockListNotesRepository_ListNotes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockListNotesRepository_ListNotes_Call) RunAndReturn(run func(context.Context, string, []dao.ListNoteFilter) ([]*entities.Note, error)) *MockListNotesRepository_ListNotes_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockListNotesRepository creates a new instance of MockListNotesRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockListNotesRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockListNotesRepository {
	mock := &MockListNotesRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
