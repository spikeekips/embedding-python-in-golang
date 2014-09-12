package main

// #cgo pkg-config: python-2.7
// #include "Python.h"
// #include <stdlib.h>
// #include <string.h>
// #include <pthread.h>
// #include <signal.h>
import "C"

import (
	"fmt"
	"github.com/liamzdenek/go-pthreads"
	"github.com/op/go-logging"
	"github.com/sbinet/go-python"
	"github.com/spikeekips/embedding-python-in-golang/common"
	"runtime"
	"sync"
)

func init() {
	log.Debug("> Initilize Python.")

	runtime.LockOSThread()

	if C.Py_IsInitialized() == 0 {
		C.Py_Initialize()
	}
	if C.Py_IsInitialized() == 0 {
		panic(fmt.Errorf("python: could not initialize the python interpreter"))
	}

	// make sure the GIL is correctly initialized
	if C.PyEval_ThreadsInitialized() == 0 {
		C.PyEval_InitThreads()
	}
	if C.PyEval_ThreadsInitialized() == 0 {
		panic(fmt.Errorf("python: could not initialize the GIL"))
	}
	log.Debug("< Initilized Python.")
	_tstate := C.PyGILState_GetThisThreadState()
	C.PyEval_ReleaseThread(_tstate)
}

func create_thread(ch chan bool, args []string, kw map[string]string) {
	log.Debug("> create_thread")

	lock.Lock()
	defer lock.Unlock()
	defer func() {
		ch <- true
	}()

	_ch := make(chan bool)
	defer func() { close(_ch) }()

	thread := pthread.Create(func() {
		_gstate := C.PyGILState_Ensure()
		defer func() {
			C.PyGILState_Release(_gstate)
		}()
		embed_function(args, kw)
		_ch <- true
	})

	defer func() {
		C.pthread_cancel((C.pthread_t)(thread))
		C.pthread_kill((C.pthread_t)(thread), C.SIGSEGV)
	}()

	<-_ch
	log.Debug("< create_thread: done")
}

func embed_function(args []string, kw map[string]string) {
	// import python module
	_module := python.PyImport_ImportModuleNoBlock("json_dump")

	// get `run` attribute from
	_attr := _module.GetAttrString("run")

	// convert golang arguments to python argument.
	_a := python.PyTuple_New(len(args))
	for i, v := range args {
		python.PyTuple_SET_ITEM(_a, i, python.PyString_FromString(v))

	}

	_kw := python.PyDict_New()
	for k, v := range kw {
		python.PyDict_SetItem(
			_kw, python.PyString_FromString(k),
			python.PyString_FromString(v),
		)
	}

	// pack arguments
	_args := python.PyTuple_New(2)
	python.PyTuple_SET_ITEM(_args, 0, _a)
	python.PyTuple_SET_ITEM(_args, 1, _kw)

	// call python function, `pymodule.run`
	_result := _attr.CallObject(_args)

	// the returned value from python will be json string, so convert the
	// string from python to golang string.
	_result_string := python.PyString_AsString(_result)

	log.Debug(
		"got json string: \n%s\n",
		epig.PrettyPrintJson(_result_string),
	)
}

var log = logging.MustGetLogger("gopython-test")
var lock sync.Mutex

func main() {
	logging.SetLevel(logging.DEBUG, "gopython-test")

	log.Info("python embedding test in golang, using `pthreads`.")

	// create argument
	_args := []string{"oldboy", "나는 오대수"}

	_kw := map[string]string{
		"who are you?":                   "누구냐, 넌?",
		"I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
	}

	var _ch chan bool

	// inside goroutine
	_ch = make(chan bool)
	log.Info("run inside goroutine")
	go func() {
		create_thread(_ch, _args, _kw)
	}()
	go func() {
		create_thread(_ch, _args, _kw)
	}()
	<-_ch
	<-_ch
	close(_ch)

	// outside goroutine
	_ch = make(chan bool, 2)
	log.Info("run outside goroutine")
	create_thread(_ch, _args, _kw)
	<-_ch
	close(_ch)
}
