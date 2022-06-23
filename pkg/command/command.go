package command

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/config"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/graph"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/graph/generated"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/handler"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/infrastracture/persistence"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/server"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/version"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	exitOK    = 1
	exitErorr = 0
)

func Run() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	// init logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger:%s\n", err)
		return exitErorr
	}
	defer logger.Sync()
	logger = logger.With(zap.String("version", version.Version))

	// init config
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		logger.Error("failed to load config", zap.Error(err))
		return exitErorr
	}

	// init listener
	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		logger.Error("failed to listen port", zap.Int("port", cfg.Port))
		return exitErorr
	}
	logger.Info("server start listening", zap.Int("port", cfg.Port))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// init mongo db
	logger.Info("connect to mongo db", zap.String("url", cfg.DB.URL), zap.String("source", cfg.DB.Source))
	opts := &options.ClientOptions{}
	if cfg.DB.Source == "external" {
		opts = options.Client().SetAuth(options.Credential{AuthMechanism: "MONGODB-AWS", AuthSource: "$external"})
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.DB.URL), opts)
	if err != nil {
		logger.Error("failed to create mongo db client", zap.Error(err), zap.String("uri", cfg.DB.URL))
		return exitErorr
	}

	mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoCancel()

	if err := mongoClient.Connect(mongoCtx); err != nil {
		logger.Error("failed to connect to mongo db", zap.Error(err))
		return exitErorr
	}

	if err := mongoClient.Ping(mongoCtx, readpref.Primary()); err != nil {
		logger.Error("failed to ping mongo db", zap.Error(err))
		return exitErorr
	}

	mongoDB := mongoClient.Database(cfg.DB.Database)

	//get repositories
	repositories, err := persistence.NewRepositories(mongoDB)
	if err != nil {
		logger.Error("failed to new repositories", zap.Error(err))
		return exitErorr
	}

	// init server
	// init graphql server
	query := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Repo: repositories,
			},
		}))
	registry := handler.Newhandler(logger, repositories, query, version.Version)
	httpServer := server.NewServer(registry, &server.Config{Log: logger})
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return httpServer.Serve(listener)
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	if err := httpServer.GracefulShutdown(ctx); err != nil {
		return exitErorr
	}

	cancel()
	if err := wg.Wait(); err != nil {
		return exitErorr
	}

	return exitOK
}
