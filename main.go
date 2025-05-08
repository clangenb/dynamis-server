package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := NewRouter()
	fmt.Println("Running on :8080")
	http.ListenAndServe(":8080", r)
}
