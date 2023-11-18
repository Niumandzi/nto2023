package event

import (
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type DetailsPage struct {
	eventDet service.DetailsService
	logger   logging.Logger
}

func NewDetailsPage(det service.DetailsService, logger logging.Logger) DetailsPage {
	return DetailsPage{
		eventDet: det,
		logger:   logger,
	}
}
