package service

import (
	v1 "auth/api/auth/v1"
	"context"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) Idempotent(ctx context.Context, req *emptypb.Empty) (rp *v1.IdempotentReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Idempotent")
	defer span.End()
	rp = &v1.IdempotentReply{}
	rp.Token = s.idempotent.Token(ctx)
	return
}
