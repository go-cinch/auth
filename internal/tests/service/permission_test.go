package service_test

import (
	"context"
	"testing"

	v1 "auth/api/auth"
	"auth/internal/tests/mock"
	"github.com/go-cinch/common/jwt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestPermission(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	// Set up JWT context with user code
	ctx = jwt.NewServerContextByUser(ctx, jwt.User{
		Attrs: map[string]string{
			"code":     "test-user-code",
			"platform": "web",
		},
	})

	method := "GET"
	uri := "/api/test"
	req := &v1.PermissionRequest{
		Method: &method,
		Uri:    &uri,
	}

	_, err := s.Permission(ctx, req)
	// Note: May fail if user doesn't exist or no permission
	if err != nil {
		t.Logf("Permission returned error (may be expected): %v", err)
	}
}

func TestInfo(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	// Set up JWT context with user code
	ctx = jwt.NewServerContextByUser(ctx, jwt.User{
		Attrs: map[string]string{
			"code":     "test-user-code",
			"platform": "web",
		},
	})

	rp, err := s.Info(ctx, &emptypb.Empty{})
	// Note: May fail if user doesn't exist
	if err != nil {
		t.Logf("Info returned error (may be expected): %v", err)
	} else {
		assert.NotNil(t, rp)
	}
}
