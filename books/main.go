package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const svcName = "books"

var tracer trace.Tracer

func main() {
	ctx := context.Background()

	exp, err := instrumentation.newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := instrumentation.newTraceProvider(exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer(svcName)

	r := mux.NewRouter()
	r.Use(otelmux.Middleware(svcName))

	r.HandleFunc("/books", booksListHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8001", nil))
}

func booksListHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "Books List")
	defer span.End()

	io.WriteString(w, "Hello!\n")
}
