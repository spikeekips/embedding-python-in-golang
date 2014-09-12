package main

import (
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/spikeekips/embedding-python-in-golang/wsgi-simple"
	"net/http"
)

var log = logging.MustGetLogger("gowsgi")

func default_handler(w http.ResponseWriter, r *http.Request) {
	log.Debug("> request: default_handler")

	_m, _ := json.Marshal(map[string][]string(r.Header))
	w.Write(_m)
	log.Debug("< request: default_handler")
}

func main() {
	logging.SetFormatter(logging.MustStringFormatter(
		"[%{level}] %{message} (%{module})"))
	logging.SetLevel(logging.INFO, "gowsgi")
	logging.SetLevel(logging.INFO, "epig.thread")
	logging.SetLevel(logging.DEBUG, "gowsgi")
	logging.SetLevel(logging.DEBUG, "epig.thread")

	/*
		go func() {
			ticker := time.NewTicker(3 * time.Second)
			for {
				select {
				case <-ticker.C:
					python.PyType_ClearCache()
					runtime.GC()

					runtime.ReadMemStats(m)
					log.Info("Memory: %v %v %v", m.Alloc, m.StackInuse, m.HeapAlloc)
				}
			}
		}()
	*/

	log.Info("wsgi-simple started.")

	http.HandleFunc("/go", default_handler)
	http.HandleFunc("/", wsgi.WSGIHandler)

	http.ListenAndServe(":8080", nil)
}
