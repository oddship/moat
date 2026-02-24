package main

import (
	"fmt"
	"net/http"
)

// Serve starts a static file server for local preview.
func Serve(dir, port string) error {
	fmt.Printf("Serving %s on http://localhost:%s\n", dir, port)
	return http.ListenAndServe(":"+port, http.FileServer(http.Dir(dir)))
}
