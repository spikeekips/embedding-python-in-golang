## Embedding Python In Golang

I have some experiences in programming with python and I love it, I have been
always impressed with it's simple and elegant design, and most of all, it
brings the productivity to me. But as most of us know, python is slow. I have
been used to write application for long-running task, to handle massive data,
networking related jobs, etc. Python is good at various kind of area, but not
good choice for that fields, which is needed to the performance.

Recently I have tested and managed the [serf](serfdom.io) servers, which is
written in golang, I was very impressed it's stability and performance. But one
thing it was missed, there was no python library.  As usual, I wrote
[serf-python](https://github.com/spikeekips/serf-python), it can provide all of
the serf API.

I thought it will be great to run python code inside golang, if possible I will find the way
to get productivity and performance together. This is the result to find out
the way to embed python in golang.

These tests and proof-of-concepts will be continued.


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

```
$ go get github.com/spikeekips/embedding-python-in-golang
```

ready to code.


### Test

#### Lock the python-related code with `sync.Mutex`

```
$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make sync_mutex.go
```

For more detailed, see [sync_mutex/README.md](sync_mutex/README.md)


#### Using Pthread, free from multi-threads problem, based on go-pthread and go-python

```
$ cd src/github.com/spikeekips/embedding-python-in-golang

$ make pthreads.go
```

For more detailed, see [pthreads/README.md](pthreads/README.md)


#### Using Pthread, free from multi-threads problem, based on `cgo API`, without go-pthread and go-python

```
$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make pthreads2.go
```

For more detailed, see [pthreads2/README.md](pthreads2/README.md)


#### Call the golang function in python code.

...


#### Running WSGI server inside Golang's `net/http`

For more detailed, see [wsgi-simple/README.md](wsgi-simple/README.md)



