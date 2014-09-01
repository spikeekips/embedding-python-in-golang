## Embedding Python In Golang

I tested how to embed python code in golang, if possible, It brings the
combination of productivity and performance.

### Installation

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

ready to go.


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


