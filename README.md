# echosentry
[![GoDoc](https://godoc.org/github.com/01walid/echosentry?status.svg)](https://godoc.org/github.com/01walid/echosentry)

A sentry ([raven-go](https://github.com/getsentry/raven-go)) middleware for [echo](https://github.com/labstack/echo) micro web framework.

# Usage

```go
echosentry.SetDSN("https://<key>:<secret>@app.getsentry.com/<project>")
e.Use(echosentry.Middleware())

```

By default, the middleware logs the HTTP context and sends it along with the stacktrace, this adds info about the user's browser, URL, OS, device, interface_type ..etc.

You can disable HTTP context as follow:

```go
echosentry.WithContext(false)
```

## Additional tags

You can append additional tags to be captured by Sentry. Tags content can be extracted from the current request context or just static tags, e.g. tags["app_version"] = appVersion.

```go
echosentry.SetTags(func(c *echo.Context) map[string]string {
        return map[string]string{
            "endpoint":       c.Request().URL.String(),
            "http_interface": c.Request().Proto,
            "app_version":    appVersion,
        }
    })
```

# TODO
- Log the user info (user context), currently raven-go has an issue with that...
- expose more options

# License
MIT License. A copy is included with the source.
