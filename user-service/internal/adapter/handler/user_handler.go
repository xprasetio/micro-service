package handler

import (
	"encoding/json"
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
	ForgotPassword(ctx echo.Context) error
	VerifyAccount(ctx echo.Context) error
	UpdatePassword(ctx echo.Context) error
	GetProfileUser(ctx echo.Context) error
	UpdateDataUser(ctx echo.Context) error
}

type userHandler struct {
	userService service.UserServiceInterface
}

// UpdateDataUser implements [UserHandlerInterface].
func (u *userHandler) UpdateDataUser(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		req         = request.UpdateDataUserRequest{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[UpdateDataUser-1] UpdateDataUser: %s", "user not found")
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[UpdateDataUser-2] UpdateDataUser: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	userID := jwtUserData.UserID
	if err := c.Bind(&req); err != nil {
		log.Errorf("[UpdateDataUser-3] UpdateDataUser: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UpdateDataUser-4] UpdateDataUser: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	reqEntity := entity.UserEntity{
		ID:      userID,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Phone:   req.Phone,
		Photo:   req.Photo,
		Lat:     req.Lat,
		Lng:     req.Lng,
	}

	err = u.userService.UpdateDataUser(ctx, reqEntity)
	if err != nil {
		log.Errorf("[UpdateDataUser-4] UpdateDataUser: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	return c.JSON(http.StatusOK, resp)
}

// GetProfileUser implements [UserHandlerInterface].
func (u *userHandler) GetProfileUser(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		respProfile = response.ProfileResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)
	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[GetProfileUser-1] GetProfileUser: %s", "user not found")
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[GetProfileUser-2] GetProfileUser: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	userID := jwtUserData.UserID
	dataUser, err := u.userService.GetProfileUser(ctx, userID)
	if err != nil {
		log.Errorf("[GetProfileUser-3] GetProfileUser: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	resp.Message = "success"
	respProfile.Name = dataUser.Name
	respProfile.Email = dataUser.Email
	respProfile.ID = dataUser.ID
	respProfile.Address = dataUser.Address
	respProfile.Phone = dataUser.Phone
	respProfile.Photo = dataUser.Photo
	respProfile.RoleName = dataUser.RoleName
	respProfile.Lat = dataUser.Lat
	respProfile.Lng = dataUser.Lng
	resp.Data = respProfile

	return c.JSON(http.StatusOK, resp)

}

// VerifyAccount implements [UserHandlerInterface].
func (u *userHandler) VerifyAccount(c echo.Context) error {
	var (
		resp       = response.DefaultResponse{}
		respSignin = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	token := c.QueryParam("token")
	if token == "" {
		log.Errorf("[VerifyAccount-1] VerifyAccount: %s", "missing or invalid token")
		resp.Message = "missing or invalid token"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	user, err := u.userService.VerifyToken(ctx, token)
	if err != nil {
		if err.Error() == "user not found" {
			log.Errorf("[VerifyAccount-1] VerifyAccount: %s", "user not found")
			resp.Message = "user not found"
			resp.Data = nil
			return c.JSON(http.StatusUnauthorized, resp)
		}
		log.Errorf("[VerifyAccount-2] VerifyAccount: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	respSignin.ID = user.ID
	respSignin.Name = user.Name
	respSignin.Email = user.Email
	respSignin.AccessToken = user.Token
	resp.Message = "success"
	resp.Data = respSignin
	return c.JSON(http.StatusOK, resp)
}

// ForgotPassword implements [UserHandlerInterface].
func (u *userHandler) ForgotPassword(c echo.Context) error {
	var (
		req  = request.ForgotPasswordRequest{}
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)
	if err = c.Bind(&req); err != nil {
		log.Errorf("[ForgotPassword-1] ForgotPassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[ForgotPassword-2] ForgotPassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	reqEntity := entity.UserEntity{
		Email: req.Email,
	}
	err := u.userService.ForgotPassword(ctx, reqEntity)
	if err != nil {
		if err.Error() == "user not found" {
			log.Errorf("[UserHandler-1] ForgotPassword: %v", err)
			resp.Message = err.Error()
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		log.Errorf("[UserHandler-1] ForgotPassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
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

func (u *userHandler) UpdatePassword(c echo.Context) error {
	var (
		req  = request.UpdatePasswordRequest{}
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		log.Errorf("[UpdatePassword-1] UpdatePassword: %s", "missing or invalid token")
		resp.Message = "missing or invalid token"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	if err := c.Bind(&req); err != nil {
		log.Errorf("[UpdatePassword-2] UpdatePassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UpdatePassword-3] UpdatePassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	if req.NewPassword != req.NewPasswordConfirmation {
		log.Errorf("[UpdatePassword-4] UpdatePassword: %v", err)
		resp.Message = "password and password confirmation does not match"
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	reqEntity := entity.UserEntity{
		Password: req.NewPassword,
		Token:    tokenString,
	}
	err := u.userService.UpdatePassword(ctx, reqEntity)
	if err != nil {
		if err.Error() == "user not found" {
			log.Errorf("[UserHandler-1] UpdatePassword: %v", err)
			resp.Message = err.Error()
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		log.Errorf("[UserHandler-1] UpdatePassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Message = "Success Update Password"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func NewUserHandler(e *echo.Echo, userService service.UserServiceInterface, cfg *config.Config, jwtService service.JwtServiceInterface) UserHandlerInterface {
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
	e.POST("/forgot-password", userHandler.ForgotPassword)
	e.GET("/verify-account", userHandler.VerifyAccount)
	e.PUT("/update-password", userHandler.UpdatePassword)

	mid := adapter.NewMiddlewareAdapter(cfg, jwtService)
	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/profile", userHandler.GetProfileUser)
	adminGroup.GET("/check", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	authGroup := e.Group("/auth", mid.CheckToken())
	authGroup.PUT("/profile", userHandler.UpdateDataUser)

	return userHandler
}
