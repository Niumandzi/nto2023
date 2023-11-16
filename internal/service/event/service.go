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

func NewEventService(contact repository.EventRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) EventService {
	return EventService{
		eventRepo:      contact,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
