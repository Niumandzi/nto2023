package event

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type EventService struct {
	eventRepo        repository.EventRepository
	categoryTypeRepo repository.CategoryTypeRepository
	contextTimeout   time.Duration
	logger           logging.Logger
	ctx              context.Context
}

func NewEventService(event repository.EventRepository, categoryType repository.CategoryTypeRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) EventService {
	return EventService{
		eventRepo:        event,
		categoryTypeRepo: categoryType,
		contextTimeout:   timeout,
		logger:           logger,
		ctx:              ctx,
	}
}
