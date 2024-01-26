package middlewares

import (
	"encoding/json"
	"fmt"
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
			log.Printf("ERROR MARSHAL: %s\n", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("ERROR SPLITHOST: %s\n", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		for _, dnl := range denyList {
			if ip == dnl {
				fmt.Printf("REMOTE IP -> %s not authorized\n", ip)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// AllowListMiddleware is middleware to allow list
func AllowListMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowListenv := os.Getenv("ALLOW_LIST")
		allowList := []string{}
		if err := json.Unmarshal([]byte(allowListenv), &allowList); err != nil {
			log.Printf("ERROR MARSHAL: %s\n", err.Error())
			next.ServeHTTP(w, r)
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("ERROR SPLITHOST: %s\n", err.Error())
			next.ServeHTTP(w, r)
		}
		for _, alw := range allowList {
			if ip == alw {
				next.ServeHTTP(w, r)
				return
			}
		}
		fmt.Printf("REMOTE IP -> %s not authorized\n", ip)
		w.WriteHeader(http.StatusForbidden)
		return
	})
}
