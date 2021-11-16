//go:build include

//go:generate go-include --out main_gen.go $GOFILE

package main

import (
	"fmt"
	"net/http"

	"github.com/halimath/include"
)

var (
	html    = include.Bytes("./index.html")
	htmlStr = include.String("./index.html")
)

func main() {
	fmt.Printf("Listening on :8080...\n")

	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(html)
	}))
}
