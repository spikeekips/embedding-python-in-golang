.PHONY = clean ab ab.go

PWD=$(shell pwd)
ROOT=$(shell (cd $(PWD)/../../../../; pwd))
DIRECTORY_CLEAN = $(ROOT)/lib $(PWD)
TARGET=$(shell basename $@ .go)

clean: 

	for i in $(DIRECTORY_CLEAN); do \
		find $$i -type f -name "*.pyo" -delete ; \
		find $$i -type f -name "*.pyc" -delete ; \
	done

%.go: clean
	export GOPATH=$(ROOT) PYTHONPATH=$(PWD)/$(TARGET):${PYTHONPATH}; \
	go run -v -x -compiler="gc" $(PWD)/$(TARGET)/*.go;

ab:
	ab -n 1000 -n 500 http://127.0.0.1:8080/;

ab.go:
	ab -n 1000 -n 500 http://127.0.0.1:8080/go;


