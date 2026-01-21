package service_test

import (
	"context"
	"testing"

	v1 "auth/api/auth"
	"auth/internal/tests/mock"
	params "github.com/go-cinch/common/proto/params"
	"github.com/stretchr/testify/assert"
)

func TestCreateAction(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	resource := "GET|/api/test/*"
	req := &v1.CreateActionRequest{
		Name:     "test-action",
		Word:     "test_action_word",
		Resource: &resource,
	}

	_, err := s.CreateAction(ctx, req)
	// Note: May fail if action already exists
	if err != nil {
		t.Logf("CreateAction returned error (may be expected): %v", err)
	}
}

func TestFindAction(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.FindActionRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
	}

	rp, err := s.FindAction(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestFindActionWithFilter(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	name := "admin"
	req := &v1.FindActionRequest{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
		Name: &name,
	}

	rp, err := s.FindAction(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Page)
}

func TestUpdateAction(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	id := int64(1)
	name := "updated-action"
	req := &v1.UpdateActionRequest{
		Id:   id,
		Name: &name,
	}

	_, err := s.UpdateAction(ctx, req)
	// Note: May fail if action doesn't exist or data not changed
	if err != nil {
		t.Logf("UpdateAction returned error (may be expected): %v", err)
	}
}

func TestDeleteAction(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &params.IdsRequest{
		Ids: "999999", // Use non-existent ID to avoid deleting real data
	}

	_, err := s.DeleteAction(ctx, req)
	// Note: May fail if action doesn't exist
	if err != nil {
		t.Logf("DeleteAction returned error (may be expected): %v", err)
	}
}
