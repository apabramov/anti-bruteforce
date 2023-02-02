package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/apabramov/anti-bruteforce/internal/config"
	"github.com/apabramov/anti-bruteforce/internal/server/pb"
	"github.com/apabramov/anti-bruteforce/internal/storage"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedEventServiceServer
	App  Application
	Addr string
	Log  Logger
	Srv  *grpc.Server
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Debug(msg string)
	Error(msg string)
}

type Application interface {
	AddWhiteListEvent(ctx context.Context, subnet string) error
	AddBlackListEvent(ctx context.Context, subnet string) error
	DeleteWhiteListEvent(ctx context.Context, subnet string) error
	DeleteBlackListEvent(ctx context.Context, subnet string) error

	AuthEvent(ctx context.Context, auth storage.Authorize) (bool, error)
	ResetEvent(ctx context.Context, auth storage.Authorize) error
}

func NewServer(log Logger, app Application, cfg config.GrpcServerConf) *Server {
	s := &Server{
		App:  app,
		Addr: net.JoinHostPort(cfg.Host, cfg.Port),
		Log:  log,
	}

	g := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingMiddleware(log),
		),
	)

	s.Srv = g
	pb.RegisterEventServiceServer(g, s)

	return s
}

func (s *Server) Start() error {
	list, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.Log.Info(fmt.Sprintf("GRPC starting: %s", s.Addr))
	return s.Srv.Serve(list)
}

func (s *Server) Stop() error {
	s.Log.Info(fmt.Sprintf("GRPC stopping: %s", s.Addr))
	s.Srv.GracefulStop()
	return nil
}

func (s *Server) AddBlackList(ctx context.Context, r *pb.SubnetRequest) (*pb.ResultResponse, error) {
	if err := s.App.AddBlackListEvent(ctx, r.Subnet); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}
	return &pb.ResultResponse{}, nil
}

func (s *Server) AddWhiteList(ctx context.Context, r *pb.SubnetRequest) (*pb.ResultResponse, error) {
	if err := s.App.AddWhiteListEvent(ctx, r.Subnet); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}
	return &pb.ResultResponse{}, nil
}

func (s *Server) DeleteBlackList(ctx context.Context, r *pb.SubnetRequest) (*pb.ResultResponse, error) {
	if err := s.App.DeleteBlackListEvent(ctx, r.Subnet); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}
	return &pb.ResultResponse{}, nil
}

func (s *Server) DeleteWhiteList(ctx context.Context, r *pb.SubnetRequest) (*pb.ResultResponse, error) {
	if err := s.App.DeleteWhiteListEvent(ctx, r.Subnet); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}
	return &pb.ResultResponse{}, nil
}

func (s *Server) Auth(ctx context.Context, r *pb.AuthRequest) (*pb.AuthResponse, error) {
	auth := storage.NewAuthorize(r.GetLogin(), r.GetPassword(), r.GetIp())

	res, err := s.App.AuthEvent(ctx, auth)
	if err != nil {
		return &pb.AuthResponse{Result: false, Error: err.Error()}, err
	}
	return &pb.AuthResponse{Result: res}, err
}

func (s *Server) Reset(ctx context.Context, r *pb.AuthRequest) (*pb.ResultResponse, error) {
	auth := storage.NewAuthorize(r.GetLogin(), r.GetPassword(), r.GetIp())

	if err := s.App.ResetEvent(ctx, auth); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}
	return &pb.ResultResponse{}, nil
}
