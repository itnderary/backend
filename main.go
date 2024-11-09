package main

import (
	"flag"
	"fmt"
	"net/http"
)

var Push bool

func main() {
	http.HandleFunc("/moods", moods)
	http.HandleFunc("/favicon.ico", favicon)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func moods(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func favicon(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func init() {
	flag.BoolVar(&Push, "push", true, "Push stuff")
}
