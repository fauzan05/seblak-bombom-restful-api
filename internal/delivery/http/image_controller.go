package http

import (
	"fmt"
	"os"
	"path/filepath"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ImageController struct {
	Log     *logrus.Logger
	UseCase *usecase.ImageUseCase
}

func NewImageController(useCase *usecase.ImageUseCase, logger *logrus.Logger) *ImageController {
	return &ImageController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ImageController) Creates(ctx *fiber.Ctx) error {
	// Buat direktori uploads jika belum ada
	if _, err := os.Stat("../uploads/images/products"); os.IsNotExist(err) {
		os.Mkdir("../uploads/images/products", os.ModePerm)
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		c.Log.Warnf("Failed to get image data : %+v", err)
		return err
	}

	// Validasi ekstensi file
	allowedExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	}

	ext := filepath.Ext(file.Filename)
		if !allowedExtensions[ext] {
			c.Log.Warnf("Failed to get image data : %+v", err)
			return err
		}

	// Simpan file di direktori uploads
	filePath := fmt.Sprintf("../uploads/images/products/%s", file.Filename)
	if err := ctx.SaveFile(file, filePath); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*[]model.ImageResponse]{
		Code:   201,
		Status: "Success to add images",
	})
	
}

// func (c *ImageController) Creates(ctx *fiber.Ctx) error {
// 	request := new(model.AddImagesRequest)
// 	if err := ctx.BodyParser(&request.Images); err != nil {
// 		c.Log.Warnf("Cannot parse data : %+v", err)
// 		return err
// 	}

// 	response, err := c.UseCase.Add(ctx.Context(), request)
// 	if err != nil {
// 		c.Log.Warnf("Failed to add images : %+v", err)
// 		return err
// 	}
// 	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*[]model.ImageResponse]{
// 		Code:   201,
// 		Status: "Success to add images",
// 		Data:   response,
// 	})
// }

func (c *ImageController) EditPosition(ctx *fiber.Ctx) error {
	request := new(model.UpdateImagesRequest)
	if err := ctx.BodyParser(&request.Images); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Update(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to edit images positions : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*[]model.ImageResponse]{
		Code:   200,
		Status: "Success update images position",
		Data:   response,
	})
}

func (c *ImageController) Remove(ctx *fiber.Ctx) error {
	getId := ctx.Params("imageId")
	imageId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert image id : %+v", err)
		return err
	}

	imageRequest := new(model.DeleteImageRequest)
	imageRequest.ID = uint64(imageId)

	response, err := c.UseCase.Delete(ctx.Context(), imageRequest)
	if err != nil {
		c.Log.Warnf("Failed to delete image by id : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success delete image by id",
		Data:   response,
	})
}
