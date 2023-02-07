package sqlstorage

import (
	"context"

	"github.com/apabramov/anti-bruteforce/internal/config"
	"github.com/apabramov/anti-bruteforce/internal/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Storage struct {
	DB  *sqlx.DB
	Log *logger.Logger
}

func New(log *logger.Logger, conf config.StorageConf) (*Storage, error) {
	db, err := sqlx.Open("postgres", conf.Dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{Log: log, DB: db}, nil
}

func (s *Storage) Close() error {
	if err := s.DB.Close(); err != nil {
		s.Log.Info(errors.Wrap(err, "err closing db connection").Error())
		return err
	}
	s.Log.Info("db connection gracefully closed")
	return nil
}

func (s *Storage) AddWhiteList(ctx context.Context, subnet string) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO whitelist (subnet) VALUES ($1)", subnet)
	return err
}

func (s *Storage) AddBlackList(ctx context.Context, subnet string) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO blacklist (subnet) VALUES ($1)", subnet)
	return err
}

func (s *Storage) DeleteBlackList(ctx context.Context, subnet string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM blacklist WHERE subnet = $1", subnet)
	return err
}

func (s *Storage) DeleteWhiteList(ctx context.Context, subnet string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM whitelist WHERE subnet = $1", subnet)
	return err
}

func (s *Storage) CheckIPBlackList(ctx context.Context, ip string) (bool, error) {
	var cnt int
	if err := s.DB.Get(&cnt, "SELECT count(1) FROM blacklist WHERE subnet >> $1", ip); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *Storage) CheckIPWhiteList(ctx context.Context, ip string) (bool, error) {
	var cnt int
	if err := s.DB.Get(&cnt, "SELECT count(1) FROM whitelist WHERE subnet >> $1", ip); err != nil {
		return false, err
	}
	return cnt > 0, nil
}
