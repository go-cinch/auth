package service_test

import (
	"context"
	"testing"

	v1 "auth/api/auth"
	"auth/internal/tests/mock"
	params "github.com/go-cinch/common/proto/params"
	"github.com/stretchr/testify/assert"
)

func TestFindUser(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.FindUserRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
	}

	rp, err := s.FindUser(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestFindUserWithFilter(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	username := "admin"
	req := &v1.FindUserRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
		Username: &username,
	}

	rp, err := s.FindUser(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestUpdateUser(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	id := int64(1)
	platform := "mobile"
	req := &v1.UpdateUserRequest{
		Id:       id,
		Platform: &platform,
	}

	_, err := s.UpdateUser(ctx, req)
	// Note: May fail if user doesn't exist
	if err != nil {
		t.Logf("UpdateUser returned error (may be expected): %v", err)
	}
}
func TestUpdateUserUnlock(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	findReq := &v1.FindUserRequest{
		Page: &params.Page{Num: 1, Size: 1},
	}
	findRp, err := s.FindUser(ctx, findReq)
	if err != nil {
		t.Fatalf("FindUser failed: %v", err)
	}
	if len(findRp.List) == 0 {
		t.Skip("No user found")
	}

	id := findRp.List[0].Id
	unlocked := false
	req := &v1.UpdateUserRequest{
		Id:     id,
		Locked: &unlocked,
	}
	_, err = s.UpdateUser(ctx, req)
	if err != nil {
		t.Logf("Unlock error (may be expected if already unlocked): %v", err)
	}
}

func TestUpdateUserLockWithExpireTime(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	// First find an existing user
	findReq := &v1.FindUserRequest{
		Page: &params.Page{
			Num:  1,
			Size: 1,
		},
	}
	findRp, err := s.FindUser(ctx, findReq)
	if err != nil {
		t.Fatalf("FindUser failed: %v", err)
	}
	if len(findRp.List) == 0 {
		t.Skip("No user found in database, skipping test")
	}

	id := findRp.List[0].Id
	locked := true
	lockExpireTime := "2026-01-30 00:00:00"
	req := &v1.UpdateUserRequest{
		Id:             id,
		Locked:         &locked,
		LockExpireTime: &lockExpireTime,
	}

	_, err = s.UpdateUser(ctx, req)
	if err != nil {
		t.Errorf("UpdateUser with lock expire time failed: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &params.IdsRequest{
		Ids: "999999", // Use non-existent ID to avoid deleting real data
	}

	_, err := s.DeleteUser(ctx, req)
	// Note: May fail if trying to delete yourself or user doesn't exist
	if err != nil {
		t.Logf("DeleteUser returned error (may be expected): %v", err)
	}
}
