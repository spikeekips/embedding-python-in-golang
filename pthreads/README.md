## Using Pthread, free from multi-threads problem, based on go-pthread and go-python

The first test, `sync_mutext.go` will be only available without multi-threads
environment. If you run the thread code like below, the method of
`sync_mutext.go` will be failed with segfault.

To run the threaded python code, we need to create thread in golang using
`pthread`. In this test, I will use `go-pthread`.

```
$ go get github.com/sbinet/go-python
$ go get github.com/spikeekips/go-pthreads
$ go get github.com/op/go-logging

$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make pthreads.go
2014/09/02 03:14:47 > Initilize Python.
2014/09/02 03:14:47 < Initilized Python.
2014/09/02 03:14:47 python embedding test in golang, using `pthreads`.
2014/09/02 03:14:47 run inside goroutine
2014/09/02 03:14:47 > create_thread
2014/09/02 03:14:47 got json string:
[
  4.7053070751488e+13,
  [
    "oldboy",
    "나는 오대수"
  ],
  {
    "I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
    "who are you?": "누구냐, 넌?"
  }
]
2014/09/02 03:14:47 < create_thread: done
2014/09/02 03:14:47 run outside goroutine
2014/09/02 03:14:47 > create_thread
2014/09/02 03:14:47 got json string:
[
  4.7053075789568e+13,
  [
    "oldboy",
    "나는 오대수"
  ],
  {
    "I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
    "who are you?": "누구냐, 넌?"
  }
]
2014/09/02 03:14:47 < create_thread: done
```


