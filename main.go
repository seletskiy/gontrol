package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.net/websocket"
)

type actionsConfig map[string]executor
type executor interface {
	exec(string, map[string][]string)
}

type actionHandler struct {
	root string
}

func (h actionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.root, r.URL.Path)
	action := strings.TrimPrefix(r.URL.Path, "/a/")
	if action != r.URL.Path {
		log.Println("[action]", action)

		path = h.root

		if !PerformAction(action, r.URL.Query()) {
			http.NotFound(w, r)
			return
		}
	}

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
	} else {
		if stat.IsDir() {
			path = filepath.Join(path, "index.html")
		}

		log.Println("[serve] ", path)

		file, err := os.Open(path)
		if err != nil {
			log.Fatalln(err)
		} else {
			io.Copy(w, file)
		}
	}
}

func notifyHandler(ws *websocket.Conn) {
	i := make(chan []byte)

	addClient(i)

	log.Println("[client] addr", ws.RemoteAddr())

	for {
		_, err := ws.Write(<-i)
		if err != nil {
			log.Println("[client] err", err)
			delClient(i)
		}
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	webdir := filepath.Join(cwd, "web")
	http.Handle("/", actionHandler{webdir})

	runNotifier()
	http.Handle("/s/", websocket.Handler(notifyHandler))

	log.Println("[server] starting listener at", ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
