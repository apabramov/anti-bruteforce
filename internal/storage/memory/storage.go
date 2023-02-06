package memorystorage

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/apabramov/anti-bruteforce/internal/storage"
)

var ErrInvalidIP = errors.New("invalid IP")

type Storage struct {
	blacklist map[string]string
	whitelist map[string]string
	mu        sync.RWMutex
}

func New() *Storage {
	return &Storage{
		blacklist: make(map[string]string),
		whitelist: make(map[string]string),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) AddWhiteList(ctx context.Context, subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, _, err := net.ParseCIDR(subnet); err != nil {
		return err
	}

	if _, ok := s.whitelist[subnet]; ok {
		return storage.ErrExists
	}
	s.whitelist[subnet] = subnet
	return nil
}

func (s *Storage) AddBlackList(ctx context.Context, subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, _, err := net.ParseCIDR(subnet); err != nil {
		return err
	}

	if _, ok := s.blacklist[subnet]; ok {
		return storage.ErrExists
	}
	s.blacklist[subnet] = subnet
	return nil
}

func (s *Storage) DeleteWhiteList(ctx context.Context, subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.whitelist[subnet]; !ok {
		return storage.ErrNotExists
	}
	delete(s.whitelist, subnet)
	return nil
}

func (s *Storage) DeleteBlackList(ctx context.Context, subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.blacklist[subnet]; !ok {
		return storage.ErrNotExists
	}
	delete(s.blacklist, subnet)
	return nil
}

func (s *Storage) CheckIPBlackList(ctx context.Context, ip string) (bool, error) {
	var i net.IP
	if i = net.ParseIP(ip); i == nil {
		return false, ErrInvalidIP
	}

	for l := range s.blacklist {
		_, sb, err := net.ParseCIDR(l)
		if err != nil {
			return false, err
		}
		if sb.Contains(i) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) CheckIPWhiteList(ctx context.Context, ip string) (bool, error) {
	var i net.IP
	if i = net.ParseIP(ip); i == nil {
		return false, ErrInvalidIP
	}
	for l := range s.whitelist {
		_, sb, err := net.ParseCIDR(l)
		if err != nil {
			return false, err
		}
		if sb.Contains(net.ParseIP(ip)) {
			return true, nil
		}
	}
	return false, nil
}
