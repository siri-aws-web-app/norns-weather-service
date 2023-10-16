package utils

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadAwsDefaultConfig() (aws.Config, error) {
	var cfg aws.Config
	var err error

	if os.Getenv("ENV") == "local" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("eu-central-1"),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		log.Fatalf("Failed to load config, %v", err)
	}

	return cfg, err
}
