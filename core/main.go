package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/techsysfr/paastek-poc/bo"
)

type configuration struct {
	PricingService string `envconfig:"PRICING_SERVICE" required:"true"`
	ListenAddress  string `envconfig:"LISTEN_ADDRESS" required:"true"`
}

type handler struct {
	pricing bo.PricingClient
}

func (h *handler) getLineItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := &bo.ItemID{
		IdentityLineItemID: ps.ByName("id"),
	}
	output, err := h.pricing.ListItem(context.TODO(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {
	var config configuration
	err := envconfig.Process("PAASTEK", &config)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial(config.PricingService, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot dial grpc service ", err)
	}
	pricing := bo.NewPricingClient(conn)
	h := &handler{
		pricing: pricing,
	}

	router := httprouter.New()
	router.GET("/lineItem/:id", h.getLineItem)

	log.Fatal(http.ListenAndServe(config.ListenAddress, router))
}
