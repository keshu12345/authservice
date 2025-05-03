package pkg

import (
	"github.com/keshucs12345/authservice/pkg/endpoint"
	"github.com/keshucs12345/authservice/pkg/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(endpoint.RegisterEndpoint),
	fx.Provide(service.NewAuthService),
)
