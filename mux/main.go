package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte("Hello, " + vars["name"] + "!"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", YourHandler)
	r.HandleFunc("/home/{name}", HomeHandler)

	port := "8082"
	_ = exec.Command("open", "http://localhost:"+port).Run()
	log.Fatal(http.ListenAndServe(":"+port, r))
}
