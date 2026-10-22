# api

[![Go Reference](https://pkg.go.dev/badge/go.trueardian.com/api.svg)](https://pkg.go.dev/go.trueardian.com/api)
![Go 1.25+](https://img.shields.io/badge/go-1.25%2B-00ADD8)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A Go library that carries who's calling from HTTP headers into `context.Context`, and back out again.

## Install

```sh
go get go.trueardian.com/api
```

## The packages

| Package       | What it is                                                              | Provides                 | Depends on |
| ------------- | ----------------------------------------------------------------------- | ------------------------ | ---------- |
| `api/user`    | Reads the `user-id` header into the context. `401` when it's missing.    | `user.HTTPWithID`        | `api/http` |
| `api/session` | Reads the `session-id` header into the context. `401` when it's missing. | `session.HTTPWithID`     | `api/http` |
| `api/time`    | Reads the `time-zone` header into the context. `400` when it's invalid.  | `time.HTTPWithZone`      | `api/http` |
| `api/http`    | Writes context values back as headers on outbound requests, and JSON responses. | `*http.Transport` |            |

```
user ────┐
session ─┼─▶ http
time ────┘
```

The whole module is standard library only, so `go get` pulls in nothing else. Only import what you
need. Below are the explanations and guides.

## `package user`

So your service sits behind a gateway that already knows who the caller is, and passes it along as
a header? Then every handler wants that user ID, and you don't want to thread it through every
function signature to get it there. I got you.

Wrap your mux:

```go
handler := user.HTTPWithID(mux)
```

It reads the `user-id` header and puts it into the context. Missing or empty, and the request stops
right there with a `401` and `{"detail":"user ID cannot be empty"}`, so your handler never runs
without one. Then, anywhere below it:

```go
uid, err := user.IDFromContext(r.Context())
```

Behind the middleware that error never fires. It's there for code reached some other way, like a
cron job or a test, where nobody filled the context. Not coming from HTTP at all? Fill it yourself
with `user.ContextWithID(ctx, id)`, which runs the same empty check without the header.

The header is `user-id`, not `X-User-Id`. The `X-` prefix was retired by
[RFC 6648](https://www.rfc-editor.org/rfc/rfc6648) back in 2012, and HTTP/2 lowercases every header
name on the wire anyway, so this is just the honest name. Go matches header names
case-insensitively, so a caller sending `User-Id` works too.

## `package session`

Exactly the same as `user`, for the `session-id` header:

```go
handler := session.HTTPWithID(mux)

sid, err := session.IDFromContext(r.Context())
```

Missing or empty is a `401` here too, and `session.ContextWithID(ctx, id)` is there for when you
fill it yourself.

## `package time`

So you want to show times in the caller's own timezone, but the server runs in UTC? Have them send
an IANA name like `Asia/Jakarta` in the `time-zone` header:

```go
import apitime "go.trueardian.com/api/time"

handler := apitime.HTTPWithZone(mux)

tz, err := apitime.ZoneFromContext(r.Context())
```

The name is checked with `time.LoadLocation` before it goes into the context, so `Mars/Olympus`
never reaches your handler. That one is a `400`, not a `401`, because a bad timezone is a malformed
request, not an unknown caller. What you get back is the name itself, so `time.LoadLocation(tz)` it
when you need the `*time.Location`.

Do note the package is named `time`, same as the standard library, so alias it on import.

Stack all three in any order. Each one stops the request on its own:

```go
handler := user.HTTPWithID(session.HTTPWithID(apitime.HTTPWithZone(mux)))
```

## `package http`

So the context is filled, and now your handler calls the next service, which wants those same
headers? I got you.

```go
import apihttp "go.trueardian.com/api/http"

client := &http.Client{
    Transport: &apihttp.Transport{
        Propagators: []apihttp.Propagator{
            apihttp.WithHeader(user.ContextKey, "user-id"),
            apihttp.WithHeader(session.ContextKey, "session-id"),
            apihttp.WithHeader(apitime.ContextKey, "time-zone"),
        },
    },
}

req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://downstream.example.com/internal", nil)
resp, err := client.Do(req)
```

Every request through that client gets the values from its context written back as headers. A value
that isn't there is skipped, so the same client works for anonymous calls too. `Base` takes your own
`RoundTripper`, and leaving it nil means `http.DefaultTransport`. The request is cloned before any
header is set, because a `RoundTripper` must never modify the request it was handed.

`user.ContextKey` is exported so `WithHeader` can name it from another package, but its type isn't,
so no other package can make a key that collides with it. `WithHeader` only copies a non-empty
string. Want anything else on the wire? A `Propagator` is just `func(ctx, header)`, so write your
own.

The package also has the two writers the middleware answers with:

```go
apihttp.WriteJSON(w, http.StatusOK, payload)
apihttp.WriteProblem(w, http.StatusUnprocessableEntity, map[string]any{"detail": "..."})
```

They are identical today, and that's on purpose. `WriteProblem` is the one that grows into
[RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) problem details, so every error call site is
already pointing at the right function when it does. A body that fails to encode gets logged,
because the status is already sent by then.

Same as `time`, the package name collides with the standard library, so alias it.

## What this library doesn't do

- **Authenticate anyone.** It carries an identity, it doesn't verify one. Anyone who can reach your
  service can send whatever `user-id` they like, so these headers have to come from something you
  trust at the edge, one that strips them off outside traffic.
- **Route.** Every middleware is a plain `func(http.Handler) http.Handler`. Bring your own mux.
- **Logging, tracing or metrics.** A `Propagator` can carry a trace ID just fine, but none ships.
- **Carry more than these three.** Need a request ID, a locale, a tenant? Write your own, buddy:
  copy `session`, rename it, and you're done. It's two files. Pull requests are welcome too, though
  I'm fairly strict about how code is written here, so expect me to refactor yours before merging.

## Roadmap

- [ ] Full RFC 9457 problem details behind `WriteProblem`, served as `application/problem+json`.

## License

[MIT](LICENSE) © 2026 Ardian
