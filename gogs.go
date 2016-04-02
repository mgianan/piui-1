package main

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const poolSliceSize = 64 * 1024

func gogsProxy(gogsHostPort string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse("http://" + gogsHostPort + "/")
	if err != nil {
		return nil, err
	}

	p := httputil.NewSingleHostReverseProxy(u)
	p.BufferPool = newBytePool(poolSliceSize)
	return p, nil
}

func gogsNotInstalled(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "Gogs handler is not registered")
}
