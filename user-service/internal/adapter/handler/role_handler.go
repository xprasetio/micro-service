package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

type RoleHandlerInterface interface {
	Create(ctx echo.Context) error
	GetAll(ctx echo.Context) error
	GetByID(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type roleHandler struct {
	roleService service.RoleServiceInterface
}

// CreateRole implements [RoleHandlerInterface].
func (r *roleHandler) Create(c echo.Context) error {
	var (
		req  = request.RoleRequest{}
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)
	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[CreateRole-0] CreateRole: %s", "user not found")
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	jwtUserData := entity.JwtUserData{}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[CreateRole-0] CreateRole: %v", err)
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	if jwtUserData.RoleName != "admin" {
		log.Errorf("[CreateRole-0] CreateRole: %s", "user is not admin")
		resp.Message = "user is not admin"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	if err := c.Bind(&req); err != nil {
		log.Errorf("[CreateRole-1] CreateRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err := c.Validate(req); err != nil {
		log.Errorf("[CreateRole-2] CreateRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	reqEntity := entity.RoleEntity{
		Name: req.Name,
	}
	err = r.roleService.Create(ctx, reqEntity)
	if err != nil {
		log.Errorf("[CreateRole-3] CreateRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

// DeleteRole implements [RoleHandlerInterface].
func (r *roleHandler) Delete(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)
	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[DeleteRole-0] DeleteRole: %s", "user not found")
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	jwtUserData := entity.JwtUserData{}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[DeleteRole-0] DeleteRole: %v", err)
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	if jwtUserData.RoleName != "admin" {
		log.Errorf("[DeleteRole-0] DeleteRole: %s", "user is not admin")
		resp.Message = "user is not admin"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	roleIDString := c.Param("id")

	if err != nil {
		log.Errorf("[DeleteRole-1] DeleteRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	roleID, err := strconv.Atoi(roleIDString)
	if err != nil {
		log.Errorf("[DeleteRole-2] DeleteRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	err = r.roleService.Delete(ctx, int64(roleID))
	if err != nil {
		log.Errorf("[DeleteRole-2] DeleteRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

// GetAllRole implements [RoleHandlerInterface].
func (r *roleHandler) GetAll(c echo.Context) error {
	var (
		resp     = response.DefaultResponse{}
		respRole = []response.RoleResponse{}
		ctx      = c.Request().Context()
	)
	// user := c.Get("user").(string)
	// if user == "" {
	// 	log.Errorf("[GetAllRole-0] GetAllRole: %s", "user not found")
	// 	resp.Message = "user not found"
	// 	resp.Data = nil
	// 	return c.JSON(http.StatusUnauthorized, resp)
	// }
	// jwtUserData := entity.JwtUserData{}
	// err := json.Unmarshal([]byte(user), &jwtUserData)
	// if err != nil {
	// 	log.Errorf("[GetAllRole-0] GetAllRole: %v", err)
	// 	resp.Message = "user not found"
	// 	resp.Data = nil
	// 	return c.JSON(http.StatusUnauthorized, resp)
	// }
	// if jwtUserData.RoleName != "admin" {
	// 	log.Errorf("[GetAllRole-0] GetAllRole: %s", "user is not admin")
	// 	resp.Message = "user is not admin"
	// 	resp.Data = nil
	// 	return c.JSON(http.StatusUnauthorized, resp)
	// }
	search := c.QueryParam("search")
	if err != nil {
		log.Errorf("[GetAllRole-1] GetAllRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	roles, err := r.roleService.GetAll(ctx, search)
	if err != nil {
		log.Errorf("[GetAllRole-2] GetAllRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	for _, role := range roles {
		respRole = append(respRole, response.RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	resp.Message = "success"
	resp.Data = respRole
	return c.JSON(http.StatusOK, resp)
}

// GetRole implements [RoleHandlerInterface].
func (r *roleHandler) GetByID(c echo.Context) error {
	var (
		resp     = response.DefaultResponse{}
		respRole = response.RoleResponse{}
		ctx      = c.Request().Context()
	)
	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[GetRole-0] GetRole: %s", "user not found")
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	jwtUserData := entity.JwtUserData{}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[GetRole-0] GetRole: %v", err)
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	if jwtUserData.RoleName != "admin" {
		log.Errorf("[GetRole-0] GetRole: %s", "user is not admin")
		resp.Message = "user is not admin"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	roleIDString := c.Param("id")
	roleID, err := strconv.Atoi(roleIDString)
	if err != nil {
		log.Errorf("[GetRole-1] GetRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	role, err := r.roleService.GetByID(ctx, int64(roleID))
	if err != nil {
		log.Errorf("[GetRole-2] GetRole: %v", err)
		if err.Error() == "role not found" {
			resp.Message = err.Error()
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	respRole.ID = role.ID
	respRole.Name = role.Name
	resp.Message = "success"
	resp.Data = respRole
	return c.JSON(http.StatusOK, resp)
}

// UpdateRole implements [RoleHandlerInterface].
func (r *roleHandler) Update(c echo.Context) error {
	var (
		req         = request.RoleRequest{}
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)
	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[UpdateRole-0] UpdateRole: %s", "user not found")
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[UpdateRole-0] UpdateRole: %v", err)
		resp.Message = "user not found"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	if jwtUserData.RoleName != "admin" {
		log.Errorf("[UpdateRole-0] UpdateRole: %s", "user is not admin")
		resp.Message = "user is not admin"
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}
	roleIDString := c.Param("id")
	if roleIDString == "" {
		log.Errorf("[UpdateRole-1] UpdateRole: %s", "role id not found")
		resp.Message = "role id not found"
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	roleID, err := strconv.Atoi(roleIDString)
	if err != nil {
		log.Errorf("[UpdateRole-2] UpdateRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	if err := c.Bind(&req); err != nil {
		log.Errorf("[UpdateRole-3] UpdateRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	if err := c.Validate(req); err != nil {
		log.Errorf("[UpdateRole-4] UpdateRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	reqEntity := entity.RoleEntity{
		ID:   int64(roleID),
		Name: req.Name,
	}
	err = r.roleService.Update(ctx, reqEntity)
	if err != nil {
		log.Errorf("[UpdateRole-5] UpdateRole: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func NewRoleHandler(e *echo.Echo, roleService service.RoleServiceInterface, cfg *config.Config, jwtService service.JwtServiceInterface) RoleHandlerInterface {
	roleHandler := &roleHandler{
		roleService: roleService,
	}
	e.Use(middleware.Recover())
	mid := adapter.NewMiddlewareAdapter(cfg, jwtService)
	admintGroup := e.Group("/admin", mid.CheckToken())
	admintGroup.POST("/roles", roleHandler.Create)
	admintGroup.GET("/roles", roleHandler.GetAll)
	admintGroup.GET("/roles/:id", roleHandler.GetByID)
	admintGroup.PUT("/roles/:id", roleHandler.Update)
	admintGroup.DELETE("/roles/:id", roleHandler.Delete)
	return roleHandler
}
