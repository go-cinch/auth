package service

import (
	"context"

	"auth/api/auth"
	"auth/internal/biz"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) Idempotent(ctx context.Context, _ *emptypb.Empty) (rp *auth.IdempotentReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Idempotent")
	defer span.End()
	rp = &auth.IdempotentReply{}
	rp.Token = s.idempotent.Token(ctx)
	return
}

func (s *AuthService) CheckIdempotent(ctx context.Context, req *auth.CheckIdempotentRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CheckIdempotent")
	defer span.End()
	rp = &emptypb.Empty{}
	if req.Token == "" {
		err = biz.ErrIdempotentMissingToken(ctx)
		return
	}
	pass := s.idempotent.Check(ctx, req.Token)
	if !pass {
		err = biz.ErrIdempotentTokenExpired(ctx)
		return
	}
	return
}
