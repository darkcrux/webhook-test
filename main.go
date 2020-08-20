package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	address := getAddress()
	fmt.Printf("Started Webhook Test at %s\n", address)
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(address, nil)
}

func getAddress() string {
	if len(os.Args) < 2 {
		fmt.Println("missing address environment, defaulting to :8081")
		return ":8081"
	}
	return os.Args[1]
}

func getUniqueKey() string {
	if len(os.Args) < 3 {
		fmt.Println("missing unique key. stopping.")
		os.Exit(1)
	}
	return os.Args[2]
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	uniqueKey := getUniqueKey()
	fmt.Println("New Notification received")
	// should be POST
	if req.Method != http.MethodPost {
		fmt.Println("Method should be POST")
		res.WriteHeader(http.StatusForbidden)
		return
	}

	uniqueKeyFromHeader := req.Header.Get("WEBHOOK-UNIQUE-KEY")
	idemKeyFromHeader := req.Header.Get("NOTIF-IDEM-KEY")

	if uniqueKey != uniqueKeyFromHeader {
		fmt.Println("Unique Key does not match")
		res.WriteHeader(http.StatusForbidden)
		return
	}

	if idemKeyFromHeader == "" {
		fmt.Println("Idem Key is empty")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("UNIQUE KEY: ", uniqueKeyFromHeader)
	fmt.Println("IDEM KEY: ", idemKeyFromHeader)

	fmt.Println("PAYLOAD:")
	// check if JSON
	var data map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		fmt.Println("Body is not JSON")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	out, _ := json.MarshalIndent(&data, "", "  ")
	fmt.Println(string(out))

	res.WriteHeader(http.StatusOK)
}
