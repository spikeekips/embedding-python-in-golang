package main

/*
#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <pthread.h>
#include <unistd.h>
#include <stdio.h>

extern void createThreadCallback();

static PyObject* thread_callback() {
    PyObject *_module_name, *_module;
	PyGILState_STATE _gstate;

    // Initialize python GIL state
	_gstate = PyGILState_Ensure();

    // Now execute some python code (call python functions)
    _module_name = PyString_FromString("json_dump");
    _module = PyImport_Import(_module_name);

    // Call a method of the class with no parameters
	PyObject *_attr, *_result;
    _attr = PyObject_GetAttr(_module, PyString_FromString("run"));
    _result = PyObject_CallObject(_attr, NULL);

    // Clean up
    Py_DECREF(_module);
    Py_DECREF(_module_name);
    Py_DECREF(_attr);

	PyGILState_Release(_gstate);

	return _result;
}

static void createThread(pthread_t* pid) {
	pthread_create(pid, NULL, (void*)createThreadCallback, pid);
}

static void initialize_python () {

	if (Py_IsInitialized() == 0) {
		Py_Initialize();
		//fprintf(stdout, "> Py_Initialize\n");
	}

	// make sure the GIL is correctly initialized
	if (PyEval_ThreadsInitialized() == 0) {
		PyEval_InitThreads();
		//fprintf(stdout, "> PyEval_ThreadsInitialized\n");
	}

	PyEval_ReleaseThread(PyGILState_GetThisThreadState());
}

*/
import "C"

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"sync"
	//"time"
	//"runtime"
)

var callbacks map[*C.pthread_t]ThreadCallback

type ThreadCallback func(a *C.PyObject)

var lock sync.Mutex

//export createThreadCallback
func createThreadCallback(pid *C.pthread_t) {
	lock.Lock()
	defer lock.Unlock()

	_cb, _ok := callbacks[pid]
	if !_ok {
		panic(fmt.Errorf("failed to found thread callback for `%v`", pid))
	}

	//runtime.LockOSThread()
	_result := C.thread_callback()

	_cb(_result)
	//defer runtime.UnlockOSThread()

	defer func() {
		delete(callbacks, pid)
	}()
}

func create_thread(cb ThreadCallback) {
	_pid := new(C.pthread_t)
	log.Debug("> create_thread: %v", _pid)
	callbacks[_pid] = cb
	C.createThread(_pid)
	log.Debug("< create_thread: %v", _pid)
}

var log = logging.MustGetLogger("gopython-test")

func init() {
	logging.SetLevel(logging.INFO, "gopython-test")
	logging.SetLevel(logging.DEBUG, "gopython-test")

	log.Debug("> Initilize Python.")
	callbacks = map[*C.pthread_t]ThreadCallback{}

	C.initialize_python()
}

func main() {
	log.Info("python embedding test in golang, using `pthreads`.")

	n := 400
	var _wg sync.WaitGroup
	_wg.Add(n)

	for i := 0; i < n; i++ {
		go create_thread(func(result *C.PyObject) {
			defer _wg.Done()

			_result_string := C.GoString(C.PyString_AsString(result))
			log.Debug("< got result string: %v (%T)", _result_string, _result_string)

			var _parsed []interface{}
			if _err := json.Unmarshal([]byte(_result_string), &_parsed); _err != nil {
				panic(fmt.Errorf("got invalid result from python function, `%v`", _result_string))
			}

			log.Debug(
				"< got thread_id=%v\tnow=%s",
				int64(_parsed[0].(float64)),
				_parsed[1].(string),
			)
		})
	}
	_wg.Wait()
}
