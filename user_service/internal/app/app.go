package app

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"buf.build/go/protovalidate"
	authv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/auth/v1"
	userv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/user/v1"
	"github.com/escoutdoor/vegetable_store/common/pkg/database"
	"github.com/escoutdoor/vegetable_store/common/pkg/database/pg"
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	common_interceptor "github.com/escoutdoor/vegetable_store/common/pkg/interceptor"
	"github.com/escoutdoor/vegetable_store/common/pkg/logger"
	authv1_implementation "github.com/escoutdoor/vegetable_store/user_service/internal/api/auth/v1"
	userv1_implementation "github.com/escoutdoor/vegetable_store/user_service/internal/api/user/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/interceptor"
	"github.com/escoutdoor/vegetable_store/user_service/internal/repository"
	user_repository "github.com/escoutdoor/vegetable_store/user_service/internal/repository/user"
	"github.com/escoutdoor/vegetable_store/user_service/internal/service"
	auth_service "github.com/escoutdoor/vegetable_store/user_service/internal/service/auth"
	user_service "github.com/escoutdoor/vegetable_store/user_service/internal/service/user"
	"github.com/escoutdoor/vegetable_store/user_service/internal/utils/token"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config        *Config
	grpcServer    *grpc.Server
	gatewayServer *http.Server

	dbClient       database.Client
	tokenProvider  token.Provider
	userService    service.UserService
	authService    service.AuthService
	userRepository repository.UserRepository
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
		a.initDBClient,
		a.initRepositories,
		a.initTokenProvider,
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

func (a *App) initDBClient(ctx context.Context) error {
	client, err := pg.NewClient(ctx, a.config.Postgres.Dsn)
	if err != nil {
		return errwrap.Wrap("pg new client", err)
	}

	if err := client.DB().Ping(ctx); err != nil {
		return errwrap.Wrap("database ping", err)
	}

	a.dbClient = client
	return nil
}

func (a *App) initRepositories(ctx context.Context) error {
	a.userRepository = user_repository.NewRepository(a.dbClient)
	return nil
}

func (a *App) initTokenProvider(_ context.Context) error {
	a.tokenProvider = token.NewTokenProvider(
		a.config.Token.AccessTokenSecretKey,
		a.config.Token.RefreshTokenSecretKey,
		a.config.Token.AccessTokenTTL,
		a.config.Token.RefreshTokenTTL,
	)
	return nil
}

func (a *App) initServices(_ context.Context) error {
	a.authService = auth_service.NewService(a.userRepository, a.tokenProvider)
	a.userService = user_service.NewService(a.userRepository)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	validator, err := protovalidate.New()
	if err != nil {
		return errwrap.Wrap("new validator", err)
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.ErrorsUnaryServerInterceptor(),
		common_interceptor.ValidationUnaryServerInterceptor(validator),
		common_interceptor.LoggingUnaryServerInterceptor(),
	))

	authImpl := authv1_implementation.NewImplementation(a.authService)
	userImpl := userv1_implementation.NewImplementation(a.userService)

	authv1.RegisterAuthServiceServer(grpcServer, authImpl)
	userv1.RegisterUserServiceServer(grpcServer, userImpl)

	reflection.Register(grpcServer)

	a.grpcServer = grpcServer
	return nil
}

func (a *App) initGatewayServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, a.config.GRPC.Address(), opts); err != nil {
		return errwrap.Wrap("register user service handler from endpoint", err)
	}

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, a.config.GRPC.Address(), opts); err != nil {
		return errwrap.Wrap("register auth service handler from endpoint", err)
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
		return errwrap.Wrap("gateway http server listen and serve", err)
	}

	return nil
}
