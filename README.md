# X-Ray
X-Ray configuration for testing/explanations.

More to come.

## General
All applications are run with Docker allowing for a simple docker-compose file which calls the specific application you want or even all of them.

NOTE: The used docker-compose file does currently require an .env file for the used environment variables.

### Usage
```docker-compose up <service_or_blank> ```

Or:

```export $(cat .env) && 
docker stack deploy --compose-file ./docker-compose <stack-name>```

## Examples provided
### X-Ray daemon
```docker-compose up xray``` //Though this will be started up automatically from any other platform.
* Docker containerised

### Node.js
To run:
```docker-compose up nodejs```
* Uses express
* Traces through a front-end and a back-end application
* Uses sampling based of a configuration file
* Adds custom X-Ray annotations
* Adds customs X-Ray metadata
* Captures an synchronous call
* Capture an asynchronous call
