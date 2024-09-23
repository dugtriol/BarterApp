package graph

import (
	"log/slog"

	"github.com/dugtriol/BarterApp/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log      *slog.Logger
	Services *service.Services
}
