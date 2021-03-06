// This file was generated by counterfeiter
package uniquestring

import (
	"sync"
)

type Fake struct {
	ConstructStub        func() (uniqueString string)
	constructMutex       sync.RWMutex
	constructArgsForCall []struct{}
	constructReturns     struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Fake) Construct() (uniqueString string) {
	fake.constructMutex.Lock()
	fake.constructArgsForCall = append(fake.constructArgsForCall, struct{}{})
	fake.recordInvocation("Construct", []interface{}{})
	fake.constructMutex.Unlock()
	if fake.ConstructStub != nil {
		return fake.ConstructStub()
	} else {
		return fake.constructReturns.result1
	}
}

func (fake *Fake) ConstructCallCount() int {
	fake.constructMutex.RLock()
	defer fake.constructMutex.RUnlock()
	return len(fake.constructArgsForCall)
}

func (fake *Fake) ConstructReturns(result1 string) {
	fake.ConstructStub = nil
	fake.constructReturns = struct {
		result1 string
	}{result1}
}

func (fake *Fake) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.constructMutex.RLock()
	defer fake.constructMutex.RUnlock()
	return fake.invocations
}

func (fake *Fake) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ UniqueStringFactory = new(Fake)
