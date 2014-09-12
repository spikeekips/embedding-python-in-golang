package wsgi

/*

// most of the wsgi code was derived from
// [https://github.com/GrahamDumpleton/mod_wsgi].

#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <pthread.h>
#include <unistd.h>
#include <stdio.h>

static PyObject *wsgi_convert_string_to_bytes(PyObject *value)
{
    PyObject *result = NULL;

    if (!PyBytes_Check(value)) {
        PyErr_Format(PyExc_TypeError, "expected byte string object, "
                     "value of type %.200s found", value->ob_type->tp_name);
        return NULL;
    }

    Py_INCREF(value);
    result = value;

    return result;
}

static PyObject *wsgi_convert_headers_to_bytes(PyObject *headers)
{
    PyObject *result = NULL;

    int i;
    long size;

    if (!PyList_Check(headers)) {
        PyErr_Format(PyExc_TypeError, "expected list object for headers, "
                     "value of type %.200s found", headers->ob_type->tp_name);
        return Py_None;
    }

    size = PyList_Size(headers);
    result = PyList_New(size);

    for (i = 0; i < size; i++) {
        PyObject *header = NULL;

        PyObject *header_name = NULL;
        PyObject *header_value = NULL;

        PyObject *header_name_as_bytes = NULL;
        PyObject *header_value_as_bytes = NULL;

        PyObject *result_tuple = NULL;

        header = PyList_GetItem(headers, i);

        if (!PyTuple_Check(header)) {
            PyErr_Format(PyExc_TypeError, "list of tuple values "
                         "expected for headers, value of type %.200s found",
                         header->ob_type->tp_name);
            Py_DECREF(result);
            return Py_None;
        }

        if (PyTuple_Size(header) != 2) {
            PyErr_Format(PyExc_ValueError, "tuple of length 2 "
                         "expected for header, length is %d",
                         (int)PyTuple_Size(header));
            Py_DECREF(result);
            return Py_None;
        }

        result_tuple = PyTuple_New(2);
        PyList_SET_ITEM(result, i, result_tuple);

        header_name = PyTuple_GetItem(header, 0);
        header_value = PyTuple_GetItem(header, 1);

        header_name_as_bytes = wsgi_convert_string_to_bytes(header_name);
        if (!header_name_as_bytes)
            goto failure;

        PyTuple_SET_ITEM(result_tuple, 0, header_name_as_bytes);

        header_value_as_bytes = wsgi_convert_string_to_bytes(header_value);
        if (!header_value_as_bytes)
            goto failure;

        PyTuple_SET_ITEM(result_tuple, 1, header_value_as_bytes);
    }

    return result;

failure:
    Py_DECREF(result);
    return Py_None;
}

typedef struct {
        int status;
        PyObject *headers;
		char *body;
} Response;


typedef struct {
        PyObject_HEAD
        PyObject *log;
        int status;
        const char *status_line;
        PyObject *headers;
        PyObject *sequence;
} AdapterObject;


static PyTypeObject Adapter_Type;

static int HTTP_INTERNAL_SERVER_ERROR = 500;

static AdapterObject *newAdapterObject()
{
    AdapterObject *self;

    self = PyObject_New(AdapterObject, &Adapter_Type);
    if (self == NULL)
        return NULL;

    self->status = HTTP_INTERNAL_SERVER_ERROR;
    self->status_line = NULL;
    self->headers = NULL;

    return self;
}


static void Adapter_dealloc(AdapterObject *self)
{
    PyObject_Del(self);
}

static PyObject *Adapter_start_response(AdapterObject *self, PyObject *args)
{
    PyObject *status_line = NULL;
    PyObject *headers = NULL;
    PyObject *exc_info = Py_None;

    PyObject *status_line_as_bytes = NULL;
    PyObject *headers_as_bytes = NULL;

    if (!PyArg_ParseTuple(args, "OO!|O:start_response",
        &status_line, &PyList_Type, &headers, &exc_info)) {
        return Py_None;
    }

    if (exc_info != Py_None && !PyTuple_Check(exc_info)) {
        PyErr_SetString(PyExc_RuntimeError, "exception info must be a tuple");
        return Py_None;
    }

    if (exc_info != Py_None) {
        if (self->status_line && !self->headers) {
            PyObject *type = NULL;
            PyObject *value = NULL;
            PyObject *traceback = NULL;

            if (!PyArg_ParseTuple(exc_info, "OOO", &type,
                                  &value, &traceback)) {
                return Py_None;
            }

            Py_INCREF(type);
            Py_INCREF(value);
            Py_INCREF(traceback);

            PyErr_Restore(type, value, traceback);

            return Py_None;
        }
    }
    else if (self->status_line && !self->headers) {
        PyErr_SetString(PyExc_RuntimeError, "headers have already been sent");
        return Py_None;
    }

    headers_as_bytes = wsgi_convert_headers_to_bytes(headers);

    if (!headers_as_bytes)
       goto finally;

    self->status_line = PyString_AsString(status_line);
    self->status = (int)strtol(self->status_line, NULL, 10);

    Py_XDECREF(self->headers);
    self->headers = headers_as_bytes;
    Py_INCREF(headers_as_bytes);

finally:
    Py_XDECREF(headers_as_bytes);

    return Py_None;
}

static PyMethodDef Adapter_methods[] = {
    { "start_response", (PyCFunction)Adapter_start_response, METH_VARARGS, 0 },
    { NULL, NULL}
};

static PyTypeObject Adapter_Type = {
    PyVarObject_HEAD_INIT(NULL, 0)
    "wsgi.Adapter",     //tp_name
    sizeof(AdapterObject),  //tp_basicsize
    0,                      //tp_itemsize
    // methods
    (destructor)Adapter_dealloc, //tp_dealloc
    0,                      //tp_print
    0,                      //tp_getattr
    0,                      //tp_setattr
    0,                      //tp_compare
    0,                      //tp_repr
    0,                      //tp_as_number
    0,                      //tp_as_sequence
    0,                      //tp_as_mapping
    0,                      //tp_hash
    0,                      //tp_call
    0,                      //tp_str
    0,                      //tp_getattro
    0,                      //tp_setattro
    0,                      //tp_as_buffer
    Py_TPFLAGS_DEFAULT,     //tp_flags
    0,                      //tp_doc
    0,                      //tp_traverse
    0,                      //tp_clear
    0,                      //tp_richcompare
    0,                      //tp_weaklistoffset
    0,                      //tp_iter
    0,                      //tp_iternext
    Adapter_methods,        //tp_methods
    0,                      //tp_members
    0,                      //tp_getset
    0,                      //tp_base
    0,                      //tp_dict
    0,                      //tp_descr_get
    0,                      //tp_descr_set
    0,                      //tp_dictoffset
    0,                      //tp_init
    0,                      //tp_alloc
    0,                      //tp_new
    0,                      //tp_free
    0,                      //tp_is_gc
};

extern void createThreadCallback();

static PyGILState_STATE start_thread() {
    return PyGILState_Ensure();
}

static void end_thread(PyGILState_STATE gstate) {
    PyGILState_Release(gstate);
}

static Response run_wsgi_application(PyObject* environ) {
	Response response;

    PyObject *_module_name, *_module;

    // Now execute some python code (call python functions)
    _module_name = PyString_FromString("wsgi");
    _module = PyImport_Import(_module_name);
    if (_module == NULL) {
        fprintf(stderr, "failed to import wsgi python module\n");
        return response;
    }

    // Call a method of the class with no parameters
    PyObject *_attr, *_result;
    _attr = PyObject_GetAttr(_module, PyString_FromString("application"));

    AdapterObject *adapter = NULL;
    PyObject *args = NULL;
    PyObject *start = NULL;

    adapter = newAdapterObject();

    start = PyObject_GetAttrString((PyObject *)adapter, "start_response");

    args = Py_BuildValue("(OO)", environ, start);

    PyObject* _body = PyObject_CallObject(_attr, args);

	// get the response status and header

	PyObject* _status = PyInt_FromLong((long)adapter->status);
	PyObject* _headers = adapter->headers;

    // Clean up
    Py_DECREF(_module);
    Py_DECREF(_module_name);
    Py_DECREF(_attr);

	response.status = adapter->status;
	response.headers = _headers;
	response.body = PyString_AsString(_body);

    return response;
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
        PyEval_ReleaseThread(PyGILState_GetThisThreadState());
    }
    PyType_Ready(&Adapter_Type);
}
*/
import "C"
import "unsafe"

import (
	"fmt"
	"github.com/op/go-logging"
	"net/http"
	"runtime"
	"sync"
)

func PyString_FromString(s string) *C.PyObject {
	_s := C.CString(s)
	defer C.free(unsafe.Pointer(_s))

	return C.PyString_FromString(_s)
}

func PyString_AsString(s *C.PyObject) string {
	return C.GoString(C.PyString_AsString(s))
}

type Callbacks struct {
	sync.RWMutex
	m map[string]ThreadCallback
	r map[string]*http.Request
	w map[string]http.ResponseWriter
}

func (cb *Callbacks) Add(pid *C.pthread_t, tc ThreadCallback, r *http.Request, w http.ResponseWriter) {
	cb.Lock()
	defer cb.Unlock()

	_key := fmt.Sprintf("%p", pid)
	cb.m[_key] = tc
	cb.r[_key] = r
	cb.w[_key] = w
}

func (cb *Callbacks) Get(pid *C.pthread_t) (
	ThreadCallback,
	*http.Request,
	http.ResponseWriter,
	bool,
) {
	cb.RLock()
	defer cb.RUnlock()

	var _ok bool
	var _tc ThreadCallback
	var _r *http.Request
	var _w http.ResponseWriter

	_tc, _ok = cb.m[fmt.Sprintf("%p", pid)]
	_r, _ok = cb.r[fmt.Sprintf("%p", pid)]
	_w, _ok = cb.w[fmt.Sprintf("%p", pid)]

	return _tc, _r, _w, _ok
}

func (cb *Callbacks) Delete(pid *C.pthread_t) {
	cb.Lock()
	defer cb.Unlock()
	delete(cb.m, fmt.Sprintf("%p", pid))
	delete(cb.r, fmt.Sprintf("%p", pid))
	delete(cb.w, fmt.Sprintf("%p", pid))
}

var callbacks Callbacks

type ThreadCallback func()

var lock sync.Mutex

func GenerateEnviron(r *http.Request) *C.PyObject {
	_environ := C.PyDict_New()

	for k, _items := range r.Header {
		_values_tuple := C.PyTuple_New(C.Py_ssize_t(len(_items)))
		for i, _item := range _items {
			C.PyTuple_SetItem(
				_values_tuple,
				C.Py_ssize_t(i),
				PyString_FromString(_item),
			)
		}
		C.PyDict_SetItem(
			_environ,
			PyString_FromString(k),
			_values_tuple,
		)
	}
	//_environ = upgrade_to_wsgi(r, _environ)

	return _environ
}

//export createThreadCallback
func createThreadCallback(pid *C.pthread_t) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_gstate := C.start_thread()

	_cb, _r, _w, _ok := callbacks.Get(pid)
	defer callbacks.Delete(pid)

	if !_ok {
		panic(fmt.Errorf("failed to found thread callback for `%v`", pid))
	}

	// TODO: add special headers for WSGI.
	_environ := GenerateEnviron(_r)

	_response := C.run_wsgi_application(_environ)

	// parse header
	for i := 0; i < int(C.PyList_Size(_response.headers)); i++ {
		_h := C.PyList_GetItem(_response.headers, C.Py_ssize_t(i))
		_k := C.PyTuple_GetItem(_h, C.Py_ssize_t(0))
		_v := C.PyTuple_GetItem(_h, C.Py_ssize_t(1))
		_w.Header().Set(
			PyString_AsString(_k),
			PyString_AsString(_v),
		)
	}

	_body := C.GoString(_response.body)

	if len(_body) < 1 {
		panic(fmt.Errorf("failed to import python wsgi module."))
	}
	C.end_thread(_gstate)

	// write body
	_w.WriteHeader(int(_response.status))
	_w.Write([]byte(_body))

	_cb()
}

func CreateThread(cb ThreadCallback, r *http.Request, w http.ResponseWriter) {
	_pid := new(C.pthread_t)
	callbacks.Add(_pid, cb, r, w)

	C.createThread(_pid)
}

var log = logging.MustGetLogger("epig.thread")

func init() {
	log.Debug("> Initilize Python.")
	callbacks = Callbacks{
		m: map[string]ThreadCallback{},
		r: map[string]*http.Request{},
		w: map[string]http.ResponseWriter{},
	}

	C.initialize_python()
}

func WSGIHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("> request: wsgi_handler")
	lock.Lock()
	defer lock.Unlock()

	_ch := make(chan bool)
	CreateThread(
		func() {
			_ch <- true
		},
		r, w,
	)
	<-_ch
	close(_ch)
	_ch = nil

	log.Debug("< request: wsgi_handler")
}
