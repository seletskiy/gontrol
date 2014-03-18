package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var config = map[string]executor{
	"notify": func (name string, r *http.Request) {
		js, err := json.Marshal(r.URL.Query())
		if err != nil {
			log.Fatal("[notify] err", err)
		}

		sendNotify(js)
	},
}
