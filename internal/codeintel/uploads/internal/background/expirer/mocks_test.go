// Code generated by go-mockgen 1.3.7; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package expirer

import (
	"context"
	"sync"
	"time"

	api "github.com/khulnasoft/khulnasoft/internal/api"
	policies "github.com/khulnasoft/khulnasoft/internal/codeintel/policies"
	shared "github.com/khulnasoft/khulnasoft/internal/codeintel/policies/shared"
)

// MockPolicyMatcher is a mock implementation of the PolicyMatcher interface
// (from the package
// github.com/khulnasoft/khulnasoft/internal/codeintel/uploads/internal/background/expirer)
// used for unit testing.
type MockPolicyMatcher struct {
	// CommitsDescribedByPolicyFunc is an instance of a mock function object
	// controlling the behavior of the method CommitsDescribedByPolicy.
	CommitsDescribedByPolicyFunc *PolicyMatcherCommitsDescribedByPolicyFunc
}

// NewMockPolicyMatcher creates a new mock of the PolicyMatcher interface.
// All methods return zero values for all results, unless overwritten.
func NewMockPolicyMatcher() *MockPolicyMatcher {
	return &MockPolicyMatcher{
		CommitsDescribedByPolicyFunc: &PolicyMatcherCommitsDescribedByPolicyFunc{
			defaultHook: func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (r0 map[string][]policies.PolicyMatch, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockPolicyMatcher creates a new mock of the PolicyMatcher
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockPolicyMatcher() *MockPolicyMatcher {
	return &MockPolicyMatcher{
		CommitsDescribedByPolicyFunc: &PolicyMatcherCommitsDescribedByPolicyFunc{
			defaultHook: func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error) {
				panic("unexpected invocation of MockPolicyMatcher.CommitsDescribedByPolicy")
			},
		},
	}
}

// NewMockPolicyMatcherFrom creates a new mock of the MockPolicyMatcher
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockPolicyMatcherFrom(i PolicyMatcher) *MockPolicyMatcher {
	return &MockPolicyMatcher{
		CommitsDescribedByPolicyFunc: &PolicyMatcherCommitsDescribedByPolicyFunc{
			defaultHook: i.CommitsDescribedByPolicy,
		},
	}
}

// PolicyMatcherCommitsDescribedByPolicyFunc describes the behavior when the
// CommitsDescribedByPolicy method of the parent MockPolicyMatcher instance
// is invoked.
type PolicyMatcherCommitsDescribedByPolicyFunc struct {
	defaultHook func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error)
	hooks       []func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error)
	history     []PolicyMatcherCommitsDescribedByPolicyFuncCall
	mutex       sync.Mutex
}

// CommitsDescribedByPolicy delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockPolicyMatcher) CommitsDescribedByPolicy(v0 context.Context, v1 int, v2 api.RepoName, v3 []shared.ConfigurationPolicy, v4 time.Time, v5 ...string) (map[string][]policies.PolicyMatch, error) {
	r0, r1 := m.CommitsDescribedByPolicyFunc.nextHook()(v0, v1, v2, v3, v4, v5...)
	m.CommitsDescribedByPolicyFunc.appendCall(PolicyMatcherCommitsDescribedByPolicyFuncCall{v0, v1, v2, v3, v4, v5, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// CommitsDescribedByPolicy method of the parent MockPolicyMatcher instance
// is invoked and the hook queue is empty.
func (f *PolicyMatcherCommitsDescribedByPolicyFunc) SetDefaultHook(hook func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// CommitsDescribedByPolicy method of the parent MockPolicyMatcher instance
// invokes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *PolicyMatcherCommitsDescribedByPolicyFunc) PushHook(hook func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *PolicyMatcherCommitsDescribedByPolicyFunc) SetDefaultReturn(r0 map[string][]policies.PolicyMatch, r1 error) {
	f.SetDefaultHook(func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *PolicyMatcherCommitsDescribedByPolicyFunc) PushReturn(r0 map[string][]policies.PolicyMatch, r1 error) {
	f.PushHook(func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error) {
		return r0, r1
	})
}

func (f *PolicyMatcherCommitsDescribedByPolicyFunc) nextHook() func(context.Context, int, api.RepoName, []shared.ConfigurationPolicy, time.Time, ...string) (map[string][]policies.PolicyMatch, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *PolicyMatcherCommitsDescribedByPolicyFunc) appendCall(r0 PolicyMatcherCommitsDescribedByPolicyFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// PolicyMatcherCommitsDescribedByPolicyFuncCall objects describing the
// invocations of this function.
func (f *PolicyMatcherCommitsDescribedByPolicyFunc) History() []PolicyMatcherCommitsDescribedByPolicyFuncCall {
	f.mutex.Lock()
	history := make([]PolicyMatcherCommitsDescribedByPolicyFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// PolicyMatcherCommitsDescribedByPolicyFuncCall is an object that describes
// an invocation of method CommitsDescribedByPolicy on an instance of
// MockPolicyMatcher.
type PolicyMatcherCommitsDescribedByPolicyFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 api.RepoName
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 []shared.ConfigurationPolicy
	// Arg4 is the value of the 5th argument passed to this method
	// invocation.
	Arg4 time.Time
	// Arg5 is a slice containing the values of the variadic arguments
	// passed to this method invocation.
	Arg5 []string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 map[string][]policies.PolicyMatch
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation. The variadic slice argument is flattened in this array such
// that one positional argument and three variadic arguments would result in
// a slice of four, not two.
func (c PolicyMatcherCommitsDescribedByPolicyFuncCall) Args() []interface{} {
	trailing := []interface{}{}
	for _, val := range c.Arg5 {
		trailing = append(trailing, val)
	}

	return append([]interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3, c.Arg4}, trailing...)
}

// Results returns an interface slice containing the results of this
// invocation.
func (c PolicyMatcherCommitsDescribedByPolicyFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// MockPolicyService is a mock implementation of the PolicyService interface
// (from the package
// github.com/khulnasoft/khulnasoft/internal/codeintel/uploads/internal/background/expirer)
// used for unit testing.
type MockPolicyService struct {
	// GetConfigurationPoliciesFunc is an instance of a mock function object
	// controlling the behavior of the method GetConfigurationPolicies.
	GetConfigurationPoliciesFunc *PolicyServiceGetConfigurationPoliciesFunc
}

// NewMockPolicyService creates a new mock of the PolicyService interface.
// All methods return zero values for all results, unless overwritten.
func NewMockPolicyService() *MockPolicyService {
	return &MockPolicyService{
		GetConfigurationPoliciesFunc: &PolicyServiceGetConfigurationPoliciesFunc{
			defaultHook: func(context.Context, shared.GetConfigurationPoliciesOptions) (r0 []shared.ConfigurationPolicy, r1 int, r2 error) {
				return
			},
		},
	}
}

// NewStrictMockPolicyService creates a new mock of the PolicyService
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockPolicyService() *MockPolicyService {
	return &MockPolicyService{
		GetConfigurationPoliciesFunc: &PolicyServiceGetConfigurationPoliciesFunc{
			defaultHook: func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error) {
				panic("unexpected invocation of MockPolicyService.GetConfigurationPolicies")
			},
		},
	}
}

// NewMockPolicyServiceFrom creates a new mock of the MockPolicyService
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockPolicyServiceFrom(i PolicyService) *MockPolicyService {
	return &MockPolicyService{
		GetConfigurationPoliciesFunc: &PolicyServiceGetConfigurationPoliciesFunc{
			defaultHook: i.GetConfigurationPolicies,
		},
	}
}

// PolicyServiceGetConfigurationPoliciesFunc describes the behavior when the
// GetConfigurationPolicies method of the parent MockPolicyService instance
// is invoked.
type PolicyServiceGetConfigurationPoliciesFunc struct {
	defaultHook func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error)
	hooks       []func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error)
	history     []PolicyServiceGetConfigurationPoliciesFuncCall
	mutex       sync.Mutex
}

// GetConfigurationPolicies delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockPolicyService) GetConfigurationPolicies(v0 context.Context, v1 shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error) {
	r0, r1, r2 := m.GetConfigurationPoliciesFunc.nextHook()(v0, v1)
	m.GetConfigurationPoliciesFunc.appendCall(PolicyServiceGetConfigurationPoliciesFuncCall{v0, v1, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the
// GetConfigurationPolicies method of the parent MockPolicyService instance
// is invoked and the hook queue is empty.
func (f *PolicyServiceGetConfigurationPoliciesFunc) SetDefaultHook(hook func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetConfigurationPolicies method of the parent MockPolicyService instance
// invokes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *PolicyServiceGetConfigurationPoliciesFunc) PushHook(hook func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *PolicyServiceGetConfigurationPoliciesFunc) SetDefaultReturn(r0 []shared.ConfigurationPolicy, r1 int, r2 error) {
	f.SetDefaultHook(func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *PolicyServiceGetConfigurationPoliciesFunc) PushReturn(r0 []shared.ConfigurationPolicy, r1 int, r2 error) {
	f.PushHook(func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error) {
		return r0, r1, r2
	})
}

func (f *PolicyServiceGetConfigurationPoliciesFunc) nextHook() func(context.Context, shared.GetConfigurationPoliciesOptions) ([]shared.ConfigurationPolicy, int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *PolicyServiceGetConfigurationPoliciesFunc) appendCall(r0 PolicyServiceGetConfigurationPoliciesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// PolicyServiceGetConfigurationPoliciesFuncCall objects describing the
// invocations of this function.
func (f *PolicyServiceGetConfigurationPoliciesFunc) History() []PolicyServiceGetConfigurationPoliciesFuncCall {
	f.mutex.Lock()
	history := make([]PolicyServiceGetConfigurationPoliciesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// PolicyServiceGetConfigurationPoliciesFuncCall is an object that describes
// an invocation of method GetConfigurationPolicies on an instance of
// MockPolicyService.
type PolicyServiceGetConfigurationPoliciesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 shared.GetConfigurationPoliciesOptions
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []shared.ConfigurationPolicy
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 int
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c PolicyServiceGetConfigurationPoliciesFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c PolicyServiceGetConfigurationPoliciesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}
