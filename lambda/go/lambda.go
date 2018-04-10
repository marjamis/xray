package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	CWSVC *cloudwatch.CloudWatch
)

type request struct {
	ID    float64 `json:"id"`
	Value string  `json:"value"`
}

type response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func handler(request request) (response, error) {
	ctx, seg := xray.BeginSegment(context.Background(), "MyGoFunc")
	defer seg.Close(nil)

	log.Println("here")

	sess := session.Must(session.NewSession())
	CWSVC := cloudwatch.New(sess)
	xray.AWS(CWSVC.Client)
	xray.Configure(xray.Config{LogLevel: "info"})

	alarmNamePrefix := "ECS"
	alarminput := &cloudwatch.DescribeAlarmsInput{
		AlarmNamePrefix: &alarmNamePrefix,
	}
	data, _ := CWSVC.DescribeAlarmsWithContext(ctx, alarminput)

	log.Println("data")
	log.Println(data)

	return response{
		Message: fmt.Sprintf("Processed request ID %f", request.ID),
		Ok:      true,
	}, nil
}

func main() {
	lambda.Start(handler)
}
