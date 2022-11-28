package main

import "net/http"

func main() {
	http.HandleFunc("/company", func(writer http.ResponseWriter, request *http.Request) {

	})
	http.HandleFunc("/student", func(writer http.ResponseWriter, request *http.Request) {

	})
}
