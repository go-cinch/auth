package service

import (
	"auth/api/auth"
	"context"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) Idempotent(ctx context.Context, req *emptypb.Empty) (rp *auth.IdempotentReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Idempotent")
	defer span.End()
	rp = &auth.IdempotentReply{}
	rp.Token = s.idempotent.Token(ctx)
	return
}
