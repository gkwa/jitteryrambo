package jitteryrambo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	LogFormat string `long:"log-format" choice:"text" choice:"json" default:"text" required:"false"`
	Verbose   []bool `short:"v" long:"verbose" description:"Show verbose debug information, each -v bumps log level"`
	logLevel  slog.Level
}

func Execute() int {
	if err := parseFlags(); err != nil {
		return 1
	}

	if err := setLogLevel(); err != nil {
		return 1
	}

	if err := setupLogger(); err != nil {
		return 1
	}

	if err := run(); err != nil {
		slog.Error("run failed", "error", err)
		return 1
	}

	return 0
}

func parseFlags() error {
	_, err := flags.Parse(&opts)
	return err
}

func run() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-gov-east-1"))
	if err != nil {
		return err
	}

	// Create an EC2 client
	client := ec2.NewFromConfig(cfg)

	// Attempt to list EC2 instances in the region
	resp, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		// Check if the error is an API error
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			fmt.Println("API Error Code:", apiErr.ErrorCode())
			fmt.Println("API Error Message:", apiErr.ErrorMessage())

			// Handle specific error codes, e.g., "AccessDeniedException"
			if apiErr.ErrorCode() == "AccessDeniedException" {
				fmt.Println("Unauthorized operation:", apiErr.ErrorMessage())
				// Handle unauthorized operation error here
			} else {
				// Handle other API errors
				fmt.Println("Other API Error:", err.Error())
			}
		} else {
			// Handle non-API errors
			fmt.Println("Error:", err.Error())
		}
		return err
	}

	// Handle success
	fmt.Println("Successfully listed EC2 instances in the region:")
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println("Instance ID:", *instance.InstanceId)
			// Add more details as needed
		}
	}

	return nil
}
