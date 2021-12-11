package logger

import (
	"log"
	"net/http"
	"time"
)

func Log(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("start: %v - %s - %s\n", time.Now(), r.Method, r.URL.Path)
		f(w, r)
		log.Printf("end: %v - %s - %s\n", time.Now(), r.Method, r.URL.Path)
	}
}
