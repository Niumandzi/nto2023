package part

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type PartService struct {
	partRepo       repository.PartRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewPartService(part repository.PartRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) PartService {
	return PartService{
		partRepo:       part,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
