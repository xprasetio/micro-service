package repository

import (
	"context"

	"user-service/internal/core/domain/entity"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type VerificationTokenRepositoryInterface interface {
	CreateVerificationToken(ctx context.Context, req entity.VerificationToken) error
}

type verificationTokenRepository struct {
	db *gorm.DB
}

func NewVerificationTokenRepository(db *gorm.DB) VerificationTokenRepositoryInterface {
	return &verificationTokenRepository{
		db: db,
	}
}

func (v *verificationTokenRepository) CreateVerificationToken(ctx context.Context, req entity.VerificationToken) error {
	if err := v.db.WithContext(ctx).Create(&req).Error; err != nil {
		log.Errorf("[VerificationTokenRepository-1] CreateVerificationToken: %v", err)
		return err
	}
	return nil
}
