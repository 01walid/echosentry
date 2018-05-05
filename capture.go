package echosentry

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/labstack/echo"
	"github.com/getsentry/raven-go"
)

// Sentry struct holding the raven client and some of its configs
type Sentry struct {
	withContext bool
	RavenClient *raven.Client
	Tags        map[string]string
}

// TagsFunc given a request context, extract some additional tags and return
// them as map[string]string as required by the raven client.
type TagsFunc func(c echo.Context) map[string]string

var (
	sentry   = &Sentry{}
	tagsFunc TagsFunc
)

// SetDSN creates a raven client and sets its Sentry server DSN.
func SetDSN(dsn string) {
	client, err := raven.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	sentry.RavenClient = client
}

// WithContext sets weather or not the HTTP context is sent with the log.
// This adds info about the user's browser, URL, OS, device, interface_type ..etc
func WithContext(yepnope bool) {
	sentry.withContext = yepnope
}

// Sets any other additional tags to be captured by Sentry.
// Tags can be extracted from the current request context
// or just static tags, e.g. tags["app_version"] = appVersion.
func SetTags(fn TagsFunc) {
	tagsFunc = fn
}

// Middleware returns an echo middleware which recovers from panics anywhere in the chain
// and logs to the sentry server specified in DSN.
func Middleware() echo.MiddlewareFunc {

	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if rval := recover(); rval != nil {
					debug.PrintStack()

					errorMsg := fmt.Sprint(rval)
					err := errors.New(errorMsg)

					stacktrace := raven.NewException(err, raven.NewStacktrace(2, 3, nil))

					httpContext := &raven.Http{}

					if sentry.withContext {
						httpContext = raven.NewHttp(c.Request())
					}

					// extract tags
					if tagsFunc != nil {
						sentry.Tags = tagsFunc(c)
					}

					// contruct the raven packet to be sent
					packet := raven.NewPacket(errorMsg, stacktrace, httpContext)

					// capture it and send.
					sentry.RavenClient.Capture(packet, sentry.Tags)

					// register the error back to echo.Context
					c.Error(err)
				}
			}()

			return h(c)
		}
	}
}

func init() {
	// HTTP context enabled by default for convenience
	sentry.withContext = true
}
