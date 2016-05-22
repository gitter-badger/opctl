// This file was generated by counterfeiter
package core

import (
	"sync"

	"github.com/opctl/engine/core/models"
)

type fakeAddOpUseCase struct {
	ExecuteStub        func(req models.AddOpReq) (err error)
	executeMutex       sync.RWMutex
	executeArgsForCall []struct {
		req models.AddOpReq
	}
	executeReturns struct {
		result1 error
	}
}

func (fake *fakeAddOpUseCase) Execute(req models.AddOpReq) (err error) {
	fake.executeMutex.Lock()
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct {
		req models.AddOpReq
	}{req})
	fake.executeMutex.Unlock()
	if fake.ExecuteStub != nil {
		return fake.ExecuteStub(req)
	} else {
		return fake.executeReturns.result1
	}
}

func (fake *fakeAddOpUseCase) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *fakeAddOpUseCase) ExecuteArgsForCall(i int) models.AddOpReq {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return fake.executeArgsForCall[i].req
}

func (fake *fakeAddOpUseCase) ExecuteReturns(result1 error) {
	fake.ExecuteStub = nil
	fake.executeReturns = struct {
		result1 error
	}{result1}
}
