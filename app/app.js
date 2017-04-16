var express    = require("express");
var AWSXRay = require('aws-xray-sdk');
var AWS = require('aws-sdk');
var winston = require('winston');
var dns = require('dns');

var app = express();

winston.level = 'debug';
AWSXRay.config([AWSXRay.plugins.EC2]);
AWSXRay.config([AWSXRay.plugins.ECS]);

const options = {
  family: 4,
  hints: dns.ADDRCONFIG | dns.V4MAPPED,
};

dns.lookup('xray', options, (err, address, family) =>
  AWSXRay.setDaemonAddress(address+':2000'));
AWSXRay.setLogger(winston);


app.use(AWSXRay.express.openSegment('webapp'));

app.get("/",function(req,resp){
  app.use(AWSXRay.express.openSegment('responseGeneration'));
  body = "This is the correct webpage\n\n";
  resp.write(body)
  resp.end();
  app.use(AWSXRay.express.closeSegment('responseGeneration'));
});

app.listen(3000);
app.use(AWSXRay.express.closeSegment('webapp'));
