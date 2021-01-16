package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.Body)
		writer.WriteHeader(http.StatusOK)
	})
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
