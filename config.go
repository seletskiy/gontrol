package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var config = map[string]executor{
	"notify": func(name string, r *http.Request) {
		js, err := json.Marshal(r.URL.Query())
		if err != nil {
			log.Fatal("[notify] err", err)
		}

		sendNotify(js)
	},

	"scan": func(name string, r *http.Request) {
		out, err := os.Create(time.Now().String() + ".tiff")
		if err != nil {
			log.Fatal("[scan] err", err)
			return
		}

		defer out.Close()

		cmd := exec.Command("scanimage", "--format=tiff")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal("[scan] err", err)
			return
		}

		if cmd.Start() != nil {
			log.Fatal("[scan] err", err)
			return
		}

		io.Copy(out, stdout)

		if cmd.Wait() != nil {
			log.Fatal("[scan] err", err)
			return
		}
	},

	"type": func(name string, r *http.Request) {
		key := r.URL.Query()["key"][0]
		modifier := r.URL.Query()["mod"][0]
		keycode := modifier + "+" + key

		log.Println("[type]  ", keycode)

		cmd := exec.Command("xdotool", "key", keycode)
		cmd.Run()
	},
}
