## Embedding Python In Golang

This is the tests to embed python code in golang runtime and proof-of-concepts,
and will be added continuously.


### Installation

create base environment using `virtualenv` of python.

```
$ pip install virtualenv
```

```
$ virtualenv ~/gopython
$ cd ~gopython
$ mkdir -p {pkg,src}
$ export GOPATH=$HOME/gopython
```

#### OSX

If you are using `homebrew`, install `pkg-config` and write `python-2.7.pc` file manually. It needs.

```
$ brew -v install pkg-config
```

and then, write new file, `/usr/local/lib/pkgconfig/python-2.7.pc`

```
prefix=/System/Library/Frameworks/Python.framework/Versions/2.7
exec_prefix=${prefix}
libdir=${exec_prefix}/lib
includedir=${prefix}/include

Name: Python
Description: Python library
Requires:
Version: 2.7
Libs.private: -lpthread -ldl  -lutil
Libs: -L${libdir} -ldl -lpython2.7
Cflags: -I${includedir}/python2.7 -I${includedir}/python2.7 -fno-strict-aliasing -fno-common -dynamic -g -Os -pipe -fno-common -fno-strict-aliasing -fwrapv -DENABLE_DTRACE -DMACOSX -DNDEBUG -Wall -Wstrict-prototypes -Wshorten-64-to-32 -DNDEBUG -g -fwrapv -Os -Wall -Wstrict-prototypes -DENABLE_DTRACE
```


#### Install Source
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


#### Using Pthread, free from multi-threads problem, based on `cgo API`

```
$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make pthreads2.go
```

For more detailed, see [pthreads2/README.md](pthreads2/README.md)


#### Running WSGI application inside Golang's `net/http`

For more detailed, see [wsgi-simple/README.md](wsgi-simple/README.md)


#### Running Django WSGI application inside Golang's `net/http`

For more detailed, see [wsgi-django/README.md](wsgi-django/README.md)


#### Call the golang function in python code.

...



