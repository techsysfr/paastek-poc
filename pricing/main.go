package main

import (
	"context"
	"log"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kelseyhightower/envconfig"
	"github.com/techsysfr/paastek-poc/bo"

	"google.golang.org/grpc"
)

type configuration struct {
	ListenAddress string `envconfig:"ADDRESS" required:"true"`
}

type pricingService struct {
	svc *dynamodb.DynamoDB
}

func (p *pricingService) ListItem(_ context.Context, itemID *bo.ItemID) (*bo.LineItem, error) {
	lineItemID := itemID.IdentityLineItemID
	// Interroge dynamodb pour récupérer l'element qui correspond à l'ID

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"IdentityLineItemID": {
				S: aws.String(lineItemID),
			},
		},
		TableName: aws.String("pricing"),
	}

	result, err := p.svc.GetItem(input)
	if err != nil {
		return nil, err
	}
	var outputBO bo.LineItem
	err = dynamodbattribute.UnmarshalMap(result.Item, &outputBO)

	if err != nil {
		return nil, err
	}
	return &outputBO, nil
}

func main() {
	var config configuration
	err := envconfig.Process("PRICING", &config)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	myServer := grpc.NewServer()
	// Create the session for dynamodb
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("failed to create session,", err)
	}

	svcDynamoDB := dynamodb.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	service := &pricingService{
		svc: svcDynamoDB,
	}
	bo.RegisterPricingServer(myServer, service)
	myServer.Serve(listener)

}
