package main

import ()

/*
	run the native go function.
*/
func run_go_in_block(number_to_run int) ([]int, []string) {
	_values_collected := make([]int, 0)
	_thread_ids := make([]string, 0)

	for i := 0; i < number_to_run; i++ {
		_f := Fibonacci(i)
		_values_collected = append(_values_collected, _f)
		_tid := get_thread_id()
		_thread_ids = append(_thread_ids, _tid)

		log.Debug("result [%5d], tid=%p, Fibonacci=%-7d", i, _tid, _f)
	}

	return _values_collected, _thread_ids
}

func main() {
	_number_to_run := 100 * 13
	//_number_to_run = 2

	_run(run_go_in_block, _number_to_run, true)
}
