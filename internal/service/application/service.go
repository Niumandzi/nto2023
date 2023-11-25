package application

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type ApplicationService struct {
	applicationRepo repository.ApplicationRepository
	contextTimeout  time.Duration
	logger          logging.Logger
	ctx             context.Context
}

func NewApplicationService(app repository.ApplicationRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) ApplicationService {
	return ApplicationService{
		applicationRepo: app,
		contextTimeout:  timeout,
		logger:          logger,
		ctx:             ctx,
	}
}
