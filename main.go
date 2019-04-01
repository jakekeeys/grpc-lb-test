package main

import (
	"context"
	"github.com/gorilla/handlers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jakekeeys/grpc-lb-test/sample"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/utilitywarehouse/swagger-ui/swaggerui"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:   "server",
			Action: server,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "grpc-bind",
					EnvVar: "GRPC_BIND",
					Value:  ":8090",
				},
				cli.StringFlag{
					Name:   "http-bind",
					EnvVar: "HTTP_BIND",
					Value:  ":8080",
				},
			},
		},
		{
			Name:   "client",
			Action: client,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "grpc-server",
					EnvVar: "GRPC_SERVER",
					Value:  ":8090",
				},
				cli.StringFlag{
					Name:   "http-bind",
					EnvVar: "HTTP_BIND",
					Value:  ":8081",
				},
				cli.IntFlag{
					Name:   "interval-ms",
					EnvVar: "INTERVAL_MS",
					Value:  100,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func client(c *cli.Context) {
	logrus.SetLevel(logrus.DebugLevel)

	conn, err := grpc.Dial(
		c.String("grpc-server"),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithUnaryInterceptor(grpc_logrus.UnaryClientInterceptor(logrus.NewEntry(logrus.StandardLogger()))),
	)
	if err != nil {
		panic(err)
	}
	client := sample.NewSampleServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(c.Int("interval-ms"))):
				_, err = client.SampleRPC(context.Background(), &sample.SampleRequest{})
				if err != nil {
					panic(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := http.Server{
		Addr:    c.String("http-bind"),
		Handler: handlers.LoggingHandler(logrus.StandardLogger().Out, mux),
	}
	defer httpServer.Shutdown(context.Background())

	go func() {
		err = httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}

func server(c *cli.Context) {
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_prometheus.UnaryServerInterceptor,
				grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
			),
		),
	)
	sampleServer := sample.Service{}
	sample.RegisterSampleServiceServer(grpcServer, &sampleServer)
	defer grpcServer.GracefulStop()

	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	rMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = sample.RegisterSampleServiceHandlerFromEndpoint(context.Background(), rMux, c.String("grpc-bind"), opts)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", rMux)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/swagger-ui/", swaggerui.UIHandler())
	mux.Handle("/swagger.json", swaggerui.FileHandler())

	httpServer := http.Server{
		Addr:    c.String("http-bind"),
		Handler: handlers.LoggingHandler(logrus.StandardLogger().Out, mux),
	}
	defer httpServer.Shutdown(context.Background())

	go func() {
		err = httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
