package main

/*
#include <pthread.h>
*/
import "C"

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/spikeekips/embedding-python-in-golang/threadpool"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"time"
	"unsafe"
)

// return Fibonacci number
func _Fibonacci(n int) int {
	if n == 0 {
		n = 1
	} else {
		n += 1
	}
	n = 33 - (33 % n)

	return Fibonacci(n)
}

func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}

	if n == 1 {
		return 1
	}

	return (Fibonacci(n-1) + Fibonacci(n-2))
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

/*
	run the native goroutine function.
*/
func run_go_in_goroutine(number_to_run int) ([]int, []string) {
	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	var _wg sync.WaitGroup
	_wg.Add(number_to_run)
	for i := 0; i < number_to_run; i++ {
		go func(j int) {
			_f := _Fibonacci(j)
			_values_collected = append(_values_collected, _f)
			_tid := C.pthread_self()
			_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

			log.Debug("result [%5d], tid=%p, Fibonacci=%-7d", j, _tid, _f)
			_wg.Done()
		}(int(i))
	}
	_wg.Wait()

	return _values_collected, _thread_ids
}

var lock sync.Mutex

/*
	run threads inside goroutine and each thread will be blocked by
	`C.pthread_join` and wait the end of execution with `channel`
*/
func run_in_goroutine_with_pthread_join_and_channel(number_to_run int) ([]int, []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	_channel_collect_seq := make(chan int)

	/*
		if you increase the `number_to_run` too high, you may face the
		segfault, it's totally up to your system.
	*/
	tp := threadpool.NewThreadPool(number_to_run)

	for i := 0; i < number_to_run; i++ {
		go func(j int, ch chan int) {
			lock.Lock()
			defer lock.Unlock()

			tid := tp.Start(
				threadpool.ThreadCallback(func(args unsafe.Pointer) {
					_args := *(**([2]interface{}))(unsafe.Pointer(&args))

					_v := _args[0].(int)
					_ch := _args[1].(chan int)

					// run go function
					_f := _Fibonacci(_v)
					_values_collected = append(_values_collected, _f)
					_tid := C.pthread_self()
					_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

					log.Debug("result: [%5d] tid=%p, Fibonacci=%-7d, args=(%T, %T)", _v, _tid, _f, _v, _ch)
					_ch <- _v
				}),
				j,
				ch,
			)
			tp.Join(tid)
		}(int(i), _channel_collect_seq)
	}

	_done := make(chan bool, 1)
	_got := []int{}
	for _v := range _channel_collect_seq {
		_got = append(_got, _v)
		if len(_got) == number_to_run {
			_done <- true
			break
		}
	}
	<-_done
	close(_channel_collect_seq)
	close(_done)

	return _values_collected, _thread_ids
}

/*
	run threads inside goroutine and each thread will be blocked by
	`C.pthread_join` and wait the end of execution with `sync.WaitGroup`
*/
func run_in_goroutine_with_pthread_join_and_WaitGroup(number_to_run int) ([]int, []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	var wg sync.WaitGroup
	wg.Add(number_to_run)

	tp := threadpool.NewThreadPool(number_to_run)

	for i := 0; i < number_to_run; i++ {
		go func(j int) {
			lock.Lock()
			defer lock.Unlock()

			tid := tp.Start(
				threadpool.ThreadCallback(func(args unsafe.Pointer) {
					_args := *(**([2]interface{}))(unsafe.Pointer(&args))

					_v := _args[0].(int)

					// run go function
					_f := _Fibonacci(_v)
					_values_collected = append(_values_collected, _f)

					_tid := C.pthread_self()
					_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

					log.Debug("result: [%5d] tid=%p, Fibonacci=%-7d, args=(%T, %T)", _v, _tid, _f, _v, wg)
					wg.Done()
				}),
				j,
			)
			tp.Join(tid)
		}(int(i))
	}

	wg.Wait()

	return _values_collected, _thread_ids
}

/*
	run threads inside goroutine and each thread will be blocked by `channel`
	and wait the end of execution with `channel`
*/
func run_in_goroutine_without_pthread_join_with_channel(number_to_run int) ([]int, []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_values_collected := make([]int, 0)
	_channel_collect_seq := make(chan int)
	_thread_ids := make([]string, 0)

	tp := threadpool.NewThreadPool(number_to_run)

	for i := 0; i < number_to_run; i++ {
		go func(j int, ch chan int) {
			lock.Lock()
			defer lock.Unlock()

			_chan_thread := make(chan bool)
			tp.Start(
				threadpool.ThreadCallback(func(args unsafe.Pointer) {
					_args := *(**([2]interface{}))(unsafe.Pointer(&args))

					_v := _args[0].(int)
					_ch := _args[1].(chan bool)

					// run go function
					_f := _Fibonacci(_v)
					_values_collected = append(_values_collected, _f)

					_tid := C.pthread_self()
					_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

					log.Debug("result: [%5d] tid=%p, Fibonacci=%-7d, args=(%T, %T)", _v, _tid, _f, _v, _ch)
					_ch <- true
				}),
				j,
				_chan_thread,
			)
			<-_chan_thread
			close(_chan_thread)
			ch <- j
		}(int(i), _channel_collect_seq)
	}

	_done := make(chan bool, 1)
	_got := []int{}
	for _v := range _channel_collect_seq {
		_got = append(_got, _v)
		if len(_got) == number_to_run {
			_done <- true
			break
		}
	}
	<-_done
	close(_channel_collect_seq)
	close(_done)

	return _values_collected, _thread_ids
}

/*
	run threads inside goroutine and each thread will be blocked by
	`sync.WaitGroup` and wait the end of execution with `WaitGroup`
*/
func run_in_goroutine_without_pthread_join_with_WaitGroup(number_to_run int) ([]int, []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	var wg sync.WaitGroup
	wg.Add(number_to_run)

	tp := threadpool.NewThreadPool(number_to_run)

	for i := 0; i < number_to_run; i++ {
		go func(j int) {
			lock.Lock()
			defer lock.Unlock()

			var wg_thread sync.WaitGroup
			wg_thread.Add(1)
			tp.Start(
				threadpool.ThreadCallback(func(args unsafe.Pointer) {
					_args := *(**([2]interface{}))(unsafe.Pointer(&args))

					_v := _args[0].(int)

					// run go function
					_f := _Fibonacci(_v)
					_values_collected = append(_values_collected, _f)

					_tid := C.pthread_self()
					_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

					log.Debug("result: [%5d] tid=%p, Fibonacci=%-7d, args=(%T)", _v, _tid, _f, _v)
					wg_thread.Done()
				}),
				j,
			)
			wg_thread.Wait()
			wg.Done()
		}(int(i))
	}

	wg.Wait()

	return _values_collected, _thread_ids
}

/*
	run the native go function.
*/
func run_go_in_block(number_to_run int) ([]int, []string) {
	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	for i := 0; i < number_to_run; i++ {
		_f := _Fibonacci(i)
		_values_collected = append(_values_collected, _f)
		_tid := C.pthread_self()
		_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

		log.Debug("result [%5d], tid=%p, Fibonacci=%-7d", i, _tid, _f)
	}

	return _values_collected, _thread_ids
}

/*
	run threads inside goroutine and each thread will be blocked by
	`C.pthread_join`.
*/
func run_in_block_with_pthread_join(number_to_run int) ([]int, []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	tp := threadpool.NewThreadPool(number_to_run)
	for i := 0; i < number_to_run; i++ {
		log.Debug("> [%5d]", i)

		tid := tp.Start(
			threadpool.ThreadCallback(func(args unsafe.Pointer) {
				_args := *(**([1]interface{}))(unsafe.Pointer(&args))
				_v, _ := _args[0].(int)

				// run go function
				_f := _Fibonacci(_v)
				_values_collected = append(_values_collected, _f)

				_tid := C.pthread_self()
				_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

				log.Debug("  [%5d] tid=%p, Fibonacci=%-7d, args=(%T, %T)", _v, _tid, _f, _v)
			}),
			int(i),
		)
		tp.Join(tid)
		log.Debug("< [%5d]", i)
	}

	return _values_collected, _thread_ids
}

/*
	run threads in golang and each thread will be blocked by `channel`.
*/
func run_in_block_with_channel(number_to_run int) ([]int, []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	tp := threadpool.NewThreadPool(number_to_run + 1)
	for i := 0; i < number_to_run; i++ {
		j := int(i)

		log.Debug("> [%5d]", j)
		_chan := make(chan bool)
		tp.Start(
			threadpool.ThreadCallback(func(args unsafe.Pointer) {
				_args := *(**([2]interface{}))(unsafe.Pointer(&args))
				_v, _ := _args[0].(int)
				_ch, _ := _args[1].(chan bool)

				// run go function
				_f := _Fibonacci(_v)
				_values_collected = append(_values_collected, _f)

				_tid := C.pthread_self()
				_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

				log.Debug("  [%5d], tid=%p, Fibonacci=%-7d, args=(%T, %T)", _v, _tid, _f, _v, _ch)
				_ch <- true
			}),
			j,
			_chan,
		)
		<-_chan
		close(_chan)
		log.Debug("< [%5d]", j)
	}

	return _values_collected, _thread_ids
}

/*
	run threads in golang and each thread will be blocked by `sync.WaitGroup`.
*/
func run_in_block_with_WaitGroup(number_to_run int) ([]int, []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	tp := threadpool.NewThreadPool(number_to_run + 1)
	for i := 0; i < number_to_run; i++ {
		j := int(i)

		log.Debug("> [%5d]", j)
		var wg sync.WaitGroup
		wg.Add(1)
		tp.Start(
			threadpool.ThreadCallback(func(args unsafe.Pointer) {
				_args := *(**([2]interface{}))(unsafe.Pointer(&args))
				_v, _ := _args[0].(int)

				// run go function
				_f := _Fibonacci(_v)
				_values_collected = append(_values_collected, _f)
				_tid := C.pthread_self()
				_thread_ids = append(_thread_ids, fmt.Sprintf("%p", _tid))

				log.Debug("  [%5d], tid=%p, Fibonacci=%-7d, args=(%T)", _v, _tid, _f, _v)
				wg.Done()
			}),
			j,
		)
		wg.Wait()
		log.Debug("< [%5d]", j)
	}

	return _values_collected, _thread_ids
}

func _run(func_name string, number_to_run int, compare_length bool) bool {
	functions := map[string]func(int) ([]int, []string){
		"run_go_in_goroutine":                                  run_go_in_goroutine,
		"run_in_goroutine_with_pthread_join_and_channel":       run_in_goroutine_with_pthread_join_and_channel,
		"run_in_goroutine_with_pthread_join_and_WaitGroup":     run_in_goroutine_with_pthread_join_and_WaitGroup,
		"run_in_goroutine_without_pthread_join_with_channel":   run_in_goroutine_without_pthread_join_with_channel,
		"run_in_goroutine_without_pthread_join_with_WaitGroup": run_in_goroutine_without_pthread_join_with_WaitGroup,
		"run_go_in_block":                                      run_go_in_block,
		"run_in_block_with_pthread_join":                       run_in_block_with_pthread_join,
		"run_in_block_with_channel":                            run_in_block_with_channel,
		"run_in_block_with_WaitGroup":                          run_in_block_with_WaitGroup,
	}

	_fib_numbers := make([]int, 0)
	for i := 0; i < number_to_run; i++ {
		_fib_numbers = append(_fib_numbers, _Fibonacci(i))
	}

	_f := reflect.ValueOf(functions[func_name])

	_args := make([]reflect.Value, 1)
	_args[0] = reflect.ValueOf(number_to_run)

	log.Debug("> %53s", func_name)

	_now := time.Now()
	_r := _f.Call(_args)
	_elapsed := time.Now().Sub(_now).Seconds()

	_fib_numbers_result := _r[0].Interface().([]int)
	_thread_ids_result := _r[1].Interface().([]string)

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

var log = logging.MustGetLogger("threadpool")

func main() {
	logging.SetLevel(logging.DEBUG, "threadpool")
	logging.SetLevel(logging.INFO, "threadpool")

	_number_to_run := 100 * 13
	//_number_to_run = 2

	log.Info("# run_go")
	_run("run_go_in_goroutine", _number_to_run, true)

	log.Info(" ")
	log.Info("# run_in_goroutine_with_pthread_join")
	_run("run_in_goroutine_with_pthread_join_and_channel", _number_to_run, true)
	_run("run_in_goroutine_with_pthread_join_and_WaitGroup", _number_to_run, true)

	log.Info(" ")
	log.Info("# run_in_goroutine_without_pthread_join")
	_run("run_in_goroutine_without_pthread_join_with_channel", _number_to_run, true)
	_run("run_in_goroutine_without_pthread_join_with_WaitGroup", _number_to_run, true)

	log.Info(" ")
	log.Info("# run_in_block")
	_run("run_go_in_block", _number_to_run, false)
	_run("run_in_block_with_pthread_join", _number_to_run, false)
	_run("run_in_block_with_channel", _number_to_run, false)
	_run("run_in_block_with_WaitGroup", _number_to_run, false)
}
