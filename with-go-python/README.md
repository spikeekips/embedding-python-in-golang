## Embedding Python In Golang, Using `go-python`

I have some experiences in programming with python and I love it, I have been
always impressed with it's simple and elegant design, and most of all, it
brings the productivity to me. But as most of us know, python is slow. I have
been used to write application for long-running task, to handle massive data,
networking related jobs, etc. Python is good at various kind of area, but not
good choice for that fields, which is needed to the performance.

Recently I have managed the [serf](serfdom.io) and fun with it, especially it's
stability and performance. But one thing I missed, there was no python library.
As usual, so I wrote [serf-python](https://github.com/spikeekips/serf-python),
it provides all of the serf API.

So I tested how to use python inside golang, if possible, I will get the
combination of productivity and performance.

I tried embedding python in golang using
[go-python](https://github.com/sbinet/go-python).

### Big Picture

To embed python in golang, I used go-python and
[go-pthreads](https://github.com/liamzdenek/go-pthreads). go-python will
provide `C-API` of python and go-pthreads also make it possible to create
pthreads in golang.

The big obstacle is the multi-thread and the GIL of python. golang does not
provide the multi-thread and I found the `GIL` made the big problems. The
goroutine always makes problems with go-python and pthreads, caused by GIL.

This test and code is just /* the proof of concept */, based on this test I
will write the new code, based on the new way :)

### Installation

At first, I tried to install `go-python` in OSX, but I could not install
because,

* The builtin python does not support `pkg-config`
* and failed to compile, I got these kind of errors,

```
pkg-config --cflags python-2.7
pkg-config --libs python-2.7
CGO_LDFLAGS="-g" "-O2" "-L/System/Library/Frameworks/Python.framework/Versions/2.7/lib" "-ldl" "-lpython2.7" /Volumes/Userland/local/Cellar/go/1.3/libexec/pkg/tool/darwin_amd64/cgo -objdir $WORK/github.com/sbinet/go-python/_test/_obj_test/ -- -fno-strict-aliasing -fno-common -dynamic -g -Os -pipe -fno-common -fno-strict-aliasing -fwrapv -DENABLE_DTRACE -DMACOSX -DNDEBUG -Wall -Wstrict-prototypes -Wshorten-64-to-32 -DNDEBUG -g -fwrapv -Os -Wall -Wstrict-prototypes -DENABLE_DTRACE -I/System/Library/Frameworks/Python.framework/Versions/2.7/include/python2.7 -I $WORK/github.com/sbinet/go-python/_test/_obj_test/ cgoflags.go dict.go exceptions.go none.go numeric.go object.go otherobjects.go python.go sequence.go type.go utilities.go veryhigh.go
# github.com/sbinet/go-python
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xb code=0x1 addr=0x0 pc=0x15672]
```

Basically go-python runs only Linux, so I tried in my ubuntu machine.

```
$ sudo apt-get install golang golang-doc golang-go.tools golang-go.tools-dev golang-src gccgo gccgo-go
```

create base environment.

```
$ mkdir ~/gopython; cd gopython
$ mkdir -p {bin,pkg,src}
$ export GOPATH=$HOME/gopython
```

install the requiments,

```
$ go get github.com/sbinet/go-python
$ go get github.com/liamzdenek/go-pthreads
$ go get github.com/op/go-logging
$ go get github.com/spikeekips/embedding-python-in-golang
```

ready to code.

### test

#### lock the python-related code with `sync.Mutex`

```
$ cd src/github.com/spikeekips/embedding-python-in-golang/with-go-python
$ make sync_mutex.go
2014/08/29 20:51:21 python embedding test in golang.
2014/08/29 20:51:21 run inside goroutine
2014/08/29 20:51:21 got json string:
[
  [
    "oldboy",
    "나는 오대수"
  ],
  {
    "I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
    "who are you?": "누구냐, 넌?"
  }
]

2014/08/29 20:51:21 run outside goroutine
2014/08/29 20:51:21 got json string:
[
  [
    "oldboy",
    "나는 오대수"
  ],
  {
    "I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
    "who are you?": "누구냐, 넌?"
  }
]
```

`sync_mutex.go` just encoding golang data in python and decoding python data in
golang. It locks the `embed_function` using `sync.Mutex` and `embed_function`
`import`s `json_dump.py` using `go-python`, so simple.

This is `json_dump.py`,

```
import json

def run (*a) :
    return json.dumps(a, )


```

Everything was fine whether goroutine or not.



