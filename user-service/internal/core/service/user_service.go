package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"user-service/config"
	"user-service/internal/adapter/message"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils/conv"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type UserServiceInterface interface {
	SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
	CreateUserAccount(ctx context.Context, req entity.UserEntity) error
	ForgotPassword(ctx context.Context, req entity.UserEntity) error
	VerifyToken(ctx context.Context, token string) (*entity.UserEntity, error)
	UpdatePassword(ctx context.Context, req entity.UserEntity) error
	GetProfileUser(ctx context.Context, userID int64) (*entity.UserEntity, error)
	UpdateDataUser(ctx context.Context, req entity.UserEntity) error
	GetCustomerAll(ctx context.Context, query entity.QueryStringCustomer) ([]entity.UserEntity, int64, int64, error)
	GetCustomerByID(ctx context.Context, userID int64) (*entity.UserEntity, error)
}

type userService struct {
	userRepository repository.UserRepositoryInterface
	cfg            *config.Config
	jwtService     JwtServiceInterface
	repoToken      repository.VerificationTokenRepositoryInterface
}

// GetCustomerByID implements [UserServiceInterface].
func (u *userService) GetCustomerByID(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	return u.userRepository.GetCustomerByID(ctx, userID)
}

// GetCustomerAll implements [UserServiceInterface].
func (u *userService) GetCustomerAll(ctx context.Context, query entity.QueryStringCustomer) ([]entity.UserEntity, int64, int64, error) {
	return u.userRepository.GetCustomerAll(ctx, query)
}

// UpdateDataUser implements [UserServiceInterface].
func (u *userService) UpdateDataUser(ctx context.Context, req entity.UserEntity) error {
	return u.userRepository.UpdateDataUser(ctx, req)
}

// GetProfileUser implements [UserServiceInterface].
func (u *userService) GetProfileUser(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	return u.userRepository.GetUserByID(ctx, userID)
}

// UpdatePassword implements [UserServiceInterface].
func (u *userService) UpdatePassword(ctx context.Context, req entity.UserEntity) error {
	token, err := u.repoToken.GetDataByToken(ctx, req.Token)
	if err != nil {
		log.Errorf("[UserService-1] UpdatePassword: %v", err)
		return err
	}
	if token.TokenType != "reset_password" {
		err = errors.New("invalid token")
		log.Errorf("[UserService-2] UpdatePassword: %v", err)
		return err
	}
	password, err := conv.HashPassword(req.Password)
	if err != nil {
		log.Errorf("[UserService-3] UpdatePassword: %v", err)
		return err
	}
	req.Password = password
	req.ID = token.UserID
	err = u.userRepository.UpdatePasswordByID(ctx, req)
	if err != nil {
		log.Errorf("[UserService-4] UpdatePassword: %v", err)
		return err
	}
	return nil
}

// VerifyToken implements [UserServiceInterface].
func (u *userService) VerifyToken(ctx context.Context, token string) (*entity.UserEntity, error) {
	verifyToken, err := u.repoToken.GetDataByToken(ctx, token)
	if err != nil {
		log.Errorf("[UserService-1] VerifyToken: %v", err)
		return nil, err
	}
	userVerify, err := u.userRepository.UpdateUserVerified(ctx, verifyToken.UserID)
	if err != nil {
		log.Errorf("[UserService-2] VerifyToken: %v", err)
		return nil, err
	}
	accessToken, err := u.jwtService.GenerateToken(userVerify.ID)
	if err != nil {
		log.Errorf("[UserService-3] VerifyToken: %v", err)
		return nil, err
	}

	sessionData := map[string]interface{}{
		"user_id":    userVerify.ID,
		"name":       userVerify.Name,
		"email":      userVerify.Email,
		"logged_in":  true,
		"created_at": time.Now().String(),
		"token":      accessToken,
		"role_name":  userVerify.RoleName,
	}

	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		log.Errorf("[UserService-4] VerifyToken Marshal: %v", err)
		return nil, err
	}

	redisConn := config.NewRedisClient()

	if err := redisConn.Set(
		ctx,
		accessToken,
		sessionJSON,
		22*time.Hour,
	).Err(); err != nil {
		log.Errorf("[UserService-5] VerifyToken Redis: %v", err)
		return nil, err
	}

	userVerify.Token = accessToken

	return userVerify, nil
}

// ForgotPassword implements [UserServiceInterface].
func (u *userService) ForgotPassword(ctx context.Context, req entity.UserEntity) error {
	user, err := u.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] ForgotPassword: %v", err)
		return err
	}
	token := uuid.New().String()
	req.Token = token
	err = u.repoToken.CreateVerificationToken(ctx, entity.VerificationToken{
		UserID:    user.ID,
		Token:     req.Token,
		TokenType: "reset_password",
		ExpiresAt: time.Now().Add(time.Hour * 1),
	})
	if err != nil {
		log.Errorf("[UserService-2] ForgotPassword: %v", err)
		return err
	}

	urlReset := fmt.Sprintf("%s/forgot-password?token=%s", u.cfg.App.UrlForgotPassword, req.Token)
	messageParam := fmt.Sprintf("Please reset your password with click link below : %v", urlReset)
	err = message.PublishMessage(req.Email, messageParam, "reset_password")
	if err != nil {
		log.Errorf("[UserService-3] ForgotPassword: %v", err)
		return err
	}

	return nil
}

// CreateUserAccount implements [UserServiceInterface].
func (u *userService) CreateUserAccount(ctx context.Context, req entity.UserEntity) error {
	password, err := conv.HashPassword(req.Password)
	if err != nil {
		log.Errorf("[UserService-1] CreateUserAccount: %v", err)
		return err
	}
	req.Password = password
	token := uuid.New().String()
	req.Token = token
	err = u.userRepository.CreateUserAccount(ctx, req)
	if err != nil {
		log.Errorf("[UserService-2] CreateUserAccount: %v", err)
		return err
	}

	urlVerify := fmt.Sprintf("http://localhost:8000/verify/%v", req.Token)
	messageParam := fmt.Sprintf("Please verify your account with click link below : %v", urlVerify)
	err = message.PublishMessage(req.Email, messageParam, "email_verification")
	if err != nil {
		log.Errorf("[UserService-3] CreateUserAccount: %v", err)
		return err
	}

	return nil
}

// SignIn implements [UserServiceInterface].
func (u *userService) SignIn(
	ctx context.Context,
	req entity.UserEntity,
) (*entity.UserEntity, string, error) {

	user, err := u.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] SignIn: %v", err)
		return nil, "", err
	}

	if user == nil {
		log.Warnf("[UserService-2] SignIn: user not found")
		return nil, "", errors.New("invalid email or password")
	}

	if !conv.CheckPasswordHash(req.Password, user.Password) {
		log.Warnf("[UserService-2] SignIn: invalid credentials")
		return nil, "", errors.New("invalid email or password")
	}

	token, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		log.Errorf("[UserService-3] SignIn: %v", err)
		return nil, "", err
	}

	sessionData := map[string]interface{}{
		"user_id":    user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"logged_in":  true,
		"created_at": time.Now().String(),
		"token":      token,
		"role_name":  user.RoleName,
	}

	redisConn := config.NewRedisClient()

	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		log.Errorf("[UserService-4] SignIn marshal session: %v", err)
		return nil, "", err
	}

	if err := redisConn.Set(
		ctx,
		token,
		sessionJSON,
		23*time.Hour,
	).Err(); err != nil {
		log.Errorf("[UserService-5] SignIn redis: %v", err)
		return nil, "", err
	}

	return user, token, nil
}

func NewUserService(userRepository repository.UserRepositoryInterface, cfg *config.Config, jwtService JwtServiceInterface, repoToken repository.VerificationTokenRepositoryInterface) UserServiceInterface {
	return &userService{
		userRepository: userRepository,
		cfg:            cfg,
		jwtService:     jwtService,
		repoToken:      repoToken,
	}
}
