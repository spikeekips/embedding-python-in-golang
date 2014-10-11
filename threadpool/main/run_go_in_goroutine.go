package main

import (
	"sync"
)

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
			_f := Fibonacci(j)
			_values_collected = append(_values_collected, _f)
			_tid := get_thread_id()
			_thread_ids = append(_thread_ids, _tid)

			log.Debug("result [%5d], tid=%p, Fibonacci=%-7d", j, _tid, _f)
			_wg.Done()
		}(int(i))
	}
	_wg.Wait()

	return _values_collected, _thread_ids
}

func main() {
	_number_to_run := 100 * 13
	//_number_to_run = 2

	_run(run_go_in_goroutine, _number_to_run, true)
}
