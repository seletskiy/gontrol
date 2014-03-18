package main

import (
	"encoding/json"
	"log"
)

var config = actionsConfig{
	"notify": notifyAction{},
}

func PerformAction(name string, args map[string][]string) bool {
	if action, ok := config[name]; ok {
		action.exec(name, args)
		return true
	} else {
		return false
	}
}

type notifyAction struct {
}

func (n notifyAction) exec(name string, args map[string][]string) {
	js, err := json.Marshal(args)
	if err != nil {
		log.Fatal("[notify]", err)
	}

	sendNotify(js)
}
