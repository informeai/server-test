package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// ParseXForwardFor is middleware to parse header X-Forward-For
func ParseXForwardFor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xForwardFor := r.Header.Get("X-Forwarded-For")
		if len(xForwardFor) > 0 {
			fmt.Printf("xForwardFor: %v\n", xForwardFor)
			fmt.Printf("remoteAddr: %v\n", r.RemoteAddr)
			fmt.Printf("request: %+v\n", r)

			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// DenyListMiddleware is middleware to deny list
func DenyListMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		denyListenv := os.Getenv("DENY_LIST")
		denyList := []string{}
		xForwardFor := r.Header.Get("X-Forwarded-For")
		if err := json.Unmarshal([]byte(denyListenv), &denyList); err != nil {
			log.Printf("ERROR MARSHAL: %s\n", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		if len(xForwardFor) > 0 {
			for _, dnl := range denyList {
				if xForwardFor == dnl {
					fmt.Printf("REMOTE IP -> %s not authorized\n", xForwardFor)
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
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
		xForwardFor := r.Header.Get("X-Forwarded-For")
		if err := json.Unmarshal([]byte(allowListenv), &allowList); err != nil {
			log.Printf("ERROR MARSHAL: %s\n", err.Error())
			next.ServeHTTP(w, r)
		}
		if len(xForwardFor) > 0 {
			for _, alw := range allowList {
				if xForwardFor == alw {
					next.ServeHTTP(w, r)
					return
				}
			}
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
