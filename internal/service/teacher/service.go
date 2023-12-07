package teacher

import (
	"context"
	"github.com/niumandzi/nto2023/internal/repository"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type TeacherService struct {
	teacherRepo    repository.TeacherRepository
	contextTimeout time.Duration
	logger         logging.Logger
	ctx            context.Context
}

func NewTeacherService(teach repository.TeacherRepository, timeout time.Duration, logger logging.Logger, ctx context.Context) TeacherService {
	return TeacherService{
		teacherRepo:    teach,
		contextTimeout: timeout,
		logger:         logger,
		ctx:            ctx,
	}
}
