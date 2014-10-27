package pyingo

/*
#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <stdio.h>

static void initialize () {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
        //fprintf(stdout, "`Py_Initialize`\n");
    }

    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
        //fprintf(stdout, "`PyEval_InitThreads`\n");
        PyEval_ReleaseThread(PyGILState_GetThisThreadState());
        //fprintf(stdout, "`PyEval_ReleaseThread`\n");
    }
}

static PyGILState_STATE thread_begin() {
    //fprintf(stdout, "`PyGILState_Ensure`\n");
    return PyGILState_Ensure();
}

static void thread_end(PyGILState_STATE gstate) {
    //fprintf(stdout, "`PyGILState_Release`\n");
    PyGILState_Release(gstate);
}

static PyObject* run_python(char *module_name, char *function_name) {
    // Now execute some python code (call python functions)
    PyObject *_module;
    _module = PyImport_Import(PyString_FromString(module_name));
    if (_module == NULL) {
        fprintf(stderr, "failed to import python module, `%s`.\n", module_name);
        return NULL;
    }

    // Call a method of the class with no parameters
    PyObject *_attr;
    _attr = PyObject_GetAttr(_module, PyString_FromString(function_name));

    PyObject *_args = NULL;
    _args = Py_BuildValue("()");

    PyObject* _result = PyObject_CallObject(_attr, _args);

    // Clean up
    Py_DECREF(_module);
    Py_DECREF(_attr);
    Py_DECREF(_args);

    return _result;
}

*/
import "C"

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

func init() {
	log.Debug("initialize python.")
	C.initialize()
}

type PyObject struct {
	Ptr *C.PyObject
}

type Python struct {
	threadpool *ThreadPool
}

func NewPython(number_of_threads int) *Python {
	py := new(Python)
	py.threadpool = NewThreadPool(number_of_threads)

	return py
}

func (self *Python) run_python(py_name string) (*PyObject, error) {
	if len(strings.Trim(py_name, "")) < 1 {
		return nil, fmt.Errorf("`py_name` is empty.")
	}

	_p := strings.Split(py_name, ":")
	if len(_p) != 2 {
		return nil, fmt.Errorf("`py_name` must be `<module>:<function>`.")
	}

	_module_name, _function_name := _p[0], _p[1]
	log.Debug("`%s.%s(...)` will be called.", _module_name, _function_name)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_result := new(C.PyObject)

	_tid := self.threadpool.Start(
		ThreadCallback(func(args unsafe.Pointer) {
			_gstate := C.thread_begin()
			defer C.thread_end(_gstate)

			_result = C.run_python(
				C.CString(_module_name),
				C.CString(_function_name),
			)

			log.Debug("in callback: %s, result: %v", py_name, _result)
		}),
	)
	self.threadpool.Join(_tid)

	return &PyObject{Ptr: _result}, nil
}

func (self *Python) Run(py_name string) (*PyObject, error) {
	return self.run_python(py_name)
}

var lock sync.Mutex

func (self *Python) RunInGoRoutine(py_name string) (*PyObject, error) {
	lock.Lock()
	defer lock.Unlock()

	return self.run_python(py_name)
}

func PyString_FromString(s string) *PyObject {
	_s := C.CString(s)
	defer C.free(unsafe.Pointer(_s))

	return &PyObject{Ptr: C.PyString_FromString(_s)}
}

func PyString_AsString(s *PyObject) string {
	return C.GoString(C.PyString_AsString(s.Ptr))
}
