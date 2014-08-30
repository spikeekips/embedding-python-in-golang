package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/sbinet/go-python"
	"github.com/spikeekips/embedding-python-in-golang"
	"sync"
)

func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
}

func embed_function(args []string, kw map[string]string) {
	lock.Lock()
	defer lock.Unlock()

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

	log.Info("python embedding test in golang.")

	_ch := make(chan bool)

	// create argument
	_args := []string{"oldboy", "나는 오대수"}

	_kw := map[string]string{
		"who are you?":                   "누구냐, 넌?",
		"I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
	}

	// in goroutine
	log.Info("run inside goroutine")
	go func() {
		embed_function(_args, _kw)
		_ch <- true
	}()
	<-_ch

	fmt.Println("")

	// outside goroutine
	log.Info("run outside goroutine")
	embed_function(_args, _kw)
}
