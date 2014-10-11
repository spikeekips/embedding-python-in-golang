package main

import (
	"github.com/spikeekips/embedding-python-in-golang/threadpool"
	"runtime"
	"unsafe"
)

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
				_f := Fibonacci(_v)
				_values_collected = append(_values_collected, _f)

				_tid := get_thread_id()
				_thread_ids = append(_thread_ids, _tid)

				log.Debug("  [%5d] tid=%p, Fibonacci=%-7d, args=(%T, %T)", _v, _tid, _f, _v)
			}),
			int(i),
		)
		tp.Join(tid)
		log.Debug("< [%5d]", i)
	}

	return _values_collected, _thread_ids
}

func main() {
	_number_to_run := 100 * 13
	//_number_to_run = 2

	_run(run_in_block_with_pthread_join, _number_to_run, true)
}
