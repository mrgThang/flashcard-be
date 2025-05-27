package main

import "net/http"

func main() {
	http.HandleFunc("hello", func(w http.ResponseWriter, r *http.Request) {})
}