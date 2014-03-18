package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.net/websocket"
)

type executor func(string, *http.Request)

func main() {
	webroot := *flag.String("webroot", filepath.Join(".", "web"),
		"directory to serve")
	listen := *flag.String("listen", ":1444", "address to listen to")
	flag.Parse()

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(webroot, r.URL.Path)
		action := strings.TrimPrefix(r.URL.Path, "/a/")
		if action != r.URL.Path {
			log.Println("[action]", action)

			path = webroot

			if executor, ok := config[action]; ok {
				executor(action, r)
			} else {
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
	})

	http.Handle("/s/", websocket.Handler(func (ws *websocket.Conn) {
		i := make(chan []byte)

		addClient(i)

		log.Println("[client] addr", ws.RemoteAddr())

		for {
			_, err := ws.Write(<-i)
			if err != nil {
				log.Println("[client] err", err)
				delClient(i)
				return
			}
		}
	}))

	runNotifier()

	log.Println("[server] starting listener at", listen)
	log.Println("[server] serving", webroot)
	log.Fatal(http.ListenAndServe(listen, nil))
}
