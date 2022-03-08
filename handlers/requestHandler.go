package handlers

import (
	"io"
	"net/http"
)

func RequestHandler() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
	http.HandleFunc("/get_map", helloHandler)
	http.ListenAndServe(":33061", nil)
}
