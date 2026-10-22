package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"
	"user-service/utils/conv"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	CreateUserAccount(ctx context.Context, req entity.UserEntity) error
	UpdateUserVerified(ctx context.Context, userID int64) (*entity.UserEntity, error)
	UpdatePasswordByID(ctx context.Context, req entity.UserEntity) error
	GetUserByID(ctx context.Context, userID int64) (*entity.UserEntity, error)
	UpdateDataUser(ctx context.Context, req entity.UserEntity) error

	// Modul Customer Admin
	GetCustomerAll(ctx context.Context, query entity.QueryStringCustomer) ([]entity.UserEntity, int64, int64, error)
	GetCustomerByID(ctx context.Context, userID int64) (*entity.UserEntity, error)
	CreateCustomer(ctx context.Context, req entity.UserEntity) error
	UpdateCustomer(ctx context.Context, req entity.UserEntity) error
	DeleteCustomer(ctx context.Context, customerID int64) error
}

type userRepository struct {
	db *gorm.DB
}

// DeleteCustomer implements [UserRepositoryInterface].
func (u *userRepository) DeleteCustomer(ctx context.Context, customerID int64) error {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id = ?", customerID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Errorf("[UserRepository-1] DeleteCustomer: %v", err)
			return err
		}
		log.Errorf("[UserRepository-1] DeleteCustomer: %v", err)
		return err
	}

	if err := u.db.WithContext(ctx).Delete(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-2] DeleteCustomer: %v", err)
		return err
	}
	return nil
}

// UpdateCustomer implements [UserRepositoryInterface].
func (u *userRepository) UpdateCustomer(ctx context.Context, req entity.UserEntity) error {
	if req.Password != "" {
		password, err := conv.HashPassword(req.Password)
		if err != nil {
			log.Errorf("[UserRepository-1] UpdateCustomer: %v", err)
			return err
		}
		req.Password = password
	}

	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id = ?", req.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Errorf("[UserRepository-2] UpdateCustomer: %v", err)
			return err
		}
		log.Errorf("[UserRepository-2] UpdateCustomer: %v", err)
		return err
	}

	modelUser.Name = req.Name
	modelUser.Email = req.Email
	modelUser.Address = req.Address
	modelUser.Phone = req.Phone
	modelUser.Photo = req.Photo
	modelUser.Lat = req.Lat
	modelUser.Lng = req.Lng
	if req.Password != "" {
		modelUser.Password = req.Password
	}
	if err := u.db.WithContext(ctx).Save(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-3] UpdateCustomer: %v", err)
		return err
	}
	return nil
}

// CreateCustomer implements [UserRepositoryInterface].
func (u *userRepository) CreateCustomer(ctx context.Context, req entity.UserEntity) error {
	modelRole := model.Role{}
	err := u.db.WithContext(ctx).Where("id = ?", req.RoleID).First(&modelRole).Error
	if err != nil {
		log.Errorf("[UserRepository-1] CreateCustomer: %v", err)
		return err
	}

	modelUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
		Lat:      req.Lat,
		Lng:      req.Lng,
		Roles:    []model.Role{modelRole},
	}
	if err := u.db.WithContext(ctx).Create(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-1] CreateCustomer: %v", err)
		return err
	}
	return nil
}

// GetCustomerByID implements [UserRepositoryInterface].
func (u *userRepository) GetCustomerByID(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id = ? AND is_verified = TRUE", userID).Preload("Roles", "name = ?", "Customer").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Errorf("[UserRepository-1] GetCustomerByID: %v", err)
			return nil, err
		}
		log.Errorf("[UserRepository-1] GetCustomerByID: %v", err)
		return nil, err
	}
	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      modelUser.Email,
		RoleName:   modelUser.Roles[0].Name,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Address:    modelUser.Address,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: modelUser.IsVerified,
		RoleID:     modelUser.Roles[0].ID,
	}, nil
}

// GetAllUser implements [UserRepositoryInterface].
func (u *userRepository) GetCustomerAll(ctx context.Context, query entity.QueryStringCustomer) ([]entity.UserEntity, int64, int64, error) {
	modelUsers := []model.User{}
	var countData int64
	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := u.db.WithContext(ctx).Model(&model.User{}).Preload("Roles", "name = ?", "Customer").
		Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%").
		Order(order).
		Offset(offset).
		Limit(query.Limit)

	if err := sqlMain.Model(&modelUsers).Count(&countData).Error; err != nil {
		log.Errorf("[UserRepository-1] GetCustomerAll: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.Order(order).Limit(query.Limit).Find(&modelUsers).Error; err != nil {
		log.Errorf("[UserRepository-1] GetCustomerAll: %v", err)
		return nil, 0, 0, err
	}

	userEntities := make([]entity.UserEntity, len(modelUsers))
	for i, modelUser := range modelUsers {
		userEntities[i] = entity.UserEntity{
			ID:         modelUser.ID,
			Name:       modelUser.Name,
			Email:      modelUser.Email,
			Phone:      modelUser.Phone,
			Photo:      modelUser.Photo,
			RoleName:   modelUser.Roles[0].Name,
			IsVerified: modelUser.IsVerified,
		}
		if len(modelUser.Roles) > 0 {
			err := errors.New("404")
			if len(modelUser.Roles) == 0 {
				log.Errorf("[UserRepository-2] GetCustomerAll: %v", err)
				return nil, 0, 0, err
			}
			userEntities[i].RoleName = modelUser.Roles[0].Name
		}
	}

	return userEntities, countData, int64(totalPage), nil
}

// UpdateDataUser implements [UserRepositoryInterface].
func (u *userRepository) UpdateDataUser(ctx context.Context, req entity.UserEntity) error {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id = ? AND is_verified = TRUE", req.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Errorf("[UserRepository-1] UpdateDataUser: %v", err)
			return err
		}
		log.Errorf("[UserRepository-1] UpdateDataUser: %v", err)
		return err
	}

	modelUser.Name = req.Name
	modelUser.Email = req.Email
	modelUser.Address = req.Address
	modelUser.Phone = req.Phone
	modelUser.Photo = req.Photo
	modelUser.Lat = req.Lat
	modelUser.Lng = req.Lng
	if err := u.db.WithContext(ctx).Save(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-2] UpdateDataUser: %v", err)
		return err
	}
	return nil
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
