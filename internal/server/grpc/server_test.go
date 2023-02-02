package internalgrpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/apabramov/anti-bruteforce/internal/app"
	"github.com/apabramov/anti-bruteforce/internal/bucket"
	"github.com/apabramov/anti-bruteforce/internal/config"
	"github.com/apabramov/anti-bruteforce/internal/logger"
	"github.com/apabramov/anti-bruteforce/internal/redis"
	"github.com/apabramov/anti-bruteforce/internal/server/pb"
	ms "github.com/apabramov/anti-bruteforce/internal/storage/memory"
)

var (
	client *internalredis.RedisClient
	cfg    config.Config
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	c := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	client = &internalredis.RedisClient{Client: c}

	logg, err := logger.New("info")
	if err != nil {
		log.Println(err)
	}

	cfg = config.Config{
		GrpsServ: config.GrpcServerConf{
			Host: "",
			Port: "9000",
		},
		Redis: config.RedisConf{
			Host: "",
			Port: "6379",
		},
		Limit: config.LimitConf{
			LimitLogin: 1,
			LimitIP:    100,
			LimitPass:  1000,
		},
	}

	storage := ms.New()

	b := bucket.New(client, cfg.Limit)
	a := app.New(logg, storage, b)

	go runGrpc(a, &cfg, logg)

	code := m.Run()
	os.Exit(code)
}

func TestGRPC(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	clientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	_, err = grpc.Dial(l.Addr().String(), clientOptions...)
	require.NoError(t, err)

	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := ms.New()

	b := bucket.New(client, cfg.Limit)
	a := app.New(logg, storage, b)

	server := NewServer(logg, a, cfg.GrpsServ)
	require.NoError(t, err)

	go func() {
		server.Srv.Serve(l)
	}()
}

func TestGRPCServerAddBlackList(t *testing.T) {
	t.Run("add blacklist", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.AddBlackList(ctx, &pb.SubnetRequest{Subnet: "192.168.1.0/24"})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())
	})
}

func TestGRPCServerAddWhiteList(t *testing.T) {
	t.Run("add whitelist", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.AddWhiteList(ctx, &pb.SubnetRequest{Subnet: "192.168.1.0/24"})
		require.NoError(t, err)
		require.Equal(t, "", r.GetError())
	})
}

//nolint:dupl
func TestGRPCServerDeleteBlackList(t *testing.T) {
	t.Run("delete blacklist", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		res, err := c.AddBlackList(ctx, &pb.SubnetRequest{Subnet: "127.0.1.0/24"})
		require.NoError(t, err)
		require.Equal(t, "", res.GetError())

		res, err = c.DeleteBlackList(ctx, &pb.SubnetRequest{Subnet: "127.0.1.0/24"})
		require.NoError(t, err)
		require.Equal(t, "", res.GetError())
	})
}

//nolint:dupl
func TestGRPCServerDeleteWhiteList(t *testing.T) {
	t.Run("delete whitelist", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		res, err := c.AddWhiteList(ctx, &pb.SubnetRequest{Subnet: "10.10.1.0/24"})
		require.NoError(t, err)
		require.Equal(t, "", res.GetError())

		res, err = c.DeleteWhiteList(ctx, &pb.SubnetRequest{Subnet: "10.10.1.0/24"})
		require.NoError(t, err)
		require.Equal(t, "", res.GetError())
	})
}

func TestGRPCServerAuth(t *testing.T) {
	t.Run("auth limit login", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx := context.Background()

		res, err := c.Auth(ctx, &pb.AuthRequest{Login: "login_lim", Password: "pass", Ip: "198.168.25.25"})
		require.NoError(t, err)
		require.True(t, res.GetResult())

		res, err = c.Auth(ctx, &pb.AuthRequest{Login: "login_lim", Password: "pass1", Ip: "127.168.25.25"})
		require.NoError(t, err)
		require.False(t, res.GetResult())
	})
}

func TestGRPCServerReset(t *testing.T) {
	t.Run("auth reset login", func(t *testing.T) {
		conn, err := grpc.Dial(net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		res, err := c.Auth(ctx, &pb.AuthRequest{Login: "login", Password: "pass", Ip: "198.168.25.25"})
		require.NoError(t, err)
		require.True(t, res.GetResult())

		result, err := c.Reset(ctx, &pb.AuthRequest{
			Login:    "login",
			Password: "pass0",
			Ip:       "198.168.25.25",
		})
		require.NoError(t, err)
		require.True(t, result.GetError() == "")

		res, err = c.Auth(ctx, &pb.AuthRequest{Login: "login", Password: "pass1", Ip: "195.165.25.25"})
		require.NoError(t, err)
		require.True(t, res.GetResult())
	})
}

func runGrpc(app *app.App, cfg *config.Config, log *logger.Logger) {
	srv := NewServer(log, app, cfg.GrpsServ)

	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GrpsServ.Host, cfg.GrpsServ.Port))
	if err != nil {
		log.Info(err.Error())
	}
	if err := srv.Srv.Serve(lis); err != nil {
		panic(err)
	}
}
