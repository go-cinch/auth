package service

import (
	"context"

	"auth/api/auth"
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
