package main

import (
	"github.com/spikeekips/embedding-python-in-golang/threadpool"
	"runtime"
	"sync"
	"unsafe"
)

var lock sync.Mutex

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
					_f := Fibonacci(_v)
					_values_collected = append(_values_collected, _f)

					_tid := get_thread_id()
					_thread_ids = append(_thread_ids, _tid)

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

func main() {
	_number_to_run := 100 * 13
	//_number_to_run = 2

	_run(run_in_goroutine_without_pthread_join_with_channel, _number_to_run, true)
}
