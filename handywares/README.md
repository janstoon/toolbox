# Handywares [![Coverage Badge][bdg-cov-handywares]][action-tests]

Reusable middleware for HTTP, gRPC, NATS, and Asynq, plus mappers between `bricks` errors
and HTTP statuses, gRPC status codes, and Redis client errors.
Full documentation and examples are available in [pkg.go.dev][go-pkg].

# Usage

```shell
go get github.com/janstoon/toolbox/handywares
```

## HTTP

`HttpMiddlewareStack` chains `net/http` handlers. `HttpPanicRecoverMiddleware` turns a panic
into HTTP 500. `HttpCrossOriginResourceSharingPolicyMiddleware` applies a CORS policy built
with `CorsAllowOrigins`, `CorsAllowMethods`, and `CorsAllowHeaders`.
`HttpOpenTelemetryMiddleware` starts a server span for matched OpenAPI routes.

`HttpTripperwareStack` chains `http.RoundTripper` values. `HttpErrorMapperTripperware` and
`HttpToBricksErrorMapper` translate client responses into `bricks` errors.

## gRPC, NATS, and Asynq

Each transport has a middleware stack, a panic recoverer, and an OpenTelemetry middleware.
Asynq, NATS, and the generic message stack can also run a `bricks.Compensator` carried on
the handler error. `AsynqErrorDecoratorMiddleware` keeps `bricks.ErrRetryable` errors
retryable and marks every other error with `asynq.SkipRetry`.

## Errors

`HttpStatusToBricksError`, `BricksErrorToHttpStatusMapper`, `BricksToGrpcErrorMapper`,
`GrpcToBricksErrorMapper`, and `RedisToBricksError` convert between the transports and
the coded errors in `bricks`.

[go-pkg]: https://pkg.go.dev/github.com/janstoon/toolbox/handywares
[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-cov-handywares]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_handywares_refs_heads_master.json
