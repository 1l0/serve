package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s\n", r.RemoteAddr, r.URL.Path)
		w.Header().Add("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
		h.ServeHTTP(w, r)
	})
}

func main() {
	var dir string
	var err error
	if len(os.Args) > 1 {
		if dir, err = filepath.Abs(os.Args[1]); err != nil {
			log.Fatalf("directory not found: %s\n", err)
		}
	} else if dir, err = os.Getwd(); err != nil {
		log.Fatalln(err)
	}
	port := os.Getenv("TLSPORT")
	if port == "" {
		log.Fatalln("no $TLSPORT in env")
	}
	pem := os.Getenv("TLSCERT")
	if pem == "" {
		log.Fatalln("no $TLSCERT in env")
	}
	key := os.Getenv("TLSKEY")
	if key == "" {
		log.Fatalln("no $TLSKEY in env")
	}
	log.Printf("serving at https://localhost:%s\n", port)
	log.Fatalln(
		http.ListenAndServeTLS(
			":"+port, pem, key, handler(http.FileServer(http.Dir(dir))),
		),
	)
}
