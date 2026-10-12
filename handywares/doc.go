// Package handywares provides reusable middleware for HTTP servers and clients, gRPC
// unary calls, NATS messages, Asynq tasks, and generic message handlers.
//
// Stacks are tricks.Middleware values composed in FIFO order. Error helpers translate
// between bricks errors and HTTP statuses, gRPC status codes, and Redis client errors.
package handywares
