// moat — markdown + oat static site generator
//
// Usage:
//
//	moat build <src> <dst>    Build static site from markdown source
//	moat serve <dir> [--port] Serve static files for local preview
package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "Usage: moat build <src> <dst> [--site-name NAME] [--base-path PATH]\n")
			os.Exit(1)
		}
		src := os.Args[2]
		dst := os.Args[3]
		siteName := ""
		basePath := ""
		for i, arg := range os.Args {
			if arg == "--site-name" && i+1 < len(os.Args) {
				siteName = os.Args[i+1]
			}
			if arg == "--base-path" && i+1 < len(os.Args) {
				basePath = os.Args[i+1]
			}
		}
		if err := Build(src, dst, siteName, basePath); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "serve":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: moat serve <dir> [--port PORT]\n")
			os.Exit(1)
		}
		dir := os.Args[2]
		port := "8080"
		for i, arg := range os.Args {
			if arg == "--port" && i+1 < len(os.Args) {
				port = os.Args[i+1]
			}
		}
		if err := Serve(dir, port); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "version":
		fmt.Printf("moat %s\n", version)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `moat %s — markdown + oat static site generator

Usage:
  moat build <src> <dst> [flags]               Build static site
    --site-name NAME   Site name for templates (default: "Site")
    --base-path PATH   URL prefix for GitHub project pages (e.g. /moat)
  moat serve <dir> [--port PORT]              Serve for local preview
  moat version                                Print version
`, version)
}
