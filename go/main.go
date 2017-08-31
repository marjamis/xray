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

func init() {
  xray.Configure(xray.Config{
    DaemonAddr:     "xray:2000", // default
    LogLevel:       "trace",     // default
    ServiceVersion: "1.2.3",
  })
}

func main() {
  go func() {
    log.Info(http.ListenAndServe(":6060", nil))
  }()

  ws := new(restful.WebService)
  ws.Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
  ws.Route(ws.GET("/true").To(trueFunction))
  ws.Route(ws.GET("/false").To(falseFunction))
  restful.Add(ws)

  log.Info("Failed webserver", http.ListenAndServe(":8080", nil))
}

func trueFunction(request *restful.Request, response *restful.Response) {
  ctx := context.Background()
  ctx, seg := xray.BeginSegment(ctx, "GoApp")
  defer seg.Close(nil)
  xray.Capture(ctx, "trueFunction", func(ctx1 context.Context) (err error) {
    response.WriteEntity("Here")
    return
  })
}

func falseFunction(request *restful.Request, response *restful.Response) {
  ctx := context.Background()
  ctx, seg := xray.BeginSegment(ctx, "GoApp")
  defer seg.Close(nil)
  xray.Capture(ctx, "falseFunction", func(ctx1 context.Context) (err error) {
    response.WriteEntity("NotHere")
    return
  })
}
