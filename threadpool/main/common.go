package main

/*
#include <pthread.h>
*/
import "C"

import (
	"fmt"
	"github.com/op/go-logging"
	"reflect"
	"runtime"
	"sort"
	"time"
)

var log = logging.MustGetLogger("threadpool")

func init() {
	logging.SetLevel(logging.DEBUG, "threadpool")
	logging.SetLevel(logging.INFO, "threadpool")
}

func get_thread_id() string {
	return fmt.Sprintf("%p", C.pthread_self())
}

// return Fibonacci number
func Fibonacci(n int) int {
	if n == 0 {
		n = 1
	} else {
		n += 1
	}
	n = 33 - (33 % n)

	return _Fibonacci(n)
}

func _Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}

	if n == 1 {
		return 1
	}

	return (_Fibonacci(n-1) + _Fibonacci(n-2))
}

func make_array_unique(a []string) []string {
	b := []string{}
	for _, i := range a {
		if sort.SearchStrings(b, i) == len(b) {
			b = append(b, i)
		}
	}

	return b
}

func _run(function func(a int) ([]int, []string), number_to_run int, compare_length bool) bool {
	func_name := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()

	_fib_numbers := make([]int, 0)
	for i := 0; i < number_to_run; i++ {
		_fib_numbers = append(_fib_numbers, Fibonacci(i))
	}

	log.Debug("> %53s", func_name)

	_now := time.Now()
	_fib_numbers_result, _thread_ids_result := function(number_to_run)
	_elapsed := time.Now().Sub(_now).Seconds()

	// compare with numbers

	_thread_id_is_unique := ""
	if len(_thread_ids_result) != len(make_array_unique(_thread_ids_result)) {
		_thread_id_is_unique = "not unique"
	} else {
		_thread_id_is_unique = "unique"
	}
	log.Debug("< %53s: thread ids are %s", " ", _thread_id_is_unique)

	_error := false
	if compare_length {
		_error = len(_fib_numbers_result) != len(_fib_numbers)
	} else {
		for i := 0; i < len(_fib_numbers_result); i++ {
			if _fib_numbers_result[i] != _fib_numbers[i] {
				_error = true
				break
			}
		}
	}

	if !_error {
		log.Info("< %53s: elapsed %3.10f", func_name, _elapsed)
	} else {
		//fmt.Printf(">> %v\n", _fib_numbers_result)
		//fmt.Printf(">> %v\n", _fib_numbers)
		log.Error(
			"< %53s: elapsed %3.10f, but result is weired",
			func_name,
			_elapsed,
		)
		return false
	}

	return true
}
