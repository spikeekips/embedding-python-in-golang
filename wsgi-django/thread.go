package wsgi

/*

#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <pthread.h>
#include <unistd.h>
#include <stdio.h>

static void initialize_thread () {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
        fprintf(stdout, "> `Py_Initialize`\n");
    }

    // make sure the GIL is correctly initialized
    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
        fprintf(stdout, "> `PyEval_InitThreads`\n");
        PyEval_ReleaseThread(PyGILState_GetThisThreadState());
        fprintf(stdout, "> `PyEval_ReleaseThread`\n");
    }
}

extern void createCallback(pthread_t* pid);

static void createThread(pthread_t* pid) {
    pthread_create(pid, NULL, (void*)createCallback, pid);
}

static PyGILState_STATE start_thread() {
    return PyGILState_Ensure();
}

static void end_thread(PyGILState_STATE gstate) {
    PyGILState_Release(gstate);
}

*/
import "C"

import (
	"fmt"
	"github.com/op/go-logging"
	//"runtime"
	"sync"
)

var log_thread = logging.MustGetLogger("thread")
var callbacks Callbacks

func init() {
	log_thread.Debug("> initialize python threads.")
	callbacks = Callbacks{
		function: map[string]func(){},
	}

	C.initialize_thread()
}

func CreateThread(function func()) {
	_pid := new(C.pthread_t)
	callbacks.Add(_pid, function)

	C.createThread(_pid)
}

//export createCallback
func createCallback(pid *C.pthread_t) {
	//runtime.LockOSThread()
	//defer runtime.UnlockOSThread()

	_gstate := C.start_thread()
	defer C.end_thread(_gstate)

	_func, _ok := callbacks.Get(pid)
	defer callbacks.Delete(pid)

	if !_ok {
		panic(fmt.Errorf("failed to found thread callback for `%v`", pid))
	}

	_func()
}

type Callbacks struct {
	sync.RWMutex
	function map[string]func()
}

func (cb *Callbacks) Add(pid *C.pthread_t, function func()) {
	cb.Lock()
	defer cb.Unlock()

	_key := fmt.Sprintf("%p", pid)
	cb.function[_key] = function
}

func (cb *Callbacks) Get(pid *C.pthread_t) (
	func(),
	bool,
) {
	cb.RLock()
	defer cb.RUnlock()

	_function, _ok := cb.function[fmt.Sprintf("%p", pid)]
	if !_ok {
		return nil, false
	}

	return _function, _ok
}

func (cb *Callbacks) Delete(pid *C.pthread_t) {
	cb.Lock()
	defer cb.Unlock()

	_key := fmt.Sprintf("%p", pid)
	delete(cb.function, _key)
}
