package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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

	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
