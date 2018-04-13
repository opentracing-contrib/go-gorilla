package gorilla_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/uber/jaeger-client-go/config"
)

func ExampleTracingMiddleware() {
	// http listen port
	httpAddress := ":8800"

	//jaeger agent port
	jaegerHostPort := ":6831"
	myHandler := func(w http.ResponseWriter, r *http.Request) {
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
	tracer, _, err := cfg.New(
		"ExampleMiddleware", //operationName
	)

	if err != nil {
		log.Fatal("cannot initialize Jaeger Tracer", err)
	}

	r := mux.NewRouter()

	pattern := "/v1/products/{productId}"
	middleware := gorilla.NewTracingMiddleware(tracer)

	r.HandleFunc(pattern, myHandler)
	r.Use(middleware.With)
	log.Fatal(http.ListenAndServe(httpAddress, r))
}
