# ddzerolog

## Description

Package ddzerolog provides a log/span correlation func for the [github.com/rs/zerolog](https://github.com/rs/zerolog) and [github.com/DataDog/dd-trace-go](https://github.com/DataDog/dd-trace-go).

## Example

```go
package ddzerolog_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/johejo/ddzerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Example() {
	l := zerolog.New(os.Stdout)

	ctx := context.Background()

	span, ctx := tracer.StartSpanFromContext(ctx, "testSpan")
	defer span.Finish()

	l.UpdateContext(ddzerolog.UpdateContext(ctx))

	ctx = l.WithContext(ctx)

	log.Ctx(ctx).Info().Msg("test")
	// Output: {"level":"info","dd.trace_id":0,"dd.span_id":0,"message":"test"}
}

func Example_nethttp() {
	mux := http.NewServeMux()
	log.Logger = zerolog.New(os.Stdout)

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msg("hello")
		w.Write([]byte("hello\n"))
	})

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			l := log.Logger
			l.UpdateContext(ddzerolog.UpdateContext(ctx))
			ctx = l.WithContext(ctx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	handler := httptrace.WrapHandler(middleware(mux), "myService", "myResource")

	s := httptest.NewServer(handler)
	defer s.Close()

	_, err := http.Get(s.URL + "/hello")
	if err != nil {
		panic(err)
	}

	// Output: {"level":"info","dd.trace_id":0,"dd.span_id":0,"message":"hello"}
}
```

## See Also

- https://docs.datadoghq.com/tracing/connect_logs_and_traces/
- https://docs.datadoghq.com/tracing/faq/why-cant-i-see-my-correlated-logs-in-the-trace-id-panel/?tabs=jsonlogs
