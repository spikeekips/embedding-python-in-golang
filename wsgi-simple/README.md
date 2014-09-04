## Run simple `wsgi` application in Golang.

Run http server.

```
$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make wsgi-simple.go
...
```

Send test request,

```
$ curl -H "X-A: $(date)" -v http://127.0.0.1:8080/
```

or run `ab`,

```
$ make ab
```

For comparison with native `net/http` in golang, I added go handler, you can
send request to go handler,

```
$ curl -H "X-A: $(date)" -v http://127.0.0.1:8080/go

or 

$ make ab.go
```

