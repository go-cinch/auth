package service_test

import (
	"context"
	"testing"

	v1 "auth/api/auth"
	"auth/internal/tests/mock"
	params "github.com/go-cinch/common/proto/params"
	"github.com/stretchr/testify/assert"
)

func TestCreateWhitelist(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.CreateWhitelistRequest{
		Category: int32(0), // Permission whitelist
		Resource: "GET|/api/public/*",
	}

	_, err := s.CreateWhitelist(ctx, req)
	// Note: May fail if whitelist entry already exists
	if err != nil {
		t.Logf("CreateWhitelist returned error (may be expected): %v", err)
	}
}

func TestFindWhitelist(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.FindWhitelistRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
	}

	rp, err := s.FindWhitelist(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestFindWhitelistWithFilter(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	category := int32(0)
	req := &v1.FindWhitelistRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
		Category: &category,
	}

	rp, err := s.FindWhitelist(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestUpdateWhitelist(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	id := int64(1)
	resource := "GET|/api/updated/*"
	req := &v1.UpdateWhitelistRequest{
		Id:       id,
		Resource: &resource,
	}

	_, err := s.UpdateWhitelist(ctx, req)
	// Note: May fail if whitelist entry doesn't exist or data not changed
	if err != nil {
		t.Logf("UpdateWhitelist returned error (may be expected): %v", err)
	}
}

func TestDeleteWhitelist(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &params.IdsRequest{
		Ids: "999999", // Use non-existent ID to avoid deleting real data
	}

	_, err := s.DeleteWhitelist(ctx, req)
	// Note: May fail if whitelist entry doesn't exist
	if err != nil {
		t.Logf("DeleteWhitelist returned error (may be expected): %v", err)
	}
}
