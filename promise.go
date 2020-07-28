package promise

import (
	"fmt"
	"sync"
)


type Promise struct {
	run     	func(resolve func(interface{}), reject func(error))
	onPending   bool
	finalResult interface{}
	err         error
	wg          sync.WaitGroup
}

func CreatePromise(run func(resolve func(interface{}), reject func(error))) *Promise {
	var promise = &Promise{
		run:     	 run,
		onPending:   true,
		finalResult: nil,
		err:         nil,
		wg:          sync.WaitGroup{},
	}

	promise.wg.Add(1)

	go func() {
		defer promise.panicHandler()
		promise.run(promise.resolve, promise.reject)
	}()

	return promise
}



func (promise *Promise) resolve(resolution interface{}) {
	if !promise.onPending {
		return
	}

	switch result := resolution.(type) {
	case *Promise:
		resultAwait, err := result.Await()
		if err != nil {
			promise.reject(err)
			return
		}
		promise.finalResult = resultAwait
	default:
		promise.finalResult = result
	}
	promise.onPending = false

	promise.wg.Done()
}

func (promise *Promise) reject(err error) {

	if !promise.onPending {
		return
	}

	promise.err = err
	promise.onPending = false

	promise.wg.Done()
}

func (promise *Promise) panicHandler() {
	e := recover()
	if e != nil {
		switch err := e.(type) {
		case error:
			promise.reject(fmt.Errorf("error: %s", err.Error()))
		case nil:
			promise.reject(fmt.Errorf("nil"))
		default:
			promise.reject(fmt.Errorf("unknown"))
		}
	}
}

func (promise *Promise) Await() (interface{}, error) {
	promise.wg.Wait()
	return promise.finalResult, promise.err
}

func Resolve(resolution interface{}) *Promise {
	return CreatePromise(func(resolve func(interface{}), reject func(error)) {
		resolve(resolution)
	})
}

func Reject(err error) *Promise {
	return CreatePromise(func(resolve func(interface{}), reject func(error)) {
		reject(err)
	})
}

func (promise *Promise) Then(onFulfilled func(info interface{}) interface{}) *Promise {
	return CreatePromise(func(resolve func(interface{}), reject func(error)) {
		result, err := promise.Await()
		if err != nil {
			reject(err)
			return
		}
		resolve(onFulfilled(result))
	})
}

func (promise *Promise) Catch(onRejected func(err error) error) *Promise {
	return CreatePromise(func(resolve func(interface{}), reject func(error)) {
		result, err := promise.Await()
		if err != nil {
			reject(onRejected(err))
			return
		}
		resolve(result)
	})
}


func (promise *Promise) Finally(onFinally func(info interface{}) interface{}) *Promise {
	return CreatePromise(func(resolve func(interface{}), reject func(error)) {
		result, err := promise.Await()
		if err != nil {
			reject(err)
			return
		}
		resolve(onFinally(result))
	})
}


