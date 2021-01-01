// Code generated by go-mockgen 0.1.0; DO NOT EDIT.

package mocks

import (
	"context"
	process "github.com/go-nacelle/process"
	"sync"
)

// MockFinalizer is a mock implementation of the Finalizer interface (from
// the package github.com/go-nacelle/process) used for unit testing.
type MockFinalizer struct {
	// FinalizeFunc is an instance of a mock function object controlling the
	// behavior of the method Finalize.
	FinalizeFunc *FinalizerFinalizeFunc
}

// NewMockFinalizer creates a new mock of the Finalizer interface. All
// methods return zero values for all results, unless overwritten.
func NewMockFinalizer() *MockFinalizer {
	return &MockFinalizer{
		FinalizeFunc: &FinalizerFinalizeFunc{
			defaultHook: func(context.Context) error {
				return nil
			},
		},
	}
}

// NewMockFinalizerFrom creates a new mock of the MockFinalizer interface.
// All methods delegate to the given implementation, unless overwritten.
func NewMockFinalizerFrom(i process.Finalizer) *MockFinalizer {
	return &MockFinalizer{
		FinalizeFunc: &FinalizerFinalizeFunc{
			defaultHook: i.Finalize,
		},
	}
}

// FinalizerFinalizeFunc describes the behavior when the Finalize method of
// the parent MockFinalizer instance is invoked.
type FinalizerFinalizeFunc struct {
	defaultHook func(context.Context) error
	hooks       []func(context.Context) error
	history     []FinalizerFinalizeFuncCall
	mutex       sync.Mutex
}

// Finalize delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockFinalizer) Finalize(v0 context.Context) error {
	r0 := m.FinalizeFunc.nextHook()(v0)
	m.FinalizeFunc.appendCall(FinalizerFinalizeFuncCall{v0, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Finalize method of
// the parent MockFinalizer instance is invoked and the hook queue is empty.
func (f *FinalizerFinalizeFunc) SetDefaultHook(hook func(context.Context) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Finalize method of the parent MockFinalizer instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *FinalizerFinalizeFunc) PushHook(hook func(context.Context) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *FinalizerFinalizeFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *FinalizerFinalizeFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context) error {
		return r0
	})
}

func (f *FinalizerFinalizeFunc) nextHook() func(context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *FinalizerFinalizeFunc) appendCall(r0 FinalizerFinalizeFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of FinalizerFinalizeFuncCall objects
// describing the invocations of this function.
func (f *FinalizerFinalizeFunc) History() []FinalizerFinalizeFuncCall {
	f.mutex.Lock()
	history := make([]FinalizerFinalizeFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// FinalizerFinalizeFuncCall is an object that describes an invocation of
// method Finalize on an instance of MockFinalizer.
type FinalizerFinalizeFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c FinalizerFinalizeFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c FinalizerFinalizeFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}
