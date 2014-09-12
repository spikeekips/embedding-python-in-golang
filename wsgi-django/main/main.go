package main

import (
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/spikeekips/embedding-python-in-golang/wsgi-django"
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
	//logging.SetLevel(logging.INFO, "wsgi")
	logging.SetLevel(logging.DEBUG, "wsgi")

	log.Info("wsgi-django started.")

	http.HandleFunc("/go", default_handler)
	http.HandleFunc("/", wsgi.WSGIHandler)

	http.ListenAndServe(":8080", nil)
}
