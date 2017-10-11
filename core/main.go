package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	productsDB = make(map[string]*product, 0)
	productsDB = map[string]*product{
		"a": &product{
			Family: "mafamille",
			Size:   "micro",
			Hardware: hardware{
				CPU: "i7",
				RAM: 2,
			},
		},
		"b": &product{
			Family: "mafamille",
			Size:   "medium",
			Hardware: hardware{
				CPU: "i7",
				RAM: 4,
			},
		},
		"c": &product{
			Family: "mafamille",
			Size:   "large",
			Hardware: hardware{
				CPU: "i5",
				RAM: 4,
			},
		},
	}
}

func main() {

	router := httprouter.New()
	router.GET("/", rootHandler)
	router.GET("/products", productsHandler)
	router.GET("/products/:ID", productsHandler)
	router.POST("/products", insertProduct)
	router.GET("/articles", articlesHandler)

	log.Fatal(http.ListenAndServe(":8080", router))

}
