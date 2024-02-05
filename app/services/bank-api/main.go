package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("from bank-api")
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	}))
}
