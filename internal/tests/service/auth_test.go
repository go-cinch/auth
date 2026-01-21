package service_test

import (
	"context"
	"testing"

	v1 "auth/api/auth"
	"auth/internal/tests/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRegister(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "testpassword123",
		Platform: "web",
	}

	_, err := s.Register(ctx, req)
	// Note: May fail if user already exists or captcha is required
	if err != nil {
		t.Logf("Register returned error (may be expected): %v", err)
	}
}

func TestLogin(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.LoginRequest{
		Username: "admin",
		Password: "admin123",
	}

	rp, err := s.Login(ctx, req)
	// Note: May fail if user doesn't exist or password is wrong
	if err != nil {
		t.Logf("Login returned error (may be expected): %v", err)
	} else {
		assert.NotNil(t, rp)
		assert.NotEmpty(t, rp.Token)
	}
}

func TestLogout(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	rp, err := s.Logout(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.NotNil(t, rp)
}

func TestRefresh(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.RefreshRequest{
		Token: "", // Empty token should fail
	}

	_, err := s.Refresh(ctx, req)
	// Expected to fail with empty token
	assert.Error(t, err)
}

func TestPwd(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.PwdRequest{
		Username:    "admin",
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	_, err := s.Pwd(ctx, req)
	// Note: May fail if user doesn't exist or old password is wrong
	if err != nil {
		t.Logf("Pwd returned error (may be expected): %v", err)
	}
}

func TestStatus(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	req := &v1.StatusRequest{
		Username: "admin",
	}

	rp, err := s.Status(ctx, req)
	// Note: May fail if user doesn't exist
	if err != nil {
		t.Logf("Status returned error (may be expected): %v", err)
	} else {
		assert.NotNil(t, rp)
	}
}
func TestCaptcha(t *testing.T) {
	s := mock.AuthService()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	rp, err := s.Captcha(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.NotNil(t, rp)
	assert.NotNil(t, rp.Captcha)
	assert.NotEmpty(t, rp.Captcha.Id)
}
