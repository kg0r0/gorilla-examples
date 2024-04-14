package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"text/template"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	CSRFMiddleware := csrf.Protect(
		[]byte("place-your-32-byte-long-key-here"),
		csrf.Secure(false),                 // false in development only!
		csrf.RequestHeader("X-CSRF-Token"), // Must be in CORS Allowed and Exposed Headers
	)

	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/api", Get).Methods(http.MethodGet)
	r.HandleFunc("/api", Post).Methods(http.MethodPost)

	port := "8083"
	_ = exec.Command("open", "http://localhost:"+port).Run()
	log.Fatal(http.ListenAndServe(":"+port, CSRFMiddleware(r)))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").ParseFiles("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Template parsing error: %v", err), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-CSRF-Token", csrf.Token(r))
	w.WriteHeader(http.StatusOK)
}

func Post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
