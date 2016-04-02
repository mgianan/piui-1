package main

import "flag"
import "net/http"

var (
	httpAddr = flag.String("http", ":8080", "Host:Port to listen on for HTTP requests")
	httpsAddr = flag.String("https", ":8443", "Host:Port to listen on for HTTPS requests")

	certFile = flag.String("cert", "", "SSL certificate file")
	keyFile = flag.String("key", "", "SSL private key file")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
