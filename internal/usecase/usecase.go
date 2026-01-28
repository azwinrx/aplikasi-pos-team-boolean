package usecase

import (
	"aplikasi-pos-team-boolean/internal/data/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UseCase struct {
	repo repository.Repository
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	return &UseCase{}
}
