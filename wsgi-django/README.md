## Run django `wsgi` application in Golang.

Test for `django` wsgi application.


### installation

install `django` first,

```
$ source bin/activate
$ pip -v install django
$ cd src/github.com/spikeekips/embedding-python-in-golang/wsgi-django/main
$ django-admin.py syncdb --settings=full.settings --pythonpath=. --noinput
```

### Run

```
$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make wsgi-django.go
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


