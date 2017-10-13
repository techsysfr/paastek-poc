package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/techsysfr/paastek-poc/bo"
	"google.golang.org/grpc/testdata"
)

type configuration struct {
	PricingService     string `envconfig:"PRICING_SERVICE" required:"true"`
	ListenAddress      string `envconfig:"LISTEN_ADDRESS" required:"true"`
	CAFile             string `envconfig:"CA_FILE" default:"google.golang.org/grpc/testdata/ca.pem" required:"true"`
	ServerHostOverride string `envconfig:"SERVER_HOST_OVERRIDE" required:"true" default:"x.test.google.fr"`
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
	usage := flag.Bool("help", false, "Display usage and exit")
	flag.Parse()
	if *usage {
		envconfig.Usage("PAASTEK", &config)
		return
	}
	err := envconfig.Process("PAASTEK", &config)
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.DialOption
	if config.CAFile == "google.golang.org/grpc/testdata/ca.pem" {
		config.CAFile = testdata.Path("ca.pem")
	}
	creds, err := credentials.NewClientTLSFromFile(config.CAFile, config.ServerHostOverride)
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(config.PricingService, opts...)
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
