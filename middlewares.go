package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
)

// DenyListMiddleware is middleware to deny list
func DenyListMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		denyListenv := os.Getenv("DENY_LIST")
		denyList := []string{}
		if err := json.Unmarshal([]byte(denyListenv), &denyList); err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			next.ServeHTTP(w, r)
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			next.ServeHTTP(w, r)
		}
		for _, dnl := range denyList {
			if ip == dnl {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
