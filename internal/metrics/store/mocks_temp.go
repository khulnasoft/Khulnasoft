// Code generated by go-mockgen 1.3.7; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package store

import (
	"sync"

	go1 "github.com/prometheus/client_model/go"
)

// MockDistributedStore is a mock implementation of the DistributedStore
// interface (from the package
// github.com/khulnasoft/khulnasoft/internal/metrics/store) used for unit
// testing.
type MockDistributedStore struct {
	// GatherFunc is an instance of a mock function object controlling the
	// behavior of the method Gather.
	GatherFunc *DistributedStoreGatherFunc
	// IngestFunc is an instance of a mock function object controlling the
	// behavior of the method Ingest.
	IngestFunc *DistributedStoreIngestFunc
}

// NewMockDistributedStore creates a new mock of the DistributedStore
// interface. All methods return zero values for all results, unless
// overwritten.
func NewMockDistributedStore() *MockDistributedStore {
	return &MockDistributedStore{
		GatherFunc: &DistributedStoreGatherFunc{
			defaultHook: func() (r0 []*go1.MetricFamily, r1 error) {
				return
			},
		},
		IngestFunc: &DistributedStoreIngestFunc{
			defaultHook: func(string, []*go1.MetricFamily) (r0 error) {
				return
			},
		},
	}
}

// NewStrictMockDistributedStore creates a new mock of the DistributedStore
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockDistributedStore() *MockDistributedStore {
	return &MockDistributedStore{
		GatherFunc: &DistributedStoreGatherFunc{
			defaultHook: func() ([]*go1.MetricFamily, error) {
				panic("unexpected invocation of MockDistributedStore.Gather")
			},
		},
		IngestFunc: &DistributedStoreIngestFunc{
			defaultHook: func(string, []*go1.MetricFamily) error {
				panic("unexpected invocation of MockDistributedStore.Ingest")
			},
		},
	}
}

// NewMockDistributedStoreFrom creates a new mock of the
// MockDistributedStore interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockDistributedStoreFrom(i DistributedStore) *MockDistributedStore {
	return &MockDistributedStore{
		GatherFunc: &DistributedStoreGatherFunc{
			defaultHook: i.Gather,
		},
		IngestFunc: &DistributedStoreIngestFunc{
			defaultHook: i.Ingest,
		},
	}
}

// DistributedStoreGatherFunc describes the behavior when the Gather method
// of the parent MockDistributedStore instance is invoked.
type DistributedStoreGatherFunc struct {
	defaultHook func() ([]*go1.MetricFamily, error)
	hooks       []func() ([]*go1.MetricFamily, error)
	history     []DistributedStoreGatherFuncCall
	mutex       sync.Mutex
}

// Gather delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockDistributedStore) Gather() ([]*go1.MetricFamily, error) {
	r0, r1 := m.GatherFunc.nextHook()()
	m.GatherFunc.appendCall(DistributedStoreGatherFuncCall{r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Gather method of the
// parent MockDistributedStore instance is invoked and the hook queue is
// empty.
func (f *DistributedStoreGatherFunc) SetDefaultHook(hook func() ([]*go1.MetricFamily, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Gather method of the parent MockDistributedStore instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *DistributedStoreGatherFunc) PushHook(hook func() ([]*go1.MetricFamily, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DistributedStoreGatherFunc) SetDefaultReturn(r0 []*go1.MetricFamily, r1 error) {
	f.SetDefaultHook(func() ([]*go1.MetricFamily, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DistributedStoreGatherFunc) PushReturn(r0 []*go1.MetricFamily, r1 error) {
	f.PushHook(func() ([]*go1.MetricFamily, error) {
		return r0, r1
	})
}

func (f *DistributedStoreGatherFunc) nextHook() func() ([]*go1.MetricFamily, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DistributedStoreGatherFunc) appendCall(r0 DistributedStoreGatherFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DistributedStoreGatherFuncCall objects
// describing the invocations of this function.
func (f *DistributedStoreGatherFunc) History() []DistributedStoreGatherFuncCall {
	f.mutex.Lock()
	history := make([]DistributedStoreGatherFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DistributedStoreGatherFuncCall is an object that describes an invocation
// of method Gather on an instance of MockDistributedStore.
type DistributedStoreGatherFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []*go1.MetricFamily
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DistributedStoreGatherFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DistributedStoreGatherFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// DistributedStoreIngestFunc describes the behavior when the Ingest method
// of the parent MockDistributedStore instance is invoked.
type DistributedStoreIngestFunc struct {
	defaultHook func(string, []*go1.MetricFamily) error
	hooks       []func(string, []*go1.MetricFamily) error
	history     []DistributedStoreIngestFuncCall
	mutex       sync.Mutex
}

// Ingest delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockDistributedStore) Ingest(v0 string, v1 []*go1.MetricFamily) error {
	r0 := m.IngestFunc.nextHook()(v0, v1)
	m.IngestFunc.appendCall(DistributedStoreIngestFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Ingest method of the
// parent MockDistributedStore instance is invoked and the hook queue is
// empty.
func (f *DistributedStoreIngestFunc) SetDefaultHook(hook func(string, []*go1.MetricFamily) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Ingest method of the parent MockDistributedStore instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *DistributedStoreIngestFunc) PushHook(hook func(string, []*go1.MetricFamily) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *DistributedStoreIngestFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(string, []*go1.MetricFamily) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *DistributedStoreIngestFunc) PushReturn(r0 error) {
	f.PushHook(func(string, []*go1.MetricFamily) error {
		return r0
	})
}

func (f *DistributedStoreIngestFunc) nextHook() func(string, []*go1.MetricFamily) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *DistributedStoreIngestFunc) appendCall(r0 DistributedStoreIngestFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of DistributedStoreIngestFuncCall objects
// describing the invocations of this function.
func (f *DistributedStoreIngestFunc) History() []DistributedStoreIngestFuncCall {
	f.mutex.Lock()
	history := make([]DistributedStoreIngestFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// DistributedStoreIngestFuncCall is an object that describes an invocation
// of method Ingest on an instance of MockDistributedStore.
type DistributedStoreIngestFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 string
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 []*go1.MetricFamily
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c DistributedStoreIngestFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c DistributedStoreIngestFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// MockStore is a mock implementation of the Store interface (from the
// package github.com/khulnasoft/khulnasoft/internal/metrics/store) used
// for unit testing.
type MockStore struct {
	// GatherFunc is an instance of a mock function object controlling the
	// behavior of the method Gather.
	GatherFunc *StoreGatherFunc
}

// NewMockStore creates a new mock of the Store interface. All methods
// return zero values for all results, unless overwritten.
func NewMockStore() *MockStore {
	return &MockStore{
		GatherFunc: &StoreGatherFunc{
			defaultHook: func() (r0 []*go1.MetricFamily, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockStore creates a new mock of the Store interface. All methods
// panic on invocation, unless overwritten.
func NewStrictMockStore() *MockStore {
	return &MockStore{
		GatherFunc: &StoreGatherFunc{
			defaultHook: func() ([]*go1.MetricFamily, error) {
				panic("unexpected invocation of MockStore.Gather")
			},
		},
	}
}

// NewMockStoreFrom creates a new mock of the MockStore interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockStoreFrom(i Store) *MockStore {
	return &MockStore{
		GatherFunc: &StoreGatherFunc{
			defaultHook: i.Gather,
		},
	}
}

// StoreGatherFunc describes the behavior when the Gather method of the
// parent MockStore instance is invoked.
type StoreGatherFunc struct {
	defaultHook func() ([]*go1.MetricFamily, error)
	hooks       []func() ([]*go1.MetricFamily, error)
	history     []StoreGatherFuncCall
	mutex       sync.Mutex
}

// Gather delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockStore) Gather() ([]*go1.MetricFamily, error) {
	r0, r1 := m.GatherFunc.nextHook()()
	m.GatherFunc.appendCall(StoreGatherFuncCall{r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Gather method of the
// parent MockStore instance is invoked and the hook queue is empty.
func (f *StoreGatherFunc) SetDefaultHook(hook func() ([]*go1.MetricFamily, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Gather method of the parent MockStore instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *StoreGatherFunc) PushHook(hook func() ([]*go1.MetricFamily, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *StoreGatherFunc) SetDefaultReturn(r0 []*go1.MetricFamily, r1 error) {
	f.SetDefaultHook(func() ([]*go1.MetricFamily, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *StoreGatherFunc) PushReturn(r0 []*go1.MetricFamily, r1 error) {
	f.PushHook(func() ([]*go1.MetricFamily, error) {
		return r0, r1
	})
}

func (f *StoreGatherFunc) nextHook() func() ([]*go1.MetricFamily, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *StoreGatherFunc) appendCall(r0 StoreGatherFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of StoreGatherFuncCall objects describing the
// invocations of this function.
func (f *StoreGatherFunc) History() []StoreGatherFuncCall {
	f.mutex.Lock()
	history := make([]StoreGatherFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// StoreGatherFuncCall is an object that describes an invocation of method
// Gather on an instance of MockStore.
type StoreGatherFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []*go1.MetricFamily
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c StoreGatherFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c StoreGatherFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
