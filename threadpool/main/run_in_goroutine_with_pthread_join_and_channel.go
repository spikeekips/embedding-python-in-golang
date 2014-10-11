package main

import (
	"github.com/spikeekips/embedding-python-in-golang/threadpool"
	"runtime"
	"sync"
	"unsafe"
)

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
					_f := Fibonacci(_v)
					_values_collected = append(_values_collected, _f)
					_tid := get_thread_id()
					_thread_ids = append(_thread_ids, _tid)

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

func main() {
	_number_to_run := 100 * 13
	//_number_to_run = 2

	_run(run_in_goroutine_with_pthread_join_and_channel, _number_to_run, true)
}
