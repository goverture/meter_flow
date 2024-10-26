# MeterFlow

Regulate the flow of your API calls to avoid rate limit errors.

# Register a resource

A resource is a rate limited endpoint, like an API endpoint.

Example: register a resource limited to 2 requests per second.
```
curl -X POST -H "Content-Type: application/json" -d '{"name": "rate_limited_resource", "request_count": 2, "time_frame": 1}' http://localhost:8080/resources
```

# Schedule your API calls

Schedule your API calls to the resource you registered.

```
curl -X POST -H "Content-Type: application/json" -d '{"resource_name": "rate_limited_resource", "num_calls": 5}' http://localhost:8080/schedule
```

This return [0,0,1,1,2] which means that the first two calls can be made immediately, the third and fourth calls should be made after 1 second and the fifth call should be made after 2 seconds.

# Benchmark

go-wrk -c 50 -d 10 -M POST -H "Content-Type: application/json" -body '{"resource_name": "rate_limited_resource", "num_calls": 5}' http://localhost:8080/schedule
