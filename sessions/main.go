package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret"))

func myHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	foo := session.Values["foo"]
	if foo == nil {
		foo = "not set"
	}
	fooString, _ := foo.(string) // Add type assertion to convert foo to string
	data := struct {
		Foo string
	}{
		Foo: fooString, // Use the converted fooString
	}
	renderTemplate(w, "home.html", data)
}

func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template parsing error: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", myHandler)
	http.HandleFunc("/home", homeHandler)

	port := "8081"
	_ = exec.Command("open", "http://localhost:"+port).Run()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
