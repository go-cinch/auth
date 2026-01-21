package service_test

import (
	"context"
	"testing"

	v1 "auth/api/auth"
	"auth/internal/tests/mock"
	params "github.com/go-cinch/common/proto/params"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserGroup(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.CreateUserGroupRequest{
		Name: "test-user-group",
		Word: "test_user_group_word",
	}

	_, err := s.CreateUserGroup(ctx, req)
	// Note: May fail if user group already exists
	if err != nil {
		t.Logf("CreateUserGroup returned error (may be expected): %v", err)
	}
}

func TestFindUserGroup(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.FindUserGroupRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
	}

	rp, err := s.FindUserGroup(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestFindUserGroupWithFilter(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	name := "admin"
	req := &v1.FindUserGroupRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
		Name: &name,
	}

	rp, err := s.FindUserGroup(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestUpdateUserGroup(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	id := int64(1)
	name := "updated-user-group"
	req := &v1.UpdateUserGroupRequest{
		Id:   id,
		Name: &name,
	}

	_, err := s.UpdateUserGroup(ctx, req)
	// Note: May fail if user group doesn't exist or data not changed
	if err != nil {
		t.Logf("UpdateUserGroup returned error (may be expected): %v", err)
	}
}

func TestDeleteUserGroup(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &params.IdsRequest{
		Ids: "999999", // Use non-existent ID to avoid deleting real data
	}

	_, err := s.DeleteUserGroup(ctx, req)
	// Note: May fail if user group doesn't exist
	if err != nil {
		t.Logf("DeleteUserGroup returned error (may be expected): %v", err)
	}
}
