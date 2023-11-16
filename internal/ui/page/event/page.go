package event

import (
	"context"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type EventPage struct {
	eventServ service.EventService
	ctx       context.Context
	logger    logging.Logger
}

func NewEventPage(event service.EventService, ctx context.Context, logger logging.Logger) EventPage {
	return EventPage{
		eventServ: event,
		ctx:       ctx,
		logger:    logger,
	}
}
