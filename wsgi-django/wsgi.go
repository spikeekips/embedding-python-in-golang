package wsgi

/*
// most of the wsgi code was derived from
// [https://github.com/GrahamDumpleton/mod_wsgi].

#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdio.h>


////////////////////////////////////////////////////////////////////////////////
// create StringIO
static PyObject *get_stringio (char *s)
{
	PyObject *_module, *_stringio, *_args;

    _module = PyImport_Import(PyString_FromString("StringIO"));
    _stringio = PyObject_GetAttr(_module, PyString_FromString("StringIO"));

	_args = Py_BuildValue("(z)", s);

    return PyObject_CallObject(_stringio, _args);
}

static PyObject *read_response_body (PyObject *iterator)
{
	PyObject *_iterator, *output;
	output = PyString_FromString("");

	_iterator = PyObject_GetIter(iterator);
	if (_iterator == NULL) {
		return output;
	}

	PyObject *item = NULL;
    long length = 0;

	while ((item = PyIter_Next(_iterator))) {
	    if (!PyString_Check(item)) {
	        PyErr_Format(PyExc_TypeError, "sequence of byte "
	                     "string values expected, value of "
	                     "type %.200s found",
	                     item->ob_type->tp_name);
	        Py_DECREF(item);
	        break;
	    }

		length = PyString_Size(item);

		if (length < 1) {
		    Py_DECREF(item);
		    break;
		}

		PyString_Concat(&output, item);
	    Py_DECREF(item);
	}

	return output;
}


////////////////////////////////////////////////////////////////////////////////


////////////////////////////////////////////////////////////////////////////////
// Adapter_Type
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

static PyTypeObject Adapter_Type;
static int HTTP_INTERNAL_SERVER_ERROR = 500;

typedef struct {
        PyObject_HEAD
        PyObject *log;
        int status;
        const char *status_line;
        PyObject *headers;
        PyObject *sequence;
		PyObject *wsgi_error;
} AdapterObject;

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
    "_wsgi.Adapter",     //tp_name
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

static AdapterObject *newAdapterObject()
{
    AdapterObject *self;

    self = PyObject_New(AdapterObject, &Adapter_Type);
    if (self == NULL)
        return NULL;

    self->status = HTTP_INTERNAL_SERVER_ERROR;
    self->status_line = NULL;
    self->headers = NULL;
	self->wsgi_error = get_stringio(NULL);

    return self;
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Response
typedef struct {
        int status;
        PyObject *headers;
		PyObject *body;
} Response;

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// run_wsgi_application
static Response run_wsgi_application(char *body, PyObject* environ) {
	Response response;
	response.status = 0;
	response.headers = NULL;
	response.body = Py_None;

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

    AdapterObject *_adapter = NULL;
    PyObject *_args = NULL;
    PyObject *_start = NULL;

    _adapter = newAdapterObject();

	// wsgi environ
    PyDict_SetItemString(environ, "wsgi.input", get_stringio(body));
    PyDict_SetItemString(environ, "wsgi.error", _adapter->wsgi_error);
    PyDict_SetItemString(environ, "wsgi.version", Py_BuildValue("(ii)", 1, 0));

    _start = PyObject_GetAttrString((PyObject *)_adapter, "start_response");

    _args = Py_BuildValue("(OO)", environ, _start);

    PyObject* _body = PyObject_CallObject(_attr, _args);

	// get the response status and header
	PyObject* _status = PyInt_FromLong((long)_adapter->status);
	PyObject* _headers = _adapter->headers;

	response.status = _adapter->status;
	response.headers = _headers;
	response.body = read_response_body(_body);

    // Clean up
    Py_DECREF(_module);
    Py_DECREF(_module_name);
    Py_DECREF(_attr);
    Py_DECREF(_start);
    Py_DECREF(_adapter);
    Py_DECREF(_body);

    return response;
}

////////////////////////////////////////////////////////////////////////////////


////////////////////////////////////////////////////////////////////////////////
static void initialize_wsgi () {
    PyType_Ready(&Adapter_Type);
}

////////////////////////////////////////////////////////////////////////////////

*/
import "C"

import (
	//"fmt"
	"github.com/op/go-logging"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"unsafe"
)

var log_wsgi = logging.MustGetLogger("wsgi")
var lock sync.Mutex

func init() {
	log_wsgi.Debug("> initialize wsgi.")
	C.initialize_wsgi()
}

func upgrade_header_to_wsgi(r *http.Request) *C.PyObject {
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

		// convert header name
		_k := strings.ToUpper(strings.Replace(k, "-", "_", -1))
		C.PyDict_SetItem(
			_environ,
			PyString_FromString("HTTP_"+_k),
			_values_tuple,
		)
	}

	C.PyDict_SetItem(
		_environ,
		PyString_FromString("X_FROM"),
		PyString_FromString("go"),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("REQUEST_METHOD"),
		PyString_FromString(r.Method),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("SCRIPT_NAME"),
		PyString_FromString(r.URL.Path),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("PATH_INFO"),
		PyString_FromString(r.URL.Path),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("QUERY_STRING"),
		PyString_FromString(r.URL.RawQuery),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("CONTENT_TYPE"),
		PyString_FromString(""),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("CONTENT_LENGTH"),
		PyString_FromString("0"),
	)

	_host, _port, _ := net.SplitHostPort(r.Host)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("SERVER_NAME"),
		PyString_FromString(_host),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("SERVER_PORT"),
		PyString_FromString(_port),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("SERVER_PROTOCOL"),
		PyString_FromString(r.Proto),
	)

	C.PyDict_SetItem(
		_environ,
		PyString_FromString("wsgi.url_scheme"),
		PyString_FromString(strings.ToLower(strings.Split(r.Proto, "/")[0])),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("wsgi.multithread"),
		C.PyBool_FromLong(1),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("wsgi.multiprocess"),
		C.PyBool_FromLong(1),
	)
	C.PyDict_SetItem(
		_environ,
		PyString_FromString("wsgi.run_once"),
		C.PyBool_FromLong(0),
	)

	return _environ
}

func wsgi_callback(r *http.Request, w http.ResponseWriter) {
	// TODO: add special headers for WSGI.
	_environ := upgrade_header_to_wsgi(r)

	_body_request, _ := ioutil.ReadAll(r.Body)
	_body_c := C.CString(string(_body_request))
	defer C.free(unsafe.Pointer(_body_c))

	_response := C.run_wsgi_application(_body_c, _environ)
	if _response.headers == nil || int(_response.status) == 0 {
		w.WriteHeader(500)
		log_wsgi.Error("failed to run python wsgi module.")
		return
	}

	// parse header
	for i := 0; i < int(C.PyList_Size(_response.headers)); i++ {
		_h := C.PyList_GetItem(_response.headers, C.Py_ssize_t(i))
		_k := C.PyTuple_GetItem(_h, C.Py_ssize_t(0))
		_v := C.PyTuple_GetItem(_h, C.Py_ssize_t(1))
		w.Header().Set(
			PyString_AsString(_k),
			PyString_AsString(_v),
		)
	}

	_body_response := PyString_AsString(_response.body)

	// write body
	w.WriteHeader(int(_response.status))
	w.Write([]byte(_body_response))
}

func WSGIHandler(w http.ResponseWriter, r *http.Request) {
	log_wsgi.Debug("> request: wsgi_handler")

	lock.Lock()
	defer lock.Unlock()

	_ch := make(chan bool)
	CreateThread(
		func() {
			wsgi_callback(r, w)
			_ch <- true
		},
	)
	<-_ch
	close(_ch)
	_ch = nil

	log_wsgi.Debug("< request: wsgi_handler")
}
