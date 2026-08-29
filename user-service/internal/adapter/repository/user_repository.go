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

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	CreateUserAccount(ctx context.Context, req entity.UserEntity) error
	UpdateUserVerified(ctx context.Context, userID int64) (*entity.UserEntity, error)
	UpdatePasswordByID(ctx context.Context, req entity.UserEntity) error
	GetUserByID(ctx context.Context, userID int64) (*entity.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

// GetUserByID implements [UserRepositoryInterface].
func (u *userRepository) GetUserByID(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id = ? AND is_verified = TRUE", userID).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Errorf("[UserRepository-1] GetUserByID: %v", err)
			return nil, err
		}
		log.Errorf("[UserRepository-1] GetUserByID: %v", err)
		return nil, err
	}
	return &entity.UserEntity{
		ID:       modelUser.ID,
		Name:     modelUser.Name,
		Email:    modelUser.Email,
		RoleName: modelUser.Roles[0].Name,
		Lat:      modelUser.Lat,
		Lng:      modelUser.Lng,
		Address:  modelUser.Address,
		Phone:    modelUser.Phone,
		Photo:    modelUser.Photo,
	}, nil

}

// UpdatePasswordByID implements [UserRepositoryInterface].
func (u *userRepository) UpdatePasswordByID(ctx context.Context, req entity.UserEntity) error {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id = ?", req.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("user not found")
			log.Errorf("[UserRepository-1] UpdatePasswordByID: %v", err)
			return err
		}
		log.Errorf("[UserRepository-1] UpdatePasswordByID: %v", err)
		return err
	}
	modelUser.Password = req.Password
	if err := u.db.WithContext(ctx).Save(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-2] UpdatePasswordByID: %v", err)
		return err
	}
	return nil
}

// UpdateUserVerified implements [UserRepositoryInterface].
func (u *userRepository) UpdateUserVerified(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id = ?", userID).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("user not found")
			log.Errorf("[UserRepository-1] UpdateUserVerified: %v", err)
			return nil, err
		}
		log.Errorf("[UserRepository-1] UpdateUserVerified: %v", err)
		return nil, err
	}
	modelUser.IsVerified = true
	if err := u.db.WithContext(ctx).Save(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-2] UpdateUserVerified: %v", err)
		return nil, err
	}
	if len(modelUser.Roles) == 0 {
		err := errors.New("user role not found")

		log.Errorf("[UserRepository-3] UpdateUserVerified: %v", err)

		return nil, err
	}
	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      modelUser.Email,
		Password:   modelUser.Password,
		RoleName:   modelUser.Roles[0].Name,
		Address:    modelUser.Address,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		IsVerified: modelUser.IsVerified,
	}, nil
}

// CreateUserAccount implements [UserRepositoryInterface].
func (u *userRepository) CreateUserAccount(ctx context.Context, req entity.UserEntity) error {
	modelRole := model.Role{}
	err := u.db.WithContext(ctx).Where("name = ?", "Customer").First(&modelRole).Error
	if err != nil {
		log.Errorf("[UserRepository-1] CreateUserAccount: %v", err)
		return err
	}

	modelUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Roles:    []model.Role{modelRole},
	}
	if err := u.db.WithContext(ctx).Create(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-1] CreateUserAccount: %v", err)
		return err
	}
	now := time.Now()
	modelVerify := model.VerificationToken{
		UserID:    modelUser.ID,
		Token:     req.Token,
		TokenType: "email_verification",
		ExpiresAt: now.Add(time.Hour * 1),
	}
	if err := u.db.WithContext(ctx).Create(&modelVerify).Error; err != nil {
		log.Errorf("[UserRepository-2] CreateUserAccount: %v", err)
		return err
	}

	return nil
}

// GetUserByEmail implements [UserRepositoryInterface].
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("email = ? and is_verified = ?", email, true).
		Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("user not found")
			log.Errorf("[UserRepository-1] GetUserByEmail: %v", err)
			return nil, err
		}
		log.Errorf("[UserRepository-1] GetUserByEmail: %v", err)
		return nil, err
	}

	userEntity := &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      modelUser.Email,
		Password:   modelUser.Password,
		RoleName:   modelUser.Roles[0].Name,
		Address:    modelUser.Address,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		IsVerified: modelUser.IsVerified,
	}

	if len(modelUser.Roles) > 0 {
		userEntity.RoleName = modelUser.Roles[0].Name
	}

	return userEntity, nil
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{
		db: db,
	}
}
