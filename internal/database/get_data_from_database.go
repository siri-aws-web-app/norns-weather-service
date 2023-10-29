package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/siri-aws-web-app/verdandi-weather-service/internal/utils"
)

// Declare the config variable package wide here
var cfg aws.Config
var err error

type TableType string

const (
	RealTime TableType = "real-time"
	Forecast TableType = "forecast"
)

// List all tables in the database

func GetCurrentWeatherDataFromDb(cities []string) ([]byte, error) {
	// Load default config from utils
	cfg, err = utils.LoadAwsDefaultConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	wd, err := QueryInputDb(cfg, cities, 1, RealTime)
	if err != nil {
		log.Fatalf("failed to query, %v", err)
	}

	jwd, err := json.Marshal(wd)
	if err != nil {
		log.Fatalf("failed to marshal, %v", err)
	}

	return jwd, nil
}

func GetForecastDataFromDb(cities []string) ([]byte, error) {
	// Load default config from utils
	cfg, err = utils.LoadAwsDefaultConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	wd, err := QueryInputDb(cfg, cities, 1, Forecast)
	if err != nil {
		log.Fatalf("failed to query, %v", err)
	}

	jwd, err := json.Marshal(wd)
	if err != nil {
		log.Fatalf("failed to marshal, %v", err)
	}

	return jwd, nil
}

func QueryInputDb(cfg aws.Config, cities []string, limit int, tableType TableType) (map[string]interface{}, error) {
	client := dynamodb.NewFromConfig(cfg)

	wd := make(map[string]interface{})

	results := make(chan struct {
		city string
		data []map[string]types.AttributeValue
		err  error
	}, len(cities))

	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()

			cityLowercase := strings.ToLower(city)

			var tableToQuery string

			switch tableType {
			case RealTime:
				tableToQuery = fmt.Sprintf("%s-real-time-weather", cityLowercase)
			case Forecast:
				tableToQuery = fmt.Sprintf("%s-forecast", cityLowercase)
			default:
				err = fmt.Errorf("invalid table type")
				results <- struct {
					city string
					data []map[string]types.AttributeValue
					err  error
				}{city, nil, err}
				return
			}

			input := &dynamodb.QueryInput{
				TableName:              aws.String(tableToQuery),
				KeyConditionExpression: aws.String("city = :cityValue"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":cityValue": &types.AttributeValueMemberS{Value: cityLowercase},
				},
				ScanIndexForward: aws.Bool(false),
				Limit:            aws.Int32(int32(limit)),
			}

			q, err := client.Query(context.Background(), input)
			if err != nil {
				results <- struct {
					city string
					data []map[string]types.AttributeValue
					err  error
				}{city, nil, err}
			}

			results <- struct {
				city string
				data []map[string]types.AttributeValue
				err  error
			}{city, q.Items, nil}
		}(city)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			return nil, result.err
		}
		if len(result.data) > 0 {
			wd[result.city] = result.data
		}
	}

	return wd, nil
}
