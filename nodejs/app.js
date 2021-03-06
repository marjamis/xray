var express = require('express');
var AWSXRay = require('aws-xray-sdk');
var AWS = require('aws-sdk');
var winston = require('winston');
var dns = require('dns');
var http = require('http');
var tracedHttp = AWSXRay.captureHTTPs(require('http'));

var app = express();

// BEGIN: AWS X-Ray configuration details
var ip;
const options = {
  family: 4,
  hints: dns.ADDRCONFIG | dns.V4MAPPED,
};
dns.lookup('xray-daemon', options, (err, address, family) =>
  AWSXRay.setDaemonAddress(address+':2000'));

dns.lookup('nodejs-backend', options, (err, address, family) =>
  ip = address);

//AWSXRay.config([AWSXRay.plugins.EC2, AWSXRay.plugins.ECS]);
winston.level = 'debug'; AWSXRay.setLogger(winston);
AWSXRay.middleware.setSamplingRules('/xray_sampling_rules.json');
//AWSXRay.middleware.enableDynamicNaming('*');
// END: AWS X-Ray configuration details
console.log("Starting: " + process.env.APP_NAME);
app.use(AWSXRay.express.openSegment(process.env.APP_NAME));

app.get("/",function(req,resp){
  var seg = AWSXRay.getSegment();
  seg.addAnnotation('User-Agent', req.get('User-Agent'));
  seg.addMetadata('MetadataKey', 'MetadataValue', 'general');

  AWSXRay.captureFunc('responseGeneration-/', function(subsegment){
    body = "This is the correct webpage\n\n";
  });

  //Some supported AWS API call for tracingthat? prob DDB

  var options = {
    host: ip,
    path: '/data',
    port: '3000',
    method: 'GET'
  }
  var outbound_req = tracedHttp.request(options, (res) => {
    res.on('data', (chunk) => {
      console.log(`BODY: ${chunk}`);
    });
    res.on('end', () => {
      console.log('No more data');
    });
  });
  outbound_req.write('temp');
  outbound_req.end();

  AWSXRay.captureAsyncFunc('responseWriting-/', function(subsegment){
    resp.write(body)
    resp.end();
    subsegment.close();
  });
});

app.get("/data",function(req,resp){
  var seg = AWSXRay.getSegment();
  seg.addAnnotation('User-Agent', req.get('User-Agent'));
  seg.addMetadata('MetadataKey', 'MetadataValue', 'general');
  AWSXRay.captureFunc('responseGeneration-/data', function(subsegment){
    body = '{"data": true}';
  });

  AWSXRay.captureAsyncFunc('responseWriting-/data', function(subsegment){
    resp.write(body)
    resp.end();
    subsegment.close();
  });
});

app.listen(3000);
