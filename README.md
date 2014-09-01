## Embedding Python In Golang

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
$ go get github.com/sbinet/go-python
$ go get github.com/liamzdenek/go-pthreads
$ go get github.com/op/go-logging

$ cd src/github.com/spikeekips/embedding-python-in-golang

$ make sync_mutex.go
```


#### Using Pthread, free from multi-threads problem, using go-pthread and go-python

```
$ go get github.com/sbinet/go-python
$ go get github.com/liamzdenek/go-pthreads
$ go get github.com/op/go-logging

$ cd src/github.com/spikeekips/embedding-python-in-golang

$ make pthreads.go
```


#### Using Pthread, free from multi-threads problem, using `cgo API`

```
$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make pthreads2.go
```


