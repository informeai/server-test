package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	log.Printf("%s", r.URL.Path)
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

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/health", health)
	http.HandleFunc("/test", test)
	fmt.Printf("listen in port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
