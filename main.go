package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/KexinLu/tracer"
)

// define templateHandler conform to http.Handler
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//attaching ServeHttp to templateHandler
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The Address of the Application")
	flag.Parse()
	r := newRoom()
	r.tracer = tracer.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting chat server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
