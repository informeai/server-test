package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/informeai/server-test/middlewares"
)

type ResponseHealth struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RequestTest map[string]interface{}

type ResponseTest struct {
	Payload map[string]interface{} `json:"payload"`
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("%s", r.URL.Path)
	if err := json.NewEncoder(w).Encode(ResponseHealth{Status: "success", Message: "active"}); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var ip string
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if len(xForwardedFor) > 0 {
		ip = xForwardedFor
	} else {
		ipHost, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}
		ip = ipHost
	}
	log.Printf("IP -> %v\n", ip)
	log.Printf("PATH -> %s\n", r.URL.Path)
	log.Printf("QUERY PARAMS -> %v\n", r.URL.Query())
	log.Printf("HEADERS -> %v\n", r.Header)
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		payload := RequestTest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}
		bytesPayload, err := json.MarshalIndent(&payload, "", " ")
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}
		fmt.Printf("PAYLOAD: %v\n", string(bytesPayload))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	if err := json.NewEncoder(w).Encode(ResponseHealth{Status: "error", Message: "method not allowed"}); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}

func dash(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("%s", r.URL.Path)
	byFile, err := ioutil.ReadFile("./template/index.html")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byFile)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}

func quantityDynamicPath(num int) string {
	path := ""
	for i := 0; i < num; i++ {
		path += fmt.Sprintf("/{%d}", i)
	}
	return path
}

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()
	router.Use(middlewares.ParseXForwardFor, middlewares.AllowListMiddleware)
	router.HandleFunc("/health", health)
	router.HandleFunc("/{first}", test)
	router.HandleFunc("/{first}/{second}", test)
	router.HandleFunc("/{first}/{second}/{terciary}", test)
	router.HandleFunc("/{first}/{second}/{terciary}/{four}", test)
	fmt.Printf("listen in port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
