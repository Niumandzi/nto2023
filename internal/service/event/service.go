package event

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type EventService struct {
	eventRepo      repository.EventRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewEventService(event repository.EventRepository, timeout time.Duration, logger logging.Logger) EventService {
	return EventService{
		eventRepo:      event,
		contextTimeout: timeout,
		logger:         logger,
	}
}
