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
//      nethttp.OperationNameFunc(func(r *http.Request) string {
//	        return r.Proto + " " + r.Method + ":" + pattern
//      }),
//      nethttp.MWSpanObserver(func(sp opentracing.Span, r *http.Request) {
//			sp.SetTag("http.uri", r.URL.EscapedPath())
//		}),
//   )
//   r := mux.NewRouter()
//   r.HandleFunc(pattern, handler)
//   r.Use(middleware.With)
func Middleware(tr opentracing.Tracer, h http.Handler, options ...nethttp.MWOption) http.Handler {
	opNameFunc := func(r *http.Request) string {
		if tpl, err := mux.CurrentRoute(r).GetPathTemplate(); nil == err {
			return r.Proto + " " + r.Method + " " + tpl
		}
		return r.Proto + " " + r.Method + " " + r.URL.Path 
	}
	var opts []nethttp.MWOption
	opts = append(opts,nethttp.OperationNameFunc(opNameFunc))
	opts = append(opts, options...)
	return nethttp.Middleware(tr,h,opts...)
}
