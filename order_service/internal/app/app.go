package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	orderv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/order/v1"
	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/common/pkg/database"
	"github.com/escoutdoor/vegetable_store/common/pkg/database/pg"
	"github.com/escoutdoor/vegetable_store/common/pkg/database/txmanager"
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	common_interceptor "github.com/escoutdoor/vegetable_store/common/pkg/interceptor"
	"github.com/escoutdoor/vegetable_store/common/pkg/logger"
	"github.com/escoutdoor/vegetable_store/common/pkg/tracing"
	order_implementation "github.com/escoutdoor/vegetable_store/order_service/internal/api/order/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/client"
	"github.com/escoutdoor/vegetable_store/order_service/internal/client/vegetable"
	"github.com/escoutdoor/vegetable_store/order_service/internal/interceptor"
	"github.com/escoutdoor/vegetable_store/order_service/internal/metrics"
	"github.com/escoutdoor/vegetable_store/order_service/internal/repository"
	order_repository "github.com/escoutdoor/vegetable_store/order_service/internal/repository/order"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service"
	order_service "github.com/escoutdoor/vegetable_store/order_service/internal/service/order"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config           *Config
	grpcServer       *grpc.Server
	gatewayServer    *http.Server
	prometheusServer *http.Server

	vegetableServiceClientConnection *grpc.ClientConn
	vegetableGrpcClient              vegetablev1.VegetableServiceClient
	vegetableClient                  client.VegetableClient

	dbClient           database.Client
	transactionManager database.TxManager
	orderService       service.OrderService
	orderRepository    repository.OrderRepository
}

func New(ctx context.Context, cfg *Config) (*App, error) {
	app := &App{config: cfg}

	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	quitch := make(chan os.Signal, 1)
	signal.Notify(quitch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info(ctx, "grpc server is running")
		if err := a.runGRPCServer(); err != nil {
			logger.Fatal(ctx, "run grpc server", err)
		}
	}()

	go func() {
		logger.Info(ctx, "grpc gateway server is running")
		if err := a.runGatewayServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(ctx, "run grpc gateway server", err)
		}
	}()

	go func() {
		logger.Info(ctx, "prometheus server is running")
		if err := a.runPrometheusServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(ctx, "run prometheus server", err)
		}
	}()

	<-quitch
	ctx, cancel := context.WithTimeout(ctx, a.config.GracefullShutdownTimeout)
	defer cancel()

	if err := a.stop(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) stop(ctx context.Context) error {
	logger.Info(ctx, "graceful shutdown started")
	defer logger.Info(ctx, "graceful shutdown completed")
	defer a.dbClient.Close()

	return multierr.Combine(
		a.gatewayServer.Shutdown(ctx),
		a.prometheusServer.Shutdown(ctx),
		a.gracefulShutdownGrpcServer(ctx),
		a.vegetableServiceClientConnection.Close(),
	)
}

func (a *App) initDeps(ctx context.Context) error {
	deps := []func(ctx context.Context) error{
		a.initTracing,
		a.initMetrics,
		a.initVegetableGrpcClient,
		a.initVegetableClient,
		a.initDBClient,
		a.initTransactionManager,
		a.initRepositories,
		a.initServices,
		a.initGRPCServer,
		a.initGatewayServer,
		a.initPrometheusServer,
	}

	for _, d := range deps {
		if err := d(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	if err := tracing.Init(ctx, a.config.Jaeger.Address(), a.config.AppName); err != nil {
		return errwrap.Wrap("init tracing", err)
	}

	return nil
}

func (a *App) initMetrics(_ context.Context) error {
	metrics.Init("vegetable_store", a.config.AppName)
	return nil
}

func (a *App) initVegetableGrpcClient(_ context.Context) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(a.config.VegetableClient.Address(), opts...)
	if err != nil {
		return errwrap.Wrap("create new vegetable service grpc connection", err)
	}
	a.vegetableServiceClientConnection = conn
	a.vegetableGrpcClient = vegetablev1.NewVegetableServiceClient(conn)

	return nil
}
func (a *App) initVegetableClient(_ context.Context) error {
	a.vegetableClient = vegetable.NewVegetableClient(a.vegetableGrpcClient)
	return nil
}

func (a *App) initDBClient(ctx context.Context) error {
	client, err := pg.NewClient(ctx, a.config.Postgres.Dsn)
	if err != nil {
		return errwrap.Wrap("new database client", err)
	}

	if err := client.DB().Ping(ctx); err != nil {
		return errwrap.Wrap("ping database", err)
	}

	a.dbClient = client
	return nil
}

func (a *App) initTransactionManager(_ context.Context) error {
	a.transactionManager = txmanager.NewTransactionManager(a.dbClient.DB())
	return nil
}

func (a *App) initRepositories(ctx context.Context) error {
	a.orderRepository = order_repository.NewRepository(a.dbClient)
	return nil
}

func (a *App) initServices(_ context.Context) error {
	a.orderService = order_service.NewService(a.orderRepository, a.transactionManager, a.vegetableClient)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	validator, err := protovalidate.New()
	if err != nil {
		return errwrap.Wrap("new validator", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.ErrorsUnaryServerInterceptor(),
			interceptor.MetricsUnaryServerInterceptor(),
			common_interceptor.LoggingUnaryServerInterceptor(),
			common_interceptor.ValidationUnaryServerInterceptor(validator),
			common_interceptor.RecoverUnaryServerInterceptor(),
		),
	)

	orderImpl := order_implementation.NewImplementation(a.orderService)
	orderv1.RegisterOrderServiceServer(grpcServer, orderImpl)

	reflection.Register(grpcServer)

	a.grpcServer = grpcServer
	return nil
}

func (a *App) initGatewayServer(ctx context.Context) error {
	gwMux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := orderv1.RegisterOrderServiceHandlerFromEndpoint(ctx, gwMux, a.config.GRPC.Address(), opts); err != nil {
		return errwrap.Wrap("register order service handler from endpoint", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	mux.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, a.config.Swagger.FilePath())
	})
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(a.config.Swagger.Path))))

	httpServer := &http.Server{
		Addr:              a.config.Gateway.Address(),
		Handler:           mux,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
	}

	a.gatewayServer = httpServer
	return nil
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:              a.config.Prometheus.Address(),
		Handler:           mux,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
	}
	return nil
}

func (a *App) runGRPCServer() error {
	ln, err := net.Listen("tcp", a.config.GRPC.Address())
	if err != nil {
		return errwrap.Wrap("net listen", err)
	}

	if err := a.grpcServer.Serve(ln); err != nil {
		return errwrap.Wrap("grpc server serve", err)
	}

	return nil
}

func (a *App) runGatewayServer() error {
	if err := a.gatewayServer.ListenAndServe(); err != nil {
		return errwrap.Wrap("gateway server listen and serve", err)
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	if err := a.prometheusServer.ListenAndServe(); err != nil {
		return errwrap.Wrap("prometheus server listen and serve", err)
	}

	return nil
}

func (a *App) gracefulShutdownGrpcServer(ctx context.Context) error {
	qch := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(qch)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-qch:
		return nil
	}
}
