package bucket

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/apabramov/anti-bruteforce/internal/config"
	internalredis "github.com/apabramov/anti-bruteforce/internal/redis"
	"github.com/apabramov/anti-bruteforce/internal/storage"
)

const (
	loginPref = "login"
	passPref  = "password"
	ipPref    = "ip"
)

type LimitBucket struct {
	cache      internalredis.Cache
	loginLimit int
	passLimit  int
	ipLimit    int
}

func New(cache internalredis.Cache, conf config.LimitConf) *LimitBucket {
	return &LimitBucket{
		cache:      cache,
		loginLimit: conf.LimitLogin,
		passLimit:  conf.LimitPass,
		ipLimit:    conf.LimitIP,
	}
}

func getKey(pref string, val string, m string) string {
	return fmt.Sprintf("%s_%s_%s", pref, val, m)
}

func (l *LimitBucket) CheckLimit(ctx context.Context, auth storage.Authorize) (bool, error) {
	m := strconv.Itoa(time.Now().Minute())
	d := time.Minute

	kl := getKey(loginPref, auth.Login, m)
	if err := l.cache.Incr(ctx, kl, d); err != nil {
		return false, err
	}
	v, err := l.cache.Get(ctx, kl)
	if err != nil {
		return false, err
	}
	if v > l.loginLimit {
		return false, nil
	}

	kp := getKey(passPref, auth.Password, m)
	if err = l.cache.Incr(ctx, kp, d); err != nil {
		return false, err
	}
	v, err = l.cache.Get(ctx, kp)
	if err != nil {
		return false, err
	}
	if v > l.passLimit {
		return false, nil
	}

	ki := getKey(ipPref, auth.IP, m)
	if err = l.cache.Incr(ctx, ki, d); err != nil {
		return false, err
	}
	v, err = l.cache.Get(ctx, ki)
	if err != nil {
		return false, err
	}
	if v > l.ipLimit {
		return false, nil
	}

	return true, nil
}

func (l *LimitBucket) ResetBucket(ctx context.Context, auth storage.Authorize) error {
	m := strconv.Itoa(time.Now().Minute())

	if err := l.cache.Del(ctx, getKey(loginPref, auth.Login, m)); err != nil {
		return err
	}

	if err := l.cache.Del(ctx, getKey(passPref, auth.Password, m)); err != nil {
		return err
	}

	return l.cache.Del(ctx, getKey(ipPref, auth.IP, m))
}
