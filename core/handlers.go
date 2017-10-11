package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello root")
}

func productsHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	if param != nil {

		prd, err := getProduct(param.ByName("ID"))
		if err != nil {
			log.Fatal(err)
		}
		if prd == nil {
			http.NotFound(w, r)
			return
		}
		aff, err := json.MarshalIndent((*prd), "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		log.Println(prd)
		fmt.Fprintf(w, "Hello product, %s!\n", string(aff))

	} else {
		prd, err := getProducts()
		if err != nil {
			log.Fatal(err)
		}
		aff, err := json.MarshalIndent((prd), "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Hello product, %s!\n", string(aff))
	}

}

func insertProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var prd product
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&prd)

	fmt.Fprintf(w, "%v\n", err)
	fmt.Fprintf(w, "%v\n", prd)
}

func articlesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello articles")
}
