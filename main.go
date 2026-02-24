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
	"path/filepath"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "Usage: moat build <src> <dst> [--config PATH] [--site-name NAME] [--base-path PATH]\n")
			os.Exit(1)
		}
		src := os.Args[2]
		dst := os.Args[3]
		siteName := ""
		basePath := ""
		configPath := ""
		hasSiteName := false
		hasBasePath := false
		for i, arg := range os.Args {
			if arg == "--site-name" && i+1 < len(os.Args) {
				siteName = os.Args[i+1]
				hasSiteName = true
			}
			if arg == "--base-path" && i+1 < len(os.Args) {
				basePath = os.Args[i+1]
				hasBasePath = true
			}
			if arg == "--config" && i+1 < len(os.Args) {
				configPath = os.Args[i+1]
			}
		}

		// Auto-detect config.toml in src directory if --config not given
		if configPath == "" {
			candidate := filepath.Join(src, "config.toml")
			if _, err := os.Stat(candidate); err == nil {
				configPath = candidate
			}
		}

		// Load config
		cfg, err := LoadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		// Config file provides defaults; CLI flags override
		if !hasSiteName && cfg.SiteName != "" {
			siteName = cfg.SiteName
		}
		if !hasBasePath && cfg.BasePath != "" {
			basePath = cfg.BasePath
		}

		cfg.SiteName = siteName
		cfg.BasePath = basePath
		if err := Build(src, dst, siteName, basePath, cfg); err != nil {
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

	case "init":
		dir := "docs"
		if len(os.Args) >= 3 {
			dir = os.Args[2]
		}
		if err := Init(dir); err != nil {
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
  moat init [dir]                             Scaffold a docs directory (default: docs/)
  moat build <src> <dst> [flags]              Build static site
    --config PATH      Config file (default: <src>/config.toml)
    --site-name NAME   Site name for templates (default: "Site")
    --base-path PATH   URL prefix for GitHub project pages (e.g. /moat)
  moat serve <dir> [--port PORT]              Serve for local preview
  moat version                                Print version
`, version)
}
