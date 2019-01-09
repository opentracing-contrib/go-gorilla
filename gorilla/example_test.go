package gorilla_test

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/uber/jaeger-client-go/config"
)

// http listen port
const httpAddress = ":8800"

func ExampleGorilla_TracingSingleRoute() {
	tracer, _, err := getTracer()

	if err != nil {
		log.Fatal("cannot initialize Jaeger Tracer", err)
	}

	myHandler := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["productId"]
		//do something
		data := "Hello, we get " + idstr
		fmt.Fprintf(w, data)
	}

	r := mux.NewRouter()

	pattern := "/v1/products/{productId}"

	middleware := gorilla.Middleware(
		tracer,
		http.HandlerFunc(myHandler),
	)

	r.Handle(pattern, middleware)
	log.Fatal(http.ListenAndServe(httpAddress, r))
}

func ExampleGorilla_TracingAllRoutes() {
	tracer, _, err := getTracer()

	if err != nil {
		log.Fatal("cannot initialize Jaeger Tracer", err)
	}

	okHandler := func(w http.ResponseWriter, r *http.Request) {
		// do something
		data := "Hello"
		fmt.Fprintf(w, data)
	}

	r := mux.NewRouter()
	// Create multiples routes
	r.HandleFunc("/v1/products", okHandler)
	r.HandleFunc("/v1/products/{productId}", okHandler)
	r.HandleFunc("/v2/products", okHandler)
	r.HandleFunc("/v2/products/{productId}", okHandler)
	r.HandleFunc("/v3/products", okHandler)
	r.HandleFunc("/v3/products/{productId}", okHandler)
	r.HandleFunc("/v4/products", okHandler)
	r.HandleFunc("/v4/products/{productId}", okHandler)

	// Add tracing to all routes
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		route.Handler(
			gorilla.Middleware(tracer, route.GetHandler()))
		return nil
	})

	log.Fatal(http.ListenAndServe(httpAddress, r))
}

func getTracer() (opentracing.Tracer, io.Closer, error) {
	//jaeger agent port
	jaegerHostPort := ":6831"

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
	return cfg.New(
		"ExampleTracingMiddleware", //service name
	)
}
