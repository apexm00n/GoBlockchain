package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func JsonStatus (message string) []byte {
	m, _ := json.Marshal(struct{
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}