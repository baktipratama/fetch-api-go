package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Response struct {
	Data int64 `json:"data"`
}

func main() {

	//create a new router
	router := mux.NewRouter()

	//specify endpoints, handler functions and HTTP method
	router.HandleFunc("/test", testController).Methods("GET")
	http.Handle("/", router)
	fmt.Println("Server works!")

	//start and listen to requests
	http.ListenAndServe(":8080", router)
}

func testController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fib := fibbonaci(100)
	resp := Response{Data: fib}
	json.NewEncoder(w).Encode(resp)
}

func fibbonaci(n int) int64 {
	f := make([]int64, n+1, n+2)
	if n < 2 {
		f = f[0:2]
	}
	f[0] = 0
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f[n]
}
