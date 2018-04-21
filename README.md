# X-Ray
X-Ray configuration for testing/explanations.

**This is still very much a work in-progress and some things will require updating/improving.**

## General
Most of the applications are runnable via the docker-compose.yml file for simplicity.

**NOTE:** The docker-compose.yml file requires an .env file for environment variables used throughout.

### Usage
```bash
docker-compose up <service_or_blank>
```

Or:

```bash
export $(cat .env) && \
docker stack deploy --compose-file ./docker-compose <stack-name>
```

## Examples provided
All applications does/will have:
* Docker containerised
* Traces through a front-end and a back-end application
* Uses sampling based of a configuration file
* Adds custom X-Ray annotations
* Adds customs X-Ray metadata
* Captures an synchronous call
* Capture an asynchronous call

AWS Services does/will trace:
* ALB/ELB
* API Gateway
* DynamoDDB

### X-Ray daemon
```bash
docker-compose up xray # This is a dependency for all applications started via docker-compose.yml
```
### "C#"
PENDING

### Go
To run:
```bash
docker-compose up go
```
* Uses custom Segments due to non-supported frameworks

### Java
Already provided Java example:
https://github.com/awslabs/eb-java-scorekeep/tree/xray

### Node.js
To run:
```bash
docker-compose up nodejs
```
* Uses express

### Ruby
PENDING


### Generic structure of the applications
> / # Entrypoint of webapp

> /true # One internal function

> /false # Second internal function

#### Testing samples
Simple test against an endpoint with 30 concurrent users with a total of 1000 requests.
```bash
ab -k -c 30 -n 1000 <endpoint>/
```
If you're needing additional configurations/outputs the ab man page is the place to be.
