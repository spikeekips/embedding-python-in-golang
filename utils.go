package epig

import (
	"encoding/json"
)

// from https://gist.github.com/jd-boyd/119b290f881a0148b515
func PrettyPrintJson(s string) string {
	var _skel []interface{}
	if err := json.Unmarshal([]byte(s), &_skel); err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(_skel, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
