package memorystorage

import (
	"context"
	"errors"
	"testing"

	"github.com/apabramov/anti-bruteforce/internal/storage"
	"github.com/stretchr/testify/require"
)

const sub = "192.168.1.0/24"

func TestAddBlackList(t *testing.T) {
	ctx := context.Background()

	s := New()
	err := s.AddBlackList(ctx, sub)
	require.NoError(t, err)

	err = s.AddBlackList(ctx, sub)
	require.Error(t, err)
	require.True(t, errors.Is(err, storage.ErrExists))
}

func TestCheckBlackList(t *testing.T) {
	ctx := context.Background()

	s := New()
	err := s.AddBlackList(ctx, sub)
	require.NoError(t, err)

	b, err := s.CheckIPBlackList(ctx, "192.168.1.1")
	require.NoError(t, err)
	require.True(t, b)
}

func TestDeleteBlackList(t *testing.T) {
	ctx := context.Background()

	s := New()
	err := s.AddBlackList(ctx, sub)
	require.NoError(t, err)

	err = s.DeleteBlackList(ctx, sub)
	require.NoError(t, err)

	err = s.DeleteBlackList(ctx, sub)
	require.True(t, errors.Is(err, storage.ErrNotExists))
}

func TestAddWhiteList(t *testing.T) {
	ctx := context.Background()

	s := New()
	err := s.AddWhiteList(ctx, sub)
	require.NoError(t, err)

	err = s.AddWhiteList(ctx, sub)
	require.Error(t, err)
	require.True(t, errors.Is(err, storage.ErrExists))
}

func TestCheckWhiteList(t *testing.T) {
	ctx := context.Background()

	s := New()
	err := s.AddWhiteList(ctx, sub)
	require.NoError(t, err)

	b, err := s.CheckIPWhiteList(ctx, "192.168.1.2")
	require.NoError(t, err)
	require.True(t, b)
}

func TestDeleteWhiteList(t *testing.T) {
	ctx := context.Background()

	s := New()
	err := s.AddWhiteList(ctx, sub)
	require.NoError(t, err)

	err = s.DeleteWhiteList(ctx, sub)
	require.NoError(t, err)

	err = s.DeleteWhiteList(ctx, sub)
	require.True(t, errors.Is(err, storage.ErrNotExists))
}
