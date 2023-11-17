package event

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type EventService struct {
	eventRepo      repository.EventRepository
	detailsRepo    repository.DetailsRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewEventService(event repository.EventRepository, details repository.DetailsRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) EventService {
	return EventService{
		eventRepo:      event,
		detailsRepo:    details,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
