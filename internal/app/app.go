package app

import (
	"context"
	"errors"
	
	"github.com/apabramov/anti-bruteforce/internal/bucket"
	"github.com/apabramov/anti-bruteforce/internal/config"
	"github.com/apabramov/anti-bruteforce/internal/logger"
	"github.com/apabramov/anti-bruteforce/internal/storage"
	memorystorage "github.com/apabramov/anti-bruteforce/internal/storage/memory"
	sqlstorage "github.com/apabramov/anti-bruteforce/internal/storage/sql"
)

var ErrTypeStorageNotFound = errors.New("storage type not found")

type App struct {
	Log    Logger
	Store  Storage
	Bucket *bucket.LimitBucket
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

type Storage interface {
	AddWhiteList(ctx context.Context, subnet string) error
	AddBlackList(ctx context.Context, subnet string) error
	DeleteWhiteList(ctx context.Context, subnet string) error
	DeleteBlackList(ctx context.Context, subnet string) error

	CheckIPBlackList(ctx context.Context, ip string) (bool, error)
	CheckIPWhiteList(ctx context.Context, ip string) (bool, error)
}

func NewStorage(log *logger.Logger, cfg config.StorageConf) (Storage, error) {
	var (
		st  Storage
		err error
	)
	switch cfg.Type {
	case "memory":
		st = memorystorage.New()
	case "sql":
		st, err = sqlstorage.New(log, cfg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrTypeStorageNotFound
	}
	return st, nil
}

func New(logger Logger, storage Storage, bucket *bucket.LimitBucket) *App {
	return &App{Log: logger, Store: storage, Bucket: bucket}
}

func (a App) AddWhiteListEvent(ctx context.Context, subnet string) error {
	return a.Store.AddWhiteList(ctx, subnet)
}

func (a App) AddBlackListEvent(ctx context.Context, subnet string) error {
	return a.Store.AddBlackList(ctx, subnet)
}

func (a App) DeleteWhiteListEvent(ctx context.Context, subnet string) error {
	return a.Store.DeleteWhiteList(ctx, subnet)
}

func (a App) DeleteBlackListEvent(ctx context.Context, subnet string) error {
	return a.Store.DeleteBlackList(ctx, subnet)
}

func (a App) AuthEvent(ctx context.Context, auth storage.Authorize) (bool, error) {
	var (
		res bool
		err error
	)

	if res, err = a.Store.CheckIPBlackList(ctx, auth.IP); err != nil {
		return false, err
	}

	if res {
		return false, nil
	}

	if res, err = a.Store.CheckIPWhiteList(ctx, auth.IP); err != nil {
		return false, err
	}

	if res {
		return true, nil
	}

	return a.Bucket.CheckLimit(ctx, auth)
}

func (a App) ResetEvent(ctx context.Context, auth storage.Authorize) error {
	return a.Bucket.ResetBucket(ctx, auth)
}
