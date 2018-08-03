require 'sinatra'
require "aws-sdk"
require 'net/http'
require 'aws-xray-sdk'
XRay.recorder.configure({patch: %I[net_http]})

get '/' do
  erb :index
end

get '/getjoke' do
  segment = XRay.recorder.begin_segment 'NarSinny'
  # subsegment = XRay.recorder.begin_subsegment 'GetChucky', namespace: 'remote'
  uri = URI.parse("http://api.icndb.com/jokes/random")
  res = Net::HTTP.get_response(uri)
  nn = JSON.parse(res.body)
  joke = nn['value']['joke']
  # XRay.recorder.end_subsegment
  XRay.recorder.end_segment
  erb :showJoke, :locals => {:theJoke => joke}
end

get '/joke' do
  erb :chuckyIndex
end

get '/listbucket' do
  buck = params[:theBucket]
  segment = XRay.recorder.begin_segment 'NarSinny'
  s3 = Aws::S3::Client.new
  resp = s3.list_objects(bucket: buck)
  buckConts = resp.contents
  XRay.recorder.end_segment
  erb :showBucket, :locals => {:bucket => buckConts, :bucky => buck}
end
