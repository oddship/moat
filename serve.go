package main

import (
	"fmt"
	"net"
	"net/http"
)

// defaultServeHost keeps local previews private by default. Change this
// default or add a host flag if exposing the preview server is needed later.
const defaultServeHost = "127.0.0.1"

// Serve starts a static file server for local preview.
func Serve(dir, port string) error {
	addr := net.JoinHostPort(defaultServeHost, port)
	fmt.Printf("Serving %s on http://%s\n", dir, addr)
	return http.ListenAndServe(addr, http.FileServer(http.Dir(dir)))
}
