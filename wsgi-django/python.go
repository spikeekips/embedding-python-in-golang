package wsgi

/*
#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdio.h>

*/
import "C"
import "unsafe"

func PyString_FromString(s string) *C.PyObject {
	_s := C.CString(s)
	defer C.free(unsafe.Pointer(_s))

	return C.PyString_FromString(_s)
}

func PyString_AsString(s *C.PyObject) string {
	return C.GoString(C.PyString_AsString(s))
}
