package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"

	// Importing the plugins enables collection of AWS resource information at runtime
	// _ "github.com/aws/aws-xray-sdk-go/plugins/ec2"
	// _ "github.com/aws/aws-xray-sdk-go/plugins/beanstalk"
	// _ "github.com/aws/aws-xray-sdk-go/plugins/ecs"

	log "github.com/cihub/seelog"
	"github.com/emicklei/go-restful"
	_ "net/http/pprof"
)

var (
	//CONTEXT is context
	CONTEXT context.Context
)

func init() {
	xray.Configure(xray.Config{
		DaemonAddr:     "xray:2000", // default
		LogLevel:       "info",      // default
		ServiceVersion: "1.2.3",
	})

	CONTEXT = context.Background()
}

func main() {
	go func() {
		log.Info(http.ListenAndServe(":6060", nil))
	}()

	ws := new(restful.WebService)
	ws.Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/true").To(trueFunction))
	restful.Add(ws)

	xray.Capture(CONTEXT, "Internals", func(ctx1 context.Context) (err error) {
		log.Info("Test")
		xray.AddMetadata(ctx1, "ResourceResult", "Test")
		return
	})

	log.Info("Failed webserver", http.ListenAndServe(":8080", nil))
}

func trueFunction(request *restful.Request, response *restful.Response) {
	response.WriteEntity("Here")
}
