package work

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type WorkTypeService struct {
	workTypeRepo   repository.WorkTypeRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewWorkTypeService(work repository.WorkTypeRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) WorkTypeService {
	return WorkTypeService{
		workTypeRepo:   work,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
