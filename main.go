package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

var (
	httpAddr  = flag.String("http", ":8080", "Host:Port to listen on for HTTP requests")
	httpsAddr = flag.String("https", ":8443", "Host:Port to listen on for HTTPS requests")
	gogsAddr  = flag.String("gogsHttp", "127.0.0.1:3000", "Host:Port of the Gogs to forward to")

	certFile = flag.String("cert", "", "SSL certificate file")
	keyFile  = flag.String("key", "", "SSL private key file")
)

func main() {
	flag.Parse()

	gogs := http.Handler(http.HandlerFunc(gogsNotInstalled))
	if *gogsAddr != "" {
		proxy, err := gogsProxy(*gogsAddr)
		if err != nil {
			log.Fatalf("Failed to create gogs proxy: %v", err)
		}

		gogs = proxy
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Host, "git.") {
			gogs.ServeHTTP(w, r)
			return
		}

		w.Write([]byte("Hello world"))
	})
	http.Handle("/.well-known/", http.FileServer(http.Dir("./letsencrypt")))

	// Start SSL if there is a port set
	if *httpsAddr != "" {
		go func() {
			http.ListenAndServeTLS(*httpsAddr, *certFile, *keyFile, nil)
		}()
	}

	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		panic(err)
	}
}
