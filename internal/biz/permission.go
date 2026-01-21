package biz

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"

	"auth/internal/conf"
)

type Permission struct {
	Resources []string `json:"resources"`
	Menus     []string `json:"menus"`
	Btns      []string `json:"btns"`
}

type CheckPermission struct {
	UserID   int64  `json:"userId,string"`
	Resource string `json:"resource"`
	Method   string `json:"method"`
	URI      string `json:"uri"`
}

type PermissionUseCase struct {
	c    *conf.Bootstrap
	repo PermissionRepo
}

func NewPermissionUseCase(c *conf.Bootstrap, repo PermissionRepo) *PermissionUseCase {
	return &PermissionUseCase{
		c:    c,
		repo: repo,
	}
}

// Check validates a request against the caller's permissions.
// When Method is set, URI is used as the resource to check.
func (uc *PermissionUseCase) Check(ctx context.Context, req *CheckPermission) (ok bool, err error) {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Check")
	defer span.End()

	if req == nil {
		return false, ErrIllegalParameter(ctx, "permission")
	}
	if req.UserID == 0 {
		return false, ErrIllegalParameter(ctx, "userID")
	}

	resource := strings.TrimSpace(req.Resource)
	method := strings.TrimSpace(req.Method)
	if method != "" {
		resource = strings.TrimSpace(req.URI)
	}

	return uc.repo.CheckPermission(ctx, req.UserID, resource, method)
}

func (uc *PermissionUseCase) CheckPermission(ctx context.Context, userID int64, resource, method string) (bool, error) {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "CheckPermission")
	defer span.End()

	return uc.repo.CheckPermission(ctx, userID, resource, method)
}

func (uc *PermissionUseCase) GetByUserID(ctx context.Context, userID int64) (*Permission, error) {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "GetByUserID")
	defer span.End()

	if userID == 0 {
		return nil, ErrIllegalParameter(ctx, "userID")
	}

	actions, err := uc.repo.FindUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	rp := &Permission{
		Resources: make([]string, 0),
		Menus:     make([]string, 0),
		Btns:      make([]string, 0),
	}

	seenRes := make(map[string]struct{})
	seenMenu := make(map[string]struct{})
	seenBtn := make(map[string]struct{})

	addLines := func(dst *[]string, seen map[string]struct{}, s *string) {
		if s == nil {
			return
		}
		for _, line := range strings.Split(*s, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if _, ok := seen[line]; ok {
				continue
			}
			seen[line] = struct{}{}
			*dst = append(*dst, line)
		}
	}

	for _, a := range actions {
		if a == nil {
			continue
		}
		addLines(&rp.Resources, seenRes, a.Resource)
		addLines(&rp.Menus, seenMenu, a.Menu)
		addLines(&rp.Btns, seenBtn, a.Btn)
	}

	return rp, nil
}
