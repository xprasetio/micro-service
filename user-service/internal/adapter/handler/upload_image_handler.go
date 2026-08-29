package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/adapter/storage"
	"user-service/internal/core/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UploadImageInterface interface {
	UploadImage(ctx echo.Context) error
}

type uploadImage struct {
	storageHandler storage.SupabaseInterface
}

// UploadImage implements [UploadImageInterface].
func (u *uploadImage) UploadImage(ctx echo.Context) error {

	var resp = response.DefaultResponse{}
	file, err := ctx.FormFile("photo")
	if err != nil {
		log.Errorf("[UploadImage-1] UploadImage: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return ctx.JSON(http.StatusUnprocessableEntity, resp)
	}
	src, err := file.Open()
	if err != nil {
		log.Errorf("[UploadImage-2] UploadImage: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return ctx.JSON(http.StatusUnprocessableEntity, resp)
	}
	defer src.Close()
	fileBuffer := new(bytes.Buffer)
	_, err = io.Copy(fileBuffer, src)
	if err != nil {
		log.Errorf("[UploadImage-3] UploadImage: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return ctx.JSON(http.StatusUnprocessableEntity, resp)
	}
	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), getExtension(file.Filename))

	uploadPath := fmt.Sprintf("public/uploads/%s", newFileName)

	url, err := u.storageHandler.UploadFile(uploadPath, fileBuffer)
	if err != nil {
		log.Errorf("[UploadImage-4] UploadImage: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return ctx.JSON(http.StatusUnprocessableEntity, resp)
	}
	resp.Message = "success"
	resp.Data = map[string]string{"image_url": url}
	return ctx.JSON(http.StatusOK, resp)
}

func getExtension(fileName string) string {
	ext := "." + fileName[strings.LastIndex(fileName, ".")+1:]
	return ext
}

func NewUploadImageHandler(e *echo.Echo, cfg *config.Config, storageHandler storage.SupabaseInterface, jwtService service.JwtServiceInterface) UploadImageInterface {
	res := &uploadImage{
		storageHandler: storageHandler,
	}
	mid := adapter.NewMiddlewareAdapter(cfg, jwtService)
	e.POST("auth/profile/image-upload", res.UploadImage, mid.CheckToken())
	return res
}
