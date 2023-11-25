package work

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type WorkService struct {
	workRepo       repository.WorkTypeRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewWorkService(work repository.WorkTypeRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) WorkService {
	return WorkService{
		workRepo:       work,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
