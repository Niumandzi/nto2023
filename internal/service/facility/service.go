package facility

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type FacilityService struct {
	facilityRepo   repository.FacilityRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewFacilityService(fac repository.FacilityRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) FacilityService {
	return FacilityService{
		facilityRepo:   fac,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
