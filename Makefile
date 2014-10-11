.PHONY = clean ab ab.go curl curl.go wsgi-django.uwsgi

PWD = $(shell pwd)
ROOT = $(shell (cd $(PWD)/../../../../; pwd))
DIRECTORY_CLEAN = $(ROOT)/lib $(PWD)

TARGET = $(shell echo $@ | awk -F'.' '{print $$1}')
SUB_TARGET=$(shell echo $@ | awk -F'.' '{print "." $$2 "." $$3}' | sed -e 's/{.}$$//g' -e 's/.go.$$/.go/g' -e 's/^.//g')

DIRECTORY_MAIN = $(PWD)/$(TARGET)/main

PORT = 8080
URL = http://127.0.0.1:$(PORT)

clean: 

	for i in $(DIRECTORY_CLEAN); do \
		find $$i -type f -name "*.pyo" -delete ; \
		find $$i -type f -name "*.pyc" -delete ; \
	done


threadpool.%.go: clean
	export GOPATH=$(ROOT) PYTHONPATH=$(DIRECTORY_MAIN):${PYTHONPATH}; \
	go run -v -x -compiler="gc" $(DIRECTORY_MAIN)/common.go $(DIRECTORY_MAIN)/$(SUB_TARGET);

%.go: clean
	export GOPATH=$(ROOT) PYTHONPATH=$(DIRECTORY_MAIN):${PYTHONPATH}; \
	go run -v -x -compiler="gc" $(DIRECTORY_MAIN)/*.go;


wsgi-django.uwsgi: clean
	uwsgi --virtualenv $(ROOT) --pythonpath $(DIRECTORY_MAIN) --wsgi-file $(DIRECTORY_MAIN)/wsgi.py --chdir $(ROOT) --env DJANGO_SETTINGS_MODULE=full.settings --master --pidfile=$(DIRECTORY_MAIN)/uwsgi.pid --lazy-apps --http=0.0.0.0:$(PORT) --processes=3 --disable-logging --vacuum --enable-threads --max-requests=2000;


ab:
	ab -n 1000 -n 500 $(URL)/;


ab.go:
	ab -n 1000 -n 500 $(URL)/go;


curl:
	curl -H "X-A: $(date)" -v "$(URL)/?a=1&b=2&c=3"; \
	echo;


curl.go:
	curl -H "X-A: $(date)" -v $(URL)/go?a=1&b=2&c=3; \
	echo;


