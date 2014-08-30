.PHONY = clean

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
	go run $(PWD)/$(TARGET)/main.go;


