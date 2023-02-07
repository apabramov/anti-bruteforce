package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/apabramov/anti-bruteforce/internal/app"
	"github.com/apabramov/anti-bruteforce/internal/bucket"
	cfg "github.com/apabramov/anti-bruteforce/internal/config"
	"github.com/apabramov/anti-bruteforce/internal/logger"
	internalredis "github.com/apabramov/anti-bruteforce/internal/redis"
	internalgrpc "github.com/apabramov/anti-bruteforce/internal/server/grpc"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/anti-bruteforce/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := cfg.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	cache, err := internalredis.New(config.Redis)
	if err != nil {
		logg.Log.Info(err.Error())
	}
	bucket := bucket.New(cache, config.Limit)

	storage, err := app.NewStorage(logg, config.Storage)
	if err != nil {
		logg.Log.Info(err.Error())
	}
	app := app.New(logg, storage, bucket)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("anti-bruteforce is running...")
	start(ctx, &config, logg, app)
}

func start(ctx context.Context, cfg *cfg.Config, logg *logger.Logger, a *app.App) {
	g := internalgrpc.NewServer(logg, a, cfg.GrpsServ)

	go func() {
		<-ctx.Done()
		if err := g.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	if err := g.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
	}
}
