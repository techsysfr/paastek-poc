package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kelseyhightower/envconfig"
	"github.com/techsysfr/paastek-poc/bo"
)

const prefix = "PAASTEK"

type configuration struct {
	TableName string `envconfig:"tablename" required:"true"`
}

type ret struct {
	V   bo.AWSBillingItem
	Err error
}

// parse a csv file and return an array of resources
func parse(r io.Reader) chan ret {
	c := make(chan ret, 0)
	go func() {
		defer close(c)
		rd := csv.NewReader(r)
		var header []string
		header, err := rd.Read()
		if err != nil {
			c <- ret{bo.AWSBillingItem{}, err}
		}

		e := bo.AWSBillingItem{}
		et := reflect.TypeOf(e)
		var headers = make(map[string]int, et.NumField())
		for i := 0; i < et.NumField(); i++ {
			headers[et.Field(i).Name] = func(element string, array []string) int {
				for k, v := range array {
					if v == element {
						return k
					}
				}
				return -1
			}(et.Field(i).Tag.Get("csv"), header)
		}
		for {
			var e = bo.AWSBillingItem{}
			record, err := rd.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c <- ret{bo.AWSBillingItem{}, err}
			}
			for h, i := range headers {
				if i == -1 {
					continue
				}
				elem := reflect.ValueOf(&e).Elem()
				field := elem.FieldByName(h)
				if field.CanSet() {
					switch field.Type().Name() {
					case "float64":
						a, _ := strconv.ParseFloat(record[i], 64)
						field.Set(reflect.ValueOf(a))
					case "Time":
						a, _ := time.Parse("2006-01-02T00:00:00Z", record[i])
						field.Set(reflect.ValueOf(a))
					default:
						field.Set(reflect.ValueOf(record[i]))
					}
				}
			}
			c <- ret{e, nil}
		}
	}()
	return c
}

func main() {
	var config configuration
	err := envconfig.Process(prefix, &config)
	if err != nil {
		envconfig.Usage(prefix, &config)
		log.Fatal(err.Error())
	}

	// Create the session for dynamodb
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("failed to create session,", err)
	}

	svc := dynamodb.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	c := parse(os.Stdin)
	for val := range c {
		item, err := dynamodbattribute.MarshalMap(val.V)
		if err != nil {
			log.Println(err)
		}
		params := &dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(config.TableName),
		}
		log.Println(item)
		// Now put the item, discarding the result
		_, err = svc.PutItem(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err)
		}
	}
}
