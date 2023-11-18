package event

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type DetailsService struct {
	detailsRepo    repository.DetailsRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewDetailsService(det repository.DetailsRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) DetailsService {
	return DetailsService{
		detailsRepo:    det,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
