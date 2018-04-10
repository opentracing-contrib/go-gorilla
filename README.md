# go-gorilla
OpenTracing instrumentation for Gorilla framework (github.com/gorilla)

# Usage

Here is example for this repo.
```
import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
)

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
```

Then, call it:
```
func main() {
	address := ":8080"
	MyHandler = func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
	    idstr := vars["productId"]
        //do something
		fmt.Fprintf(w, string(data))
	}
	tracer := make tracer here
	r := NewServeMux(tracer)
	r.Handle("/v1/product/{productId}", http.HandlerFunc(MyHandler))
	log.Fatal(http.ListenAndServe(address, r))
}
```

