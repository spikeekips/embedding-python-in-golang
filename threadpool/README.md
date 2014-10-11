## Apply *threadpool*

```
$ cd src/github.com/spikeekips/embedding-python-in-golang

$ make threadpool.run_go_in_block.go
$ make threadpool.run_go_in_goroutine.go
$ make threadpool.run_in_block_with_WaitGroup.go
$ make threadpool.run_in_block_with_channel.go
$ make threadpool.run_in_block_with_pthread_join.go
$ make threadpool.run_in_goroutine_with_pthread_join_and_WaitGroup.go
$ make threadpool.run_in_goroutine_with_pthread_join_and_channel.go
$ make threadpool.run_in_goroutine_without_pthread_join_with_WaitGroup.go
$ make threadpool.run_in_goroutine_without_pthread_join_with_channel.go
...
```



