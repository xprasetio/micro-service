package repository

import (
	"context"
	"errors"
	"time"

	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type VerificationTokenRepositoryInterface interface {
	CreateVerificationToken(ctx context.Context, req entity.VerificationToken) error
	GetDataByToken(ctx context.Context, token string) (*entity.VerificationToken, error)
}

type verificationTokenRepository struct {
	db *gorm.DB
}

// GetDataByToken implements [VerificationTokenRepositoryInterface].
func (v *verificationTokenRepository) GetDataByToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	modelToken := model.VerificationToken{}
	if err := v.db.WithContext(ctx).Where("token = ?", token).First(&modelToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("token not found")
			log.Errorf("[VerificationTokenRepository-1] GetDataByToken: %v", err)
			return nil, err
		}
		log.Errorf("[VerificationTokenRepository-1] GetDataByToken: %v", err)
		return nil, err
	}

	currenTime := time.Now()

	if modelToken.ExpiresAt.Before(currenTime) {
		err := errors.New("token expired")
		log.Errorf("[VerificationTokenRepository-2] GetDataByToken: %v", err)
		return nil, err
	}

	return &entity.VerificationToken{
		ID:        modelToken.ID,
		UserID:    modelToken.UserID,
		Token:     modelToken.Token,
		TokenType: modelToken.TokenType,
		ExpiresAt: modelToken.ExpiresAt,
	}, nil
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
