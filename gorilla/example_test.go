package gorilla_test

import (
	"net/http"
	"log"
	"fmt"
	"time"

    "github.com/gorilla/mux"
    "github.com/opentracing-contrib/go-gorilla/gorilla"
    "github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-client-go/config"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
)

func NewServeMux(tracer opentracing.Tracer) *TracedServeMux {
    return &TracedServeMux{
        mr:   mux.NewRouter(),
        tracer: tracer,
    }
}

type TracedServeMux struct {
    mr    *mux.Router
    tracer opentracing.Tracer
}

func (tm *TracedServeMux) Handle(pattern string, handler http.Handler) {
    middleware := gorilla.Middleware(
        tm.tracer,
        handler,
        gorilla.OperationNameFunc(func(r *http.Request) string {
            return "Gorilla HTTP " + r.Method + " " + pattern
        }))
    tm.mr.Handle(pattern,middleware)
}

func (tm *TracedServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tm.mr.ServeHTTP(w, r)
}

func ExampleMiddleware() {
	httpAddress := ":8080"
	jaegerHostPort := ":6381"
	MyHandler := func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idstr := vars["productId"]
        //do something
		data := "Hello, we get " + idstr
        fmt.Fprintf(w, data)
    }
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jaegerHostPort,
		},
	}
	metricsFactory := jprom.New()
	tracer, _, err:= cfg.New(
		"ExampleMiddleware",
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)

	if err != nil {
		log.Fatal("cannot initialize Jaeger Tracer",err)
	}

    r := NewServeMux(tracer)
    r.Handle("/v1/products/{productId}", http.HandlerFunc(MyHandler))
    log.Fatal(http.ListenAndServe(httpAddress, r))
}

