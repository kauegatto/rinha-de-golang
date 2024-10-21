package main

import (
	"fmt"
	"log"
	"net/http"
)

type Payment struct {
	Valor     uint32
	Tipo      string
	Descricao string
}

func (p Payment) validar() bool {
	if p.Tipo != "c" || p.Tipo != "d" {
		return false
	}
}

func main() {
	http.Handle("/payment/{id}", deposit())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func deposit() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			fmt.Fprintf(w, "Hello %s", id)
		},
	)
}
