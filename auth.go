package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

// consider using a router package in the future, keep it simple for now
func loginHandler(w http.ResponseWriter, r *http.Request) {
	//grab the url and split it into segments
	segs := strings.Split(r.URL.Path, "/")
	// url should look like .../auth/login/google
	// maybe should strip out auth then split?
	fmt.Print(len(segs))
	fmt.Print(segs[len(segs)-1])
	if len(segs) < 4 || (len(segs) == 4 && (segs[len(segs)-1]) == "") {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Nop")
		return
	}
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Print("TODO")
		if !isSupportedProvider(&provider) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Login action for provider %s is not supported", provider)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Login action for provider %s is under construction", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s is not supported", action)
	}
}

func isSupportedProvider(s *string) bool {
	supportedProvider := []string{"google", "github", "facebook"}
	for _, provider := range supportedProvider {
		if provider == *s {
			return true
		}
	}

	return false
}

func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}
