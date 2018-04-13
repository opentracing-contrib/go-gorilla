package gorilla

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestOperationNameOption(t *testing.T) {
	mr := mux.NewRouter()
	mr.HandleFunc("/root", func(w http.ResponseWriter, r *http.Request) {})

	fn := func(r *http.Request) string {
		return "HTTP/1.1 " + r.Method + " /root"
	}

	tests := []struct {
		options []MWOption
		opName  string
	}{
		{nil, "HTTP/1.1 GET /root"},
		{[]MWOption{OperationNameFunc(fn)}, "HTTP/1.1 GET /root"},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.opName, func(t *testing.T) {
			tr := &mocktracer.MockTracer{}
			mw := NewTracingMiddleware(tr, testCase.options...)
			mr.Use(mw.With)
			srv := httptest.NewServer(mr)
			defer srv.Close()

			_, err := http.Get(srv.URL + "/root")
			if err != nil {
				t.Fatalf("server returned error: %v", err)
			}

			spans := tr.FinishedSpans()
			if got, want := len(spans), 1; got != want {
				t.Fatalf("got %d spans, expected %d", got, want)
			}

			if got, want := spans[0].OperationName, testCase.opName; got != want {
				t.Fatalf("got %s operation name, expected %s", got, want)
			}
		})
	}
}

func TestSpanObserverOption(t *testing.T) {
	mr := mux.NewRouter()
	mr.HandleFunc("/root", func(w http.ResponseWriter, r *http.Request) {})

	opNamefn := func(r *http.Request) string {
		return "HTTP/1.1 " + r.Method + " /root"
	}
	spanObserverfn := func(sp opentracing.Span, r *http.Request) {
		sp.SetTag("http.uri", r.URL.EscapedPath())
	}
	wantTags := map[string]interface{}{"http.uri": "/root"}

	tests := []struct {
		options []MWOption
		opName  string
		Tags    map[string]interface{}
	}{
		{nil, "HTTP/1.1 GET /root", nil},
		{[]MWOption{OperationNameFunc(opNamefn)}, "HTTP/1.1 GET /root", nil},
		{[]MWOption{MWSpanObserver(spanObserverfn)}, "HTTP/1.1 GET /root", wantTags},
		{[]MWOption{OperationNameFunc(opNamefn), MWSpanObserver(spanObserverfn)}, "HTTP/1.1 GET /root", wantTags},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.opName, func(t *testing.T) {
			tr := &mocktracer.MockTracer{}
			mw := NewTracingMiddleware(tr, testCase.options...)
			mr.Use(mw.With)
			srv := httptest.NewServer(mr)
			defer srv.Close()

			_, err := http.Get(srv.URL + "/root")
			if err != nil {
				t.Fatalf("server returned error: %v", err)
			}

			spans := tr.FinishedSpans()
			if got, want := len(spans), 1; got != want {
				t.Fatalf("got %d spans, expected %d", got, want)
			}

			if got, want := spans[0].OperationName, testCase.opName; got != want {
				t.Fatalf("got %s operation name, expected %s", got, want)
			}

			defaultLength := 5
			if len(spans[0].Tags()) != len(testCase.Tags)+defaultLength {
				t.Fatalf("got tag length %d, expected %d", len(spans[0].Tags()), len(testCase.Tags))
			}
			for k, v := range testCase.Tags {
				if tag := spans[0].Tag(k); v != tag.(string) {
					t.Fatalf("got %v tag, expected %v", tag, v)
				}
			}
		})
	}
}
