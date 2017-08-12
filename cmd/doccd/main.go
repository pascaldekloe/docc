package main // import "docc.io/source/cmd/doccd"

import (
	"log"
	"net/http"

	"docc.io/source/site"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", site.Router))
}
