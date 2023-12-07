package registration

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type RegistrationService struct {
	registrationRepo repository.RegistrationRepository
	contextTimeout   time.Duration
	logger           logging.Logger
	ctx              context.Context
}

func NewRegistrationService(reg repository.RegistrationRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) RegistrationService {
	return RegistrationService{
		registrationRepo: reg,
		contextTimeout:   timeout,
		logger:           logger,
		ctx:              ctx,
	}
}
