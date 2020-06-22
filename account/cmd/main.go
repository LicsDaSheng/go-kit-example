package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"net/http"
	"oncekey/go-kit-example/account/dao"
	"oncekey/go-kit-example/account/endpoints"
	"oncekey/go-kit-example/account/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/lib/pq"
	opentracinggo "github.com/opentracing/opentracing-go"
	zipkingoopentracing "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkingo "github.com/openzipkin/zipkin-go"
	zphttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const dbsource = "postgresql://postgres:123456@localhost:5432/gokitexample?sslmode=disable"

var tracer opentracinggo.Tracer

func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service", "account",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// zipkin
	var zpServerTrace httptransport.ServerOption
	{
		reporter := zphttp.NewReporter("http://localhost:9411/api/v2/spans")
		defer reporter.Close()
		endpoint, err := zipkingo.NewEndpoint("Account", "localhost:8080")
		if err != nil {
			level.Error(logger).Log("err", err)
			os.Exit(1)
		}
		localEndpoint := zipkingo.WithLocalEndpoint(endpoint)
		nativeTracer, err := zipkingo.NewTracer(reporter, localEndpoint)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		tracer = zipkingoopentracing.Wrap(nativeTracer)
		zpTracer, _ := zipkingo.NewTracer(reporter)
		zpServerTrace = zipkin.HTTPServerTrace(zpTracer)
		zpServerTrace = httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "Create", logger))
	}

	// db
	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres", dbsource)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}
	flag.Parse()

	ctx := context.Background()
	var srv service.Service
	{
		repository := dao.NewDao(db, logger)
		srv = service.NewService(repository, logger)
	}
	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	eds := endpoints.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := endpoints.NewHTTPServer(ctx, eds, zpServerTrace)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
