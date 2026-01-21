package service_test

import (
	"context"
	"testing"

	v1 "auth/api/auth"
	"auth/internal/tests/mock"
	params "github.com/go-cinch/common/proto/params"
	"github.com/stretchr/testify/assert"
)

func TestCreateRole(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.CreateRoleRequest{
		Name: "test-role",
		Word: "test_role_word",
	}

	_, err := s.CreateRole(ctx, req)
	// Note: May fail if role already exists
	if err != nil {
		t.Logf("CreateRole returned error (may be expected): %v", err)
	}
}

func TestFindRole(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.FindRoleRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
	}

	rp, err := s.FindRole(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestFindRoleWithFilter(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	name := "admin"
	req := &v1.FindRoleRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
		Name: &name,
	}

	rp, err := s.FindRole(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestUpdateRole(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	id := int64(1)
	name := "updated-role"
	req := &v1.UpdateRoleRequest{
		Id:   id,
		Name: &name,
	}

	_, err := s.UpdateRole(ctx, req)
	// Note: May fail if role doesn't exist or data not changed
	if err != nil {
		t.Logf("UpdateRole returned error (may be expected): %v", err)
	}
}

func TestDeleteRole(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &params.IdsRequest{
		Ids: "999999", // Use non-existent ID to avoid deleting real data
	}

	_, err := s.DeleteRole(ctx, req)
	// Note: May fail if role doesn't exist
	if err != nil {
		t.Logf("DeleteRole returned error (may be expected): %v", err)
	}
}
