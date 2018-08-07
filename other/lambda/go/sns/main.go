package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var cwsvc *cloudwatch.CloudWatch

func messageProcessor(ctx context.Context, msg string) error {
	name := "test"
	value := "data"
	dimensions := make([]*cloudwatch.Dimension, 2)
	dimensions[0] = &cloudwatch.Dimension{
		Name:  &name,
		Value: &value,
	}

	input := &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				Dimensions:        dimensions,
				MetricName:        aws.String(string("Task")),
				StorageResolution: aws.Int64(60),
				Timestamp:         aws.Time(time.Now()),
				Unit:              aws.String("Count"),
				Value:             aws.Float64(1),
			},
		},
		Namespace: aws.String("ECS/Task"),
	}
	_, err := cwsvc.PutMetricDataWithContext(ctx, input)
	return err
}

func handler(ctx context.Context, msgs events.SNSEvent) error {
	for _, msg := range msgs.Records {
		if err := messageProcessor(ctx, msg.SNS.Message); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	cwsvc = cloudwatch.New(session.Must(session.NewSession(nil)))
	xray.AWS(cwsvc.Client)
	xray.Configure(xray.Config{LogLevel: "debug"})
	lambda.Start(handler)
}
