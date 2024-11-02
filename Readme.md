# MeterFlow

MeterFlow is a lightweight API call scheduler designed to help you manage and regulate the flow of API requests, avoiding rate limit errors across multiple resources. With MeterFlow, you can register resources with customizable request limits to schedule API calls efficiently.

# Getting started

Register a resource (a rate limited entity) by specifying the name, request count, and time frame (ex "openai_api": 100 calls / minute). Then schedule your API calls to the registered resource and MeterFlow will return the time intervals at which you can make the calls. MeterFlow will track the timing of each calls to ensure you never exceed the rate limit of the resource.

## Register a resource
Example: register a resource limited to 2 requests per second.
```
curl -X POST -H "Content-Type: application/json" -d '{"name": "rate_limited_resource", "request_count": 2, "time_frame": 1}' http://localhost:8080/resources
```

## List all resources
```
curl -X GET -H "Content-Type: application/json" http://localhost:8080/resources
```

## Update a resource
```
curl -X PUT -H "Content-Type: application/json" -d '{"name": "rate_limited_resource", "request_count": 3, "time_frame": 1}' http://localhost:8080/resources
```

## Delete a resource
```
curl -X DELETE -H "Content-Type: application/json" -d '{"name": "rate_limited_resource"}' http://localhost:8080/resources
```

# Schedule your API calls

Schedule a number of API calls to the resource you registered. It will return the necessary delay in seconds for each call.

```
curl -X POST -H "Content-Type: application/json" -d '{"resource_name": "rate_limited_resource", "num_calls": 5}' http://localhost:8080/schedule
```

Since the rate limit is 2 requests per second for this resource, the response should be
```json
[0,0,1,1,2]
```
which means that the first two calls can be made immediately, the third and fourth calls should be made after 1 second and the fifth call should be made after 2 seconds.

# Benchmark

Use `go-wrk` (https://github.com/tsliwowicz/go-wrk) to benchmark the server.
```
go-wrk -c 50 -d 10 -M POST -H "Content-Type: application/json" -body '{"resource_name": "rate_limited_resource", "num_calls": 5}' http://localhost:8080/schedule
```