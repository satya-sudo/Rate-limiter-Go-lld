# Rate Limiter Gateway (Go)

A simple **API gateway–style rate limiter** implemented in Go using the **Token Bucket algorithm**.

The service sits in front of an API endpoint and controls how many requests a user can make over time.

---

## What is this?

This project is an in-memory **per-user rate limiter** built as HTTP middleware.

- Each user is identified by the `X-User-Id` header  
- Each user gets a token bucket  
- Requests are allowed only if a token is available  
- When tokens run out, the request is rejected with `429 Too Many Requests`

The rate limiter is concurrency-safe and runs entirely in memory.

---

## How does it work?

- Each user has a **token bucket**
- The bucket starts full (`capacity`)
- Tokens refill over time at a fixed rate (`refillRate`)
- Each request consumes **1 token**
- If no token is available, the request is blocked

Tokens refill lazily when requests arrive — no background refill loop.

---

## Current Rate Limit

```

Capacity:   5 tokens
RefillRate: 0.5 tokens per second

````

This allows:
- 5 immediate requests
- 1 new request every 2 seconds after that

---

## How is it implemented?

- **Token Bucket algorithm**
- **Per-bucket mutex** for thread safety
- **sync.Map** for storing user buckets
- **HTTP middleware** for enforcement
- **Background cleanup** removes idle buckets
- **Graceful shutdown** using a stop signal

---

## How to run

```bash
go run cmd/gateway/main.go
````

The server listens on:

```
http://localhost:8080
```

---

## Test

```bash
curl -H "X-User-Id: test-user" localhost:8080/api
```

If the rate limit is exceeded:

```
429 Too Many Requests
```

---
