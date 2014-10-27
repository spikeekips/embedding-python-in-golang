package main

/*
#cgo pkg-config: python-2.7
#include "Python.h"
*/
import "C"
import "unsafe"

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/spikeekips/embedding-python-in-golang/threadpool_python"
)

var log = logging.MustGetLogger("pyingo")

func main() {
	logging.SetLevel(logging.INFO, "pyingo")
	logging.SetLevel(logging.DEBUG, "pyingo")

	log.Info("main")

	python := pyingo.NewPython(1)
	_result, _err := python.Run("json_dumps:run")

	_unsafed := (*C.PyObject)(unsafe.Pointer(_result.Ptr))

	fmt.Printf("     C.PyString_AsString: %s, %T\n", C.GoString(C.PyString_AsString(_unsafed)), _unsafed)
	fmt.Printf("pyingo.PyString_AsString: %s, %T\n", pyingo.PyString_AsString(_result), _unsafed)

	fmt.Printf("result: %v, err: %v\n", _result, _err)
}
