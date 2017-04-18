var express    = require("express");
var AWSXRay = require('aws-xray-sdk');
var AWS = require('aws-sdk');
var winston = require('winston');
var dns = require('dns');
var http = require('http');
AWSXRay.captureHTTPs(http);

var app = express();

// BEGIN: AWS X-Ray configuration details
var ip;
const options = {
  family: 4,
  hints: dns.ADDRCONFIG | dns.V4MAPPED,
};
dns.lookup('xray', options, (err, address, family) =>
  AWSXRay.setDaemonAddress(address+':2000'));

dns.lookup('nodejs-backend', options, (err, address, family) =>
  ip = address);

AWSXRay.config([AWSXRay.plugins.EC2, AWSXRay.plugins.ECS]);
winston.level = 'debug'; AWSXRay.setLogger(winston);
AWSXRay.middleware.setSamplingRules('./xray_sampling-rules.json');
AWSXRay.middleware.enableDynamicNaming('*');
// END: AWS X-Ray configuration details
console.log("Starting: ");
console.log(process.env.APP_NAME);
app.use(AWSXRay.express.openSegment(process.env.APP_NAME));

app.get("/",function(req,resp){
  var seg = AWSXRay.getSegment();
  seg.addAnnotation('User-Agent', req.get('User-Agent'));
  seg.addMetadata('MetadataKey', 'MetadataValue');

  AWSXRay.capture('responseGeneration-/', function(subsegment){
    body = "This is the correct webpage\n\n";
  });

  //Some supported AWS API call for tracingthat? prob DDB

  var options = {
    host: ip,
    path: '/true',
    port: '3000',
    method: 'GET',
    Segment: seg
  }
  var outbound_req = http.request(options, (res) => {
    console.log('Attempting remote connection');
    res.on('data', (chunk) => {
      console.log(`BODY: ${chunk}`);
    });
    res.on('end', () => {
      console.log('No more data');
    });
  });

  outbound_req.write('temp');
  outbound_req.end();

  AWSXRay.captureAsync('responseWriting-/', function(subsegment){
    resp.write(body)
    resp.end();
    subsegment.close();
  });
});

app.get("/true",function(req,resp){
  console.log('Conection seen...');
  AWSXRay.capture('responseGeneration-/true', function(subsegment){
    body = "true";
  });

  AWSXRay.captureAsync('responseWriting-/true', function(subsegment){
    resp.write(body)
    resp.end();
    subsegment.close();
  });
});

app.listen(3000);
