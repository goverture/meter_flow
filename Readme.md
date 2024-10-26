# MeterFlow

Regulate the flow of your API calls to avoid rate limit errors.

# Register a resource

A resource is a rate limited endpoint, like an API endpoint.

```
curl -X POST http://localhost:8080/resources \
-H "Content-Type: application/json" \
-d '{"name": "openai_api"}'
```

# Schedule your API calls

Schedule your API calls to the resource you registered.

```
curl -X POST http://localhost:8080/schedule \
-H "Content-Type: application/json" \
-d '{"resource_name": "openai_api", "num_calls": 250}' 
```