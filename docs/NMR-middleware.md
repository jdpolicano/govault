# Middleware Proposal for govault

## Motivation
The `get` and `set` routes both perform the same preliminary steps:

1. Validate the `Authorization` header and fetch the session.
2. Parse a JSON body into a Go struct.

As more routes are added this duplication will grow. Centralizing these tasks as reusable middleware avoids repeated code and simplifies updates to the request pipeline.

## Goals
- Remove duplicated session and body parsing logic from route handlers.
- Provide an easy way to compose future middleware (e.g. logging).

## Proposed Approach
Create a new package `internal/server/middleware` implementing a simple middleware chain pattern.

```go
// Middleware defines a function that wraps an http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies several middleware around a final handler.
func Chain(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
    for i := len(m) - 1; i >= 0; i-- {
        h = m[i](h)
    }
    return h
}
```

### Token Validation Middleware
```go
func ValidateToken(ctx *server.Context) Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            sess, err := ctx.ValidateTokenHeader(r)
            if err != nil {
                http.Error(w, err.Error(), http.StatusUnauthorized)
                return
            }
            // store the session on the request context for later use
            r = r.WithContext(context.WithValue(r.Context(), server.SessionKey{}, sess))
            next(w, r)
        }
    }
}
```
The middleware stores the validated `Session` in the request context for downstream handlers.

### Body Parsing Middleware
Generics allow parsing any target struct:
```go
func ParseJSONBody[T any]() Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            var body T
            dec := json.NewDecoder(r.Body)
            if err := dec.Decode(&body); err != nil {
                http.Error(w, errors.InvalidRequestBody.Error(), http.StatusBadRequest)
                return
            }
            r = r.WithContext(context.WithValue(r.Context(), server.BodyKey{}, body))
            next(w, r)
        }
    }
}
```
Handlers can retrieve the parsed value from context.

### Logging Middleware (Stretch Goal)
```go
func Logging(log *log.Logger) Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next(w, r)
            log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
        }
    }
}
```

## Example Usage
```
func Handler(ctx *server.Context) http.HandlerFunc {
    handle := func(w http.ResponseWriter, r *http.Request) {
        sess := r.Context().Value(server.SessionKey{}).(server.Session)
        body := r.Context().Value(server.BodyKey{}).(SetRequest)
        // existing set logic using sess and body
    }
    return middleware.Chain(handle,
        middleware.ValidateToken(ctx),
        middleware.ParseJSONBody[SetRequest](),
        middleware.Logging(ctx.Log),
    )
}
```

## Implementation Steps
1. Add `internal/server/middleware` with the `Middleware` type and helpers above.
2. Update existing routes to use `middleware.Chain` instead of performing validation and decoding directly.
3. Store/retrieve session and body via context keys defined in `server` package (e.g. `type SessionKey struct{}`).
4. Introduce logging middleware optionally applied to routes.
5. Unit test middleware functions.

## Future Considerations
- Additional middleware for rate limiting or error handling.
- Ability to configure middleware chains per route.
- Refactor older handlers once middleware proves stable.
