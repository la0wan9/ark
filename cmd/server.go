package cmd

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/la0wan9/ark/api"
	"github.com/la0wan9/ark/internal/adoc"
	adocv1 "github.com/la0wan9/ark/pkg/adoc/v1"
	"github.com/la0wan9/ark/web"
)

// NewServerCmd creates a new server command
func NewServerCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use: "server",
		Run: func(command *cobra.Command, args []string) {
			exit := make(chan bool)
			fns := []func() error{
				startGrpcServer,
				startRestServer,
				startDebugServer,
			}
			for _, fn := range fns {
				go func(fn func() error) {
					defer func() { exit <- true }()
					if err := fn(); err != nil {
						log.Fatal(err)
					}
				}(fn)
			}
			<-exit
		},
	}
	serverCmd.Flags().SortFlags = false
	serverCmd.Flags().String("grpc-host", viper.GetString("grpc.host"), "grpc host")
	serverCmd.Flags().Int("grpc-port", viper.GetInt("grpc.port"), "grpc port")
	serverCmd.Flags().String("rest-host", viper.GetString("rest.host"), "rest host")
	serverCmd.Flags().Int("rest-port", viper.GetInt("rest.port"), "rest port")
	serverCmd.Flags().String("debug-host", viper.GetString("debug.host"), "debug host")
	serverCmd.Flags().Int("debug-port", viper.GetInt("debug.port"), "debug port")
	_ = viper.BindPFlag("grpc.host", serverCmd.Flags().Lookup("grpc-host"))
	_ = viper.BindPFlag("grpc.port", serverCmd.Flags().Lookup("grpc-port"))
	_ = viper.BindPFlag("rest.host", serverCmd.Flags().Lookup("rest-host"))
	_ = viper.BindPFlag("rest.port", serverCmd.Flags().Lookup("rest-port"))
	_ = viper.BindPFlag("debug.host", serverCmd.Flags().Lookup("debug-host"))
	_ = viper.BindPFlag("debug.port", serverCmd.Flags().Lookup("debug-port"))
	return serverCmd
}

func startGrpcServer() error {
	host := viper.GetString("grpc.host")
	port := viper.GetString("grpc.port")
	address := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	logger := log.New()
	entry := log.NewEntry(logger)
	grpc_logrus.ReplaceGrpcLogger(entry)
	decider := func(ctx context.Context, method string, server interface{}) bool {
		return true
	}
	server := grpc.NewServer(
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(entry),
			grpc_logrus.PayloadStreamServerInterceptor(entry, decider),
			grpc_prometheus.StreamServerInterceptor,
		),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(entry),
			grpc_logrus.PayloadUnaryServerInterceptor(entry, decider),
			grpc_prometheus.UnaryServerInterceptor,
		),
	)
	adocServer := &adoc.Server{}
	adocServer.Register(server)
	reflection.Register(server)
	grpc_prometheus.Register(server)
	return server.Serve(listener)
}

func startRestServer() error {
	host := viper.GetString("rest.host")
	port := viper.GetString("rest.port")
	address := net.JoinHostPort(host, port)
	host = viper.GetString("grpc.host")
	port = viper.GetString("grpc.port")
	endpoint := net.JoinHostPort(host, port)
	runtimeMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := adocv1.RegisterAdocServiceHandlerFromEndpoint(
		context.Background(), runtimeMux, endpoint, opts,
	); err != nil {
		return err
	}
	httpMux := http.NewServeMux()
	httpMux.Handle("/", runtimeMux)
	httpMux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.FS(web.FS))))
	httpMux.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.FS(api.FS))))
	return http.ListenAndServe(address, httpMux)
}

func startDebugServer() error {
	host := viper.GetString("debug.host")
	port := viper.GetString("debug.port")
	address := net.JoinHostPort(host, port)
	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.Handle("/debug/fgprof", fgprof.Handler())
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return http.ListenAndServe(address, httpMux)
}
