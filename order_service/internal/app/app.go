package app

import (
	"context"
	"net"
	"net/http"
	"sync"
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
	order_implementation "github.com/escoutdoor/vegetable_store/order_service/internal/api/order/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/client"
	"github.com/escoutdoor/vegetable_store/order_service/internal/client/vegetable"
	"github.com/escoutdoor/vegetable_store/order_service/internal/interceptor"
	"github.com/escoutdoor/vegetable_store/order_service/internal/repository"
	order_repository "github.com/escoutdoor/vegetable_store/order_service/internal/repository/order"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service"
	order_service "github.com/escoutdoor/vegetable_store/order_service/internal/service/order"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config        *Config
	grpcServer    *grpc.Server
	gatewayServer *http.Server

	vegetableGrpcClient vegetablev1.VegetableServiceClient
	vegetableClient     client.VegetableClient

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
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		logger.Info(ctx, "grpc server is running")
		defer wg.Done()
		if err := a.runGRPCServer(); err != nil {
			logger.Fatal(ctx, "run grpc server", err)
		}
	}()

	go func() {
		logger.Info(ctx, "grpc gateway server is running")
		defer wg.Done()
		if err := a.runGatewayServer(); err != nil {
			logger.Fatal(ctx, "run grpc gateway server server", err)
		}
	}()

	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	deps := []func(ctx context.Context) error{
		a.initVegetableGrpcClient,
		a.initVegetableClient,
		a.initDBClient,
		a.initTransactionManager,
		a.initRepositories,
		a.initServices,
		a.initGRPCServer,
		a.initGatewayServer,
	}

	for _, d := range deps {
		if err := d(ctx); err != nil {
			return err
		}
	}

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

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		common_interceptor.ValidationUnaryServerInterceptor(validator),
		interceptor.ErrorsUnaryServerInterceptor(),
		common_interceptor.LoggingUnaryServerInterceptor(),
	))

	orderImpl := order_implementation.NewImplementation(a.orderService)
	orderv1.RegisterOrderServiceServer(grpcServer, orderImpl)

	reflection.Register(grpcServer)

	a.grpcServer = grpcServer
	return nil
}

func (a *App) initGatewayServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := orderv1.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, a.config.GRPC.Address(), opts); err != nil {
		return errwrap.Wrap("register order service handler from endpoint", err)
	}
	httpServer := &http.Server{
		Addr:              a.config.Gateway.Address(),
		Handler:           mux,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
	}

	a.gatewayServer = httpServer
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
