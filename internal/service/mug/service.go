package mug

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type MugTypeService struct {
	mugTypeRepo    repository.MugTypeRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewMugTypeService(mug repository.MugTypeRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) MugTypeService {
	return MugTypeService{
		mugTypeRepo:    mug,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
