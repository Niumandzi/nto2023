package booking

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type BookingService struct {
	bookingRepo    repository.BookingRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewBookingService(book repository.BookingRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) BookingService {
	return BookingService{
		bookingRepo:    book,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
