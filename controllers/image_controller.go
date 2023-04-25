package controllers

import (
	"fmt"
	"go-img/config"
	"go-img/services"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	JPEG = "jpeg"
	PNG  = "png"
	WEBP = "webp"
	AVIF = "avif"
)

type ImageController struct {
	config       config.Config
	imageService services.ImageService
}

func NewImageController(imageService services.ImageService, config config.Config) *ImageController {
	return &ImageController{
		imageService: imageService,
		config:       config,
	}
}

type QueryParams struct {
	Width   int
	Quality int
	Format  string
}

func (c *ImageController) HandleImageRequest(ctx *gin.Context) {
	imagePathParam := ctx.Param("pathToImage")
	imagePath, err := filepath.Abs(path.Join(c.config.StaticDir, imagePathParam))

	fmt.Printf("imagePath: %s\n", imagePath)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	params, err := c.parseHandlerParams(ctx.Request)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	image, err := c.imageService.LoadImage(imagePath)
	if err != nil {
		// todo: handle 404
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	imageBytes, err := c.imageService.ProcessImage(image, services.ImageOptions{
		Width:   params.Width,
		Quality: params.Quality,
		Format:  params.Format,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", fmt.Sprintf("image/%s", params.Format))
	ctx.Writer.Write(imageBytes)
}

func (c *ImageController) parseHandlerParams(request *http.Request) (QueryParams, error) {
	params := QueryParams{}
	query := request.URL.Query()

	var err error
	params.Width, err = c.parseAndValidateWidth(query.Get("w"))
	if err != nil {
		return QueryParams{}, fmt.Errorf("width: %w", err)
	}

	params.Quality, err = c.parseAndValidateQuality(query.Get("q"))
	if err != nil {
		return QueryParams{}, fmt.Errorf("quality: %w", err)
	}

	params.Format, err = c.determineImageFormat(query.Get("f"), request.Header.Get("Accept"))
	if err != nil {
		return QueryParams{}, fmt.Errorf("format: %w", err)
	}

	return params, nil
}

func (c *ImageController) parseAndValidateWidth(widthStr string) (int, error) {
	if widthStr == "" {
		return 0, nil
	}
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (c *ImageController) parseAndValidateQuality(qualityStr string) (int, error) {
	if qualityStr == "" {
		return c.config.DefaultQuality, nil
	}
	quality, err := strconv.Atoi(qualityStr)
	if err != nil {
		return 0, err
	}
	if quality < 1 || quality > 100 {
		return 0, fmt.Errorf("quality must be between 1 and 100")
	}
	return quality, nil
}

func (c *ImageController) determineImageFormat(formatParam, acceptHeader string) (string, error) {
	// Determine the format from the query parameter
	if formatParam != "" {
		format, err := c.determineImageFormatFromParam(formatParam)
		if err != nil {
			return "", err
		}
		return format, nil
	}

	// If no format specified in the query parameter, determine it from the Accept header
	format, err := c.determineImageFormatFromHeader(acceptHeader)
	if err == nil {
		return format, nil
	}

	// If no format specified in the query parameter and the Accept header, use a default format
	return c.config.DefaultImageFormat, nil
}

func (c *ImageController) determineImageFormatFromParam(format string) (string, error) {
	switch format {
	case JPEG, PNG, WEBP, AVIF:
		return format, nil
	case "":
		return "", fmt.Errorf("no image format specified")
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}
}

func (c *ImageController) determineImageFormatFromHeader(acceptHeader string) (string, error) {
	supportedFormats := map[string]string{
		"image/webp": WEBP,
		"image/avif": AVIF,
		"image/png":  PNG,
		"image/jpeg": JPEG,
		"image/apng": PNG,
	}

	for _, format := range strings.Split(acceptHeader, ",") {
		trimmedFormat := strings.TrimSpace(strings.Split(format, ";")[0])
		fmt.Printf("trimmed format: %s", trimmedFormat)
		if imageFormat, ok := supportedFormats[trimmedFormat]; ok {
			return imageFormat, nil
		}
	}

	return "", fmt.Errorf("no supported image format found in accept header")
}
