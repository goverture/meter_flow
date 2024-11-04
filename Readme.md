# MeterFlow

MeterFlow is a lightweight API call scheduler designed to help you manage and regulate the flow of outgoing API requests, avoiding rate limit errors. With MeterFlow, you can register resources with customizable request limits to schedule API calls efficiently and safely.

## Example usage

```ruby
# Step 0 (only once): Register an API resource with MeterFlow, for instance "dummy_api" with 100 calls per minute
uri = URI("http://localhost:8080/resources") # Assuming you are running MeterFlow locally on port 8080
Net::HTTP.post(uri, { name: 'dummy_api', request_count: 100, time_frame: 60 }.to_json, "Content-Type" => "application/json")

# Step 1: Request the schedule from MeterFlow
uri = URI("http://localhost:8080/schedule")
response = Net::HTTP.post(uri, { resource_name: 'dummy_api', num_calls: 1000 }.to_json, "Content-Type" => "application/json")

# Step 2: Parse the response and enqueue jobs based on the delay
delays = JSON.parse(response.body)['delays']
delays.each_with_index do |delay, index|
  DummyApiCallWorker.perform_in(delay, index + 1)
end
```

## Features

Supported rate limiting algorithms:
- [x] Only "sliding window" (X calls in the past time frame) is supported at the moment.
- [ ] TODO: other algorithms ("token bucket" etc)

Supported limits:
- [x] Only the number of requests per time frame is supported at the moment.
- [ ] TODO: Support LLM "token per minute" limits.

Persistence
- [x] Save the registered resources to disk upon exist

## Getting started

Register a resource (a rate limited entity) by specifying the name, request count, and time frame (ex "openai_api": 100 calls / minute). Then schedule your API calls to the registered resource and MeterFlow will return the time intervals at which you can make the calls. MeterFlow will track the timing of each calls to ensure you never exceed the rate limit of the resource. Check the [resources wiki page](https://github.com/goverture/meter_flow/wiki/Resources) for more details.

```
curl -X POST -H "Content-Type: application/json" -d '{"name": "rate_limited_resource", "request_count": 2, "time_frame": 1}' http://localhost:8080/resources
```

Then schedule a number of API calls to the resource you registered. It will return the necessary delay in seconds for each call.

```
curl -X POST -H "Content-Type: application/json" -d '{"resource_name": "rate_limited_resource", "num_calls": 5}' http://localhost:8080/schedule
```

Since the rate limit is 2 requests per second for this resource, the response should be
```json
{"delays":[0,0,1,1,2]}
```
which means that the first two calls can be made immediately, the third and fourth calls should be made after 1 second and the fifth call should be made after 2 seconds.
