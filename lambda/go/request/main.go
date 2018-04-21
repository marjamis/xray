package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var cwsvc *cloudwatch.CloudWatch

type request struct {
	ID    float64 `json:"id"`
	Value string  `json:"value"`
}

type response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func handler(ctx context.Context, req request) error {
	alarmNamePrefix := "ECS"
	alarminput := &cloudwatch.DescribeAlarmsInput{
		AlarmNamePrefix: &alarmNamePrefix,
	}
	data, err := cwsvc.DescribeAlarmsWithContext(ctx, alarminput)
	if err != nil {
		return err
	}

	log.Println("data")
	log.Println(data)

	return nil
}

func main() {
	cwsvc = cloudwatch.New(session.Must(session.NewSession(nil)))
	xray.AWS(cwsvc.Client)
	xray.Configure(xray.Config{LogLevel: "debug"})
	lambda.Start(handler)
}
