package handler

import (
	"net/http"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type UserHandlerInterface interface {
	SignIn(ctx echo.Context) error
	CreateUserAccount(ctx echo.Context) error
}

type userHandler struct {
	userService service.UserServiceInterface
}

// CreateUserAccount implements [UserHandlerInterface].
func (u *userHandler) CreateUserAccount(c echo.Context) error {
	var (
		req  = request.SignUpRequest{}
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)
	if err := c.Bind(&req); err != nil {
		log.Errorf("[CreateUserAccount-1] CreateUserAccount: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	if err = c.Validate(req); err != nil {
		log.Errorf("[CreateUserAccount-2] CreateUserAccount: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	if req.Password != req.PasswordConfirmation {
		log.Errorf("[CreateUserAccount-3] CreateUserAccount: %v", err)
		resp.Message = "password and password confirmation does not match"
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err := u.userService.CreateUserAccount(ctx, reqEntity)
	if err != nil {
		log.Errorf("[UserHandler-1] CreateUserAccount: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

var err error

// SignIn implements [UserHandlerInterface].
func (u *userHandler) SignIn(c echo.Context) error {
	var (
		req        = request.SignInRequest{}
		respSignIn = response.SignInResponse{}
		resp       = response.DefaultResponse{}
		ctx        = c.Request().Context()
	)
	if err := c.Bind(&req); err != nil {
		log.Errorf("[SignIn-1] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[SignIn-2] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := u.userService.SignIn(ctx, reqEntity)
	if err != nil {
		if err.Error() == "user not found" {
			log.Errorf("[UserHandler-1] SignIn: %v", err)
			resp.Message = err.Error()
			resp.Data = nil
			return c.JSON(http.StatusUnprocessableEntity, resp)
		}
		log.Errorf("[UserHandler-2] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	respSignIn.ID = user.ID
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.RoleName = user.RoleName
	respSignIn.Address = user.Address
	respSignIn.Phone = user.Phone
	respSignIn.Photo = user.Photo
	respSignIn.Lat = user.Lat
	respSignIn.Lng = user.Lng
	respSignIn.IsVerified = user.IsVerified
	respSignIn.AccessToken = token

	if err = c.Validate(respSignIn); err != nil {
		if err.Error() == "user not found" {
			log.Errorf("[UserHandler-2] SignIn: %v", err)
			resp.Message = err.Error()
			resp.Data = nil
			return c.JSON(http.StatusUnprocessableEntity, resp)
		}
		log.Errorf("[UserHandler-3] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	resp.Message = "Success"
	resp.Data = respSignIn
	return c.JSON(200, resp)
}

func NewUserHandler(e *echo.Echo, userService service.UserServiceInterface, cfg *config.Config) UserHandlerInterface {
	if e == nil {
		panic("echo is nil")
	}

	if userService == nil {
		panic("userService is nil")
	}

	if cfg == nil {
		panic("config is nil")
	}

	userHandler := &userHandler{userService: userService}
	e.Use(middleware.Recover())
	e.POST("/signin", userHandler.SignIn)
	e.POST("/signup", userHandler.CreateUserAccount)

	mid := adapter.NewMiddlewareAdapter(cfg)
	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/check", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	return userHandler
}
