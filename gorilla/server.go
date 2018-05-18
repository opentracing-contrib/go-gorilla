// +build go1.7

package gorilla

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
)

const defaultComponentName = "github.com/gorilla/mux"

// Middleware wraps an http.Handler and traces incoming requests.
// Additionally, it adds the span to the request's context.
//
// By default, the operation name of the spans is set to "HTTP {method}".
// This can be overriden with options.
//
// The options allow fine tuning the behavior of the middleware.
//
// Example:
//   pattern := "/api/custerms/{id}"
//   mw := gorilla.Middleware(
//      tracer,
//      handler,
//   )
//   r := mux.NewRouter()
//   r.HandleFunc(pattern, mw)
func Middleware(tr opentracing.Tracer, h http.Handler, options ...nethttp.MWOption) http.Handler {
	opNameFunc := func(r *http.Request) string {
		if route := mux.CurrentRoute(r); route != nil {
			if tpl, err := route.GetPathTemplate(); err == nil {
				return r.Proto + " " + r.Method + " " + tpl
			}
		}
		return r.Proto + " " + r.Method
	}
	var opts = []nethttp.MWOption{nethttp.OperationNameFunc(opNameFunc), nethttp.MWComponentName(defaultComponentName)}
	opts = append(opts, options...)
	return nethttp.Middleware(tr, h, opts...)
}
