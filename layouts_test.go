package main

import "testing"

func TestSiteURL(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{name: "root asset", path: "_syntax.css", want: "/_syntax.css"},
		{name: "root path", path: "/_static/logo.svg", want: "/_static/logo.svg"},
		{name: "base asset", basePath: "/moat", path: "_syntax.css", want: "/moat/_syntax.css"},
		{name: "base root", basePath: "/moat/", want: "/moat/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := siteURL(tt.basePath, tt.path); got != tt.want {
				t.Fatalf("siteURL(%q, %q) = %q, want %q", tt.basePath, tt.path, got, tt.want)
			}
		})
	}
}
