// Code generated by go-mockgen 1.3.7; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package http

import (
	"context"
	"sync"

	uploadhandler "github.com/khulnasoft/khulnasoft/internal/uploadhandler"
)

// MockDBStore is a mock implementation of the DBStore interface (from the
// package github.com/khulnasoft/khulnasoft/internal/uploadhandler) used
// for unit testing.
type MockDBStore[T interface{}] struct {
	// AddUploadPartFunc is an instance of a mock function object
	// controlling the behavior of the method AddUploadPart.
	AddUploadPartFunc *DBStoreAddUploadPartFunc[T]
	// GetUploadByIDFunc is an instance of a mock function object
	// controlling the behavior of the method GetUploadByID.
	GetUploadByIDFunc *DBStoreGetUploadByIDFunc[T]
	// InsertUploadFunc is an instance of a mock function object controlling
	// the behavior of the method InsertUpload.
	InsertUploadFunc *DBStoreInsertUploadFunc[T]
	// MarkFailedFunc is an instance of a mock function object controlling
	// the behavior of the method MarkFailed.
	MarkFailedFunc *DBStoreMarkFailedFunc[T]
	// MarkQueuedFunc is an instance of a mock function object controlling
	// the behavior of the method MarkQueued.
	MarkQueuedFunc *DBStoreMarkQueuedFunc[T]
	// WithTransactionFunc is an instance of a mock function object
	// controlling the behavior of the method WithTransaction.
	WithTransactionFunc *DBStoreWithTransactionFunc[T]
}

// NewMockDBStore creates a new mock of the DBStore interface. All methods
// return zero values for all results, unless overwritten.
func NewMockDBStore[T interface{}]() *MockDBStore[T] {
	return &MockDBStore[T]{
		AddUploadPartFunc: &DBStoreAddUploadPartFunc[T]{
			defaultHook: func(context.Context, int, int) (r0 error) {
				return
			},
		},
		GetUploadByIDFunc: &DBStoreGetUploadByIDFunc[T]{
			defaultHook: func(context.Context, int) (r0 uploadhandler.Upload[T], r1 bool, r2 error) {
				return
			},
		},
		InsertUploadFunc: &DBStoreInsertUploadFunc[T]{
			defaultHook: func(context.Context, uploadhandler.Upload[T]) (r0 int, r1 error) {
				return
			},
		},
		MarkFailedFunc: &DBStoreMarkFailedFunc[T]{
			defaultHook: func(context.Context, int, string) (r0 error) {
				return
			},
		},
		MarkQueuedFunc: &DBStoreMarkQueuedFunc[T]{
			defaultHook: func(context.Context, int, *int64) (r0 error) {
				return
			},
		},
		WithTransactionFunc: &DBStoreWithTransactionFunc[T]{
			defaultHook: func(context.Context, func(tx uploadhandler.DBStore[T]) error) (r0 error) {
				return
			},
		},
	}
}

// NewStrictMockDBStore creates a new mock of the DBStore interface. All
// methods panic on invocation, unless overwritten.
func NewStrictMockDBStore[T interface{}]() *MockDBStore[T] {
	return &MockDBStore[T]{
		AddUploadPartFunc: &DBStoreAddUploadPartFunc[T]{
			defaultHook: func(context.Context, int, int) error {
				panic("unexpected invocation of MockDBStore.AddUploadPart")
			},
		},
		GetUploadByIDFunc: &DBStoreGetUploadByIDFunc[T]{
			defaultHook: func(context.Context, int) (uploadhandler.Upload[T], bool, error) {
				panic("unexpected invocation of MockDBStore.GetUploadByID")
			},
		},
		InsertUploadFunc: &DBStoreInsertUploadFunc[T]{
			defaultHook: func(context.Context, uploadhandler.Upload[T]) (int, error) {
				panic("unexpected invocation of MockDBStore.InsertUpload")
			},
		},
		MarkFailedFunc: &DBStoreMarkFailedFunc[T]{
			defaultHook: func(context.Context, int, string) error {
				panic("unexpected invocation of MockDBStore.MarkFailed")
			},
		},
		MarkQueuedFunc: &DBStoreMarkQueuedFunc[T]{
			defaultHook: func(context.Context, int, *int64) error {
				panic("unexpected invocation of MockDBStore.MarkQueued")
			},
		},
		WithTransactionFunc: &DBStoreWithTransactionFunc[T]{
			defaultHook: func(context.Context, func(tx uploadhandler.DBStore[T]) error) error {
				panic("unexpected invocation of MockDBStore.WithTransaction")
			},
		},
	}
}

// NewMockDBStoreFrom creates a new mock of the MockDBStore interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockDBStoreFrom[T interface{}](i uploadhandler.DBStore[T]) *MockDBStore[T] {
	return &MockDBStore[T]{
		AddUploadPartFunc: &DBStoreAddUploadPartFunc[T]{
			defaultHook: i.AddUploadPart,
		},
		GetUploadByIDFunc: &DBStoreGetUploadByIDFunc[T]{
			defaultHook: i.GetUploadByID,
		},
		InsertUploadFunc: &DBStoreInsertUploadFunc[T]{
			defaultHook: i.InsertUpload,
		},
		MarkFailedFunc: &DBStoreMarkFailedFunc[T]{
			defaultHook: i.MarkFailed,
		},
		MarkQueuedFunc: &DBStoreMarkQueuedFunc[T]{
			defaultHook: i.MarkQueued,
		},
		WithTransactionFunc: &DBStoreWithTransactionFunc[T]{
			defaultHook: i.WithTransaction,
		},
	}
}

// DBStoreAddUploadPartFunc describes the behavior when the AddUploadPart
// method of the parent MockDBStore instance is invoked.
type DBStoreAddUploadPartFunc[T interface{}] struct {
	defaultHook func(context.Context, int, int) error
	hooks       []func(context.Context, int, int) error
	history     []DBStoreAddUploadPartFuncCall[T]
	mutex       sync.Mutex
}

// AddUploadPart delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockDBStore[T]) AddUploadPart(v0 context.Context, v1 int, v2 int) error {
	r0 := m.AddUploadPartFunc.nextHook()(v0, v1, v2)
	m.AddUploadPartFunc.appendCall(DBStoreAddUploadPartFuncCall[T]{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the AddUploadPart method
// of the parent MockDBStore instance is invoked and the hook queue is
// empty.
func (f *DBStoreAddUploadPartFunc[T]) SetDefaultHook(hook func(context.Context, int, int) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// AddUploadPart method of the parent MockDBStore instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *DBStoreAddUploadPartFunc[T]) PushHook(hook func(context.Context, int, int) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DBStoreAddUploadPartFunc[T]) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, int, int) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DBStoreAddUploadPartFunc[T]) PushReturn(r0 error) {
	f.PushHook(func(context.Context, int, int) error {
		return r0
	})
}

func (f *DBStoreAddUploadPartFunc[T]) nextHook() func(context.Context, int, int) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreAddUploadPartFunc[T]) appendCall(r0 DBStoreAddUploadPartFuncCall[T]) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreAddUploadPartFuncCall objects
// describing the invocations of this function.
func (f *DBStoreAddUploadPartFunc[T]) History() []DBStoreAddUploadPartFuncCall[T] {
	f.mutex.Lock()
	history := make([]DBStoreAddUploadPartFuncCall[T], len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreAddUploadPartFuncCall is an object that describes an invocation of
// method AddUploadPart on an instance of MockDBStore.
type DBStoreAddUploadPartFuncCall[T interface{}] struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreAddUploadPartFuncCall[T]) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreAddUploadPartFuncCall[T]) Results() []interface{} {
	return []interface{}{c.Result0}
}

// DBStoreGetUploadByIDFunc describes the behavior when the GetUploadByID
// method of the parent MockDBStore instance is invoked.
type DBStoreGetUploadByIDFunc[T interface{}] struct {
	defaultHook func(context.Context, int) (uploadhandler.Upload[T], bool, error)
	hooks       []func(context.Context, int) (uploadhandler.Upload[T], bool, error)
	history     []DBStoreGetUploadByIDFuncCall[T]
	mutex       sync.Mutex
}

// GetUploadByID delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockDBStore[T]) GetUploadByID(v0 context.Context, v1 int) (uploadhandler.Upload[T], bool, error) {
	r0, r1, r2 := m.GetUploadByIDFunc.nextHook()(v0, v1)
	m.GetUploadByIDFunc.appendCall(DBStoreGetUploadByIDFuncCall[T]{v0, v1, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the GetUploadByID method
// of the parent MockDBStore instance is invoked and the hook queue is
// empty.
func (f *DBStoreGetUploadByIDFunc[T]) SetDefaultHook(hook func(context.Context, int) (uploadhandler.Upload[T], bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetUploadByID method of the parent MockDBStore instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *DBStoreGetUploadByIDFunc[T]) PushHook(hook func(context.Context, int) (uploadhandler.Upload[T], bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DBStoreGetUploadByIDFunc[T]) SetDefaultReturn(r0 uploadhandler.Upload[T], r1 bool, r2 error) {
	f.SetDefaultHook(func(context.Context, int) (uploadhandler.Upload[T], bool, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DBStoreGetUploadByIDFunc[T]) PushReturn(r0 uploadhandler.Upload[T], r1 bool, r2 error) {
	f.PushHook(func(context.Context, int) (uploadhandler.Upload[T], bool, error) {
		return r0, r1, r2
	})
}

func (f *DBStoreGetUploadByIDFunc[T]) nextHook() func(context.Context, int) (uploadhandler.Upload[T], bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreGetUploadByIDFunc[T]) appendCall(r0 DBStoreGetUploadByIDFuncCall[T]) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreGetUploadByIDFuncCall objects
// describing the invocations of this function.
func (f *DBStoreGetUploadByIDFunc[T]) History() []DBStoreGetUploadByIDFuncCall[T] {
	f.mutex.Lock()
	history := make([]DBStoreGetUploadByIDFuncCall[T], len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreGetUploadByIDFuncCall is an object that describes an invocation of
// method GetUploadByID on an instance of MockDBStore.
type DBStoreGetUploadByIDFuncCall[T interface{}] struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 uploadhandler.Upload[T]
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 bool
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreGetUploadByIDFuncCall[T]) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreGetUploadByIDFuncCall[T]) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}

// DBStoreInsertUploadFunc describes the behavior when the InsertUpload
// method of the parent MockDBStore instance is invoked.
type DBStoreInsertUploadFunc[T interface{}] struct {
	defaultHook func(context.Context, uploadhandler.Upload[T]) (int, error)
	hooks       []func(context.Context, uploadhandler.Upload[T]) (int, error)
	history     []DBStoreInsertUploadFuncCall[T]
	mutex       sync.Mutex
}

// InsertUpload delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockDBStore[T]) InsertUpload(v0 context.Context, v1 uploadhandler.Upload[T]) (int, error) {
	r0, r1 := m.InsertUploadFunc.nextHook()(v0, v1)
	m.InsertUploadFunc.appendCall(DBStoreInsertUploadFuncCall[T]{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the InsertUpload method
// of the parent MockDBStore instance is invoked and the hook queue is
// empty.
func (f *DBStoreInsertUploadFunc[T]) SetDefaultHook(hook func(context.Context, uploadhandler.Upload[T]) (int, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// InsertUpload method of the parent MockDBStore instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *DBStoreInsertUploadFunc[T]) PushHook(hook func(context.Context, uploadhandler.Upload[T]) (int, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DBStoreInsertUploadFunc[T]) SetDefaultReturn(r0 int, r1 error) {
	f.SetDefaultHook(func(context.Context, uploadhandler.Upload[T]) (int, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DBStoreInsertUploadFunc[T]) PushReturn(r0 int, r1 error) {
	f.PushHook(func(context.Context, uploadhandler.Upload[T]) (int, error) {
		return r0, r1
	})
}

func (f *DBStoreInsertUploadFunc[T]) nextHook() func(context.Context, uploadhandler.Upload[T]) (int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreInsertUploadFunc[T]) appendCall(r0 DBStoreInsertUploadFuncCall[T]) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreInsertUploadFuncCall objects
// describing the invocations of this function.
func (f *DBStoreInsertUploadFunc[T]) History() []DBStoreInsertUploadFuncCall[T] {
	f.mutex.Lock()
	history := make([]DBStoreInsertUploadFuncCall[T], len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreInsertUploadFuncCall is an object that describes an invocation of
// method InsertUpload on an instance of MockDBStore.
type DBStoreInsertUploadFuncCall[T interface{}] struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 uploadhandler.Upload[T]
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 int
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreInsertUploadFuncCall[T]) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreInsertUploadFuncCall[T]) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// DBStoreMarkFailedFunc describes the behavior when the MarkFailed method
// of the parent MockDBStore instance is invoked.
type DBStoreMarkFailedFunc[T interface{}] struct {
	defaultHook func(context.Context, int, string) error
	hooks       []func(context.Context, int, string) error
	history     []DBStoreMarkFailedFuncCall[T]
	mutex       sync.Mutex
}

// MarkFailed delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockDBStore[T]) MarkFailed(v0 context.Context, v1 int, v2 string) error {
	r0 := m.MarkFailedFunc.nextHook()(v0, v1, v2)
	m.MarkFailedFunc.appendCall(DBStoreMarkFailedFuncCall[T]{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the MarkFailed method of
// the parent MockDBStore instance is invoked and the hook queue is empty.
func (f *DBStoreMarkFailedFunc[T]) SetDefaultHook(hook func(context.Context, int, string) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// MarkFailed method of the parent MockDBStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *DBStoreMarkFailedFunc[T]) PushHook(hook func(context.Context, int, string) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DBStoreMarkFailedFunc[T]) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, int, string) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DBStoreMarkFailedFunc[T]) PushReturn(r0 error) {
	f.PushHook(func(context.Context, int, string) error {
		return r0
	})
}

func (f *DBStoreMarkFailedFunc[T]) nextHook() func(context.Context, int, string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreMarkFailedFunc[T]) appendCall(r0 DBStoreMarkFailedFuncCall[T]) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreMarkFailedFuncCall objects
// describing the invocations of this function.
func (f *DBStoreMarkFailedFunc[T]) History() []DBStoreMarkFailedFuncCall[T] {
	f.mutex.Lock()
	history := make([]DBStoreMarkFailedFuncCall[T], len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreMarkFailedFuncCall is an object that describes an invocation of
// method MarkFailed on an instance of MockDBStore.
type DBStoreMarkFailedFuncCall[T interface{}] struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreMarkFailedFuncCall[T]) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreMarkFailedFuncCall[T]) Results() []interface{} {
	return []interface{}{c.Result0}
}

// DBStoreMarkQueuedFunc describes the behavior when the MarkQueued method
// of the parent MockDBStore instance is invoked.
type DBStoreMarkQueuedFunc[T interface{}] struct {
	defaultHook func(context.Context, int, *int64) error
	hooks       []func(context.Context, int, *int64) error
	history     []DBStoreMarkQueuedFuncCall[T]
	mutex       sync.Mutex
}

// MarkQueued delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockDBStore[T]) MarkQueued(v0 context.Context, v1 int, v2 *int64) error {
	r0 := m.MarkQueuedFunc.nextHook()(v0, v1, v2)
	m.MarkQueuedFunc.appendCall(DBStoreMarkQueuedFuncCall[T]{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the MarkQueued method of
// the parent MockDBStore instance is invoked and the hook queue is empty.
func (f *DBStoreMarkQueuedFunc[T]) SetDefaultHook(hook func(context.Context, int, *int64) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// MarkQueued method of the parent MockDBStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *DBStoreMarkQueuedFunc[T]) PushHook(hook func(context.Context, int, *int64) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DBStoreMarkQueuedFunc[T]) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, int, *int64) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DBStoreMarkQueuedFunc[T]) PushReturn(r0 error) {
	f.PushHook(func(context.Context, int, *int64) error {
		return r0
	})
}

func (f *DBStoreMarkQueuedFunc[T]) nextHook() func(context.Context, int, *int64) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreMarkQueuedFunc[T]) appendCall(r0 DBStoreMarkQueuedFuncCall[T]) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreMarkQueuedFuncCall objects
// describing the invocations of this function.
func (f *DBStoreMarkQueuedFunc[T]) History() []DBStoreMarkQueuedFuncCall[T] {
	f.mutex.Lock()
	history := make([]DBStoreMarkQueuedFuncCall[T], len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreMarkQueuedFuncCall is an object that describes an invocation of
// method MarkQueued on an instance of MockDBStore.
type DBStoreMarkQueuedFuncCall[T interface{}] struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 *int64
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreMarkQueuedFuncCall[T]) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreMarkQueuedFuncCall[T]) Results() []interface{} {
	return []interface{}{c.Result0}
}

// DBStoreWithTransactionFunc describes the behavior when the
// WithTransaction method of the parent MockDBStore instance is invoked.
type DBStoreWithTransactionFunc[T interface{}] struct {
	defaultHook func(context.Context, func(tx uploadhandler.DBStore[T]) error) error
	hooks       []func(context.Context, func(tx uploadhandler.DBStore[T]) error) error
	history     []DBStoreWithTransactionFuncCall[T]
	mutex       sync.Mutex
}

// WithTransaction delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockDBStore[T]) WithTransaction(v0 context.Context, v1 func(tx uploadhandler.DBStore[T]) error) error {
	r0 := m.WithTransactionFunc.nextHook()(v0, v1)
	m.WithTransactionFunc.appendCall(DBStoreWithTransactionFuncCall[T]{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the WithTransaction
// method of the parent MockDBStore instance is invoked and the hook queue
// is empty.
func (f *DBStoreWithTransactionFunc[T]) SetDefaultHook(hook func(context.Context, func(tx uploadhandler.DBStore[T]) error) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// WithTransaction method of the parent MockDBStore instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *DBStoreWithTransactionFunc[T]) PushHook(hook func(context.Context, func(tx uploadhandler.DBStore[T]) error) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DBStoreWithTransactionFunc[T]) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, func(tx uploadhandler.DBStore[T]) error) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DBStoreWithTransactionFunc[T]) PushReturn(r0 error) {
	f.PushHook(func(context.Context, func(tx uploadhandler.DBStore[T]) error) error {
		return r0
	})
}

func (f *DBStoreWithTransactionFunc[T]) nextHook() func(context.Context, func(tx uploadhandler.DBStore[T]) error) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DBStoreWithTransactionFunc[T]) appendCall(r0 DBStoreWithTransactionFuncCall[T]) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DBStoreWithTransactionFuncCall objects
// describing the invocations of this function.
func (f *DBStoreWithTransactionFunc[T]) History() []DBStoreWithTransactionFuncCall[T] {
	f.mutex.Lock()
	history := make([]DBStoreWithTransactionFuncCall[T], len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DBStoreWithTransactionFuncCall is an object that describes an invocation
// of method WithTransaction on an instance of MockDBStore.
type DBStoreWithTransactionFuncCall[T interface{}] struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 func(tx uploadhandler.DBStore[T]) error
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DBStoreWithTransactionFuncCall[T]) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DBStoreWithTransactionFuncCall[T]) Results() []interface{} {
	return []interface{}{c.Result0}
}
