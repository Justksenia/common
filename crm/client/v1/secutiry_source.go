package v1

import (
	"context"

	"gitlab.com/adstail/ts-common/crm/client/v1/gen"
)

type ConstJWTSecuritySource struct {
	Token string
}

func NewConstJWTSecuritySource(token string) *ConstJWTSecuritySource {
	return &ConstJWTSecuritySource{Token: token}
}

func (c *ConstJWTSecuritySource) JwtBearer(_ context.Context, _ string) (gen.JwtBearer, error) {
	return gen.JwtBearer{Token: c.Token}, nil
}
