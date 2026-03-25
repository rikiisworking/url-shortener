package handler

import (
	"gitgub.com/rikiisworking/url-shortener/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type URLHandler struct {
	service   *service.URLService
	validator *validator.Validate
}

func NewURLHandler(svc *service.URLService) *URLHandler {
	return &URLHandler{
		service:   svc,
		validator: validator.New(),
	}
}

type ShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func (h *URLHandler) Shorten(c fiber.Ctx) error {
	var req ShortenRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	shortCode, err := h.service.Shorten(c.Context(), req.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	shortURLHost := "http://" + c.Hostname()
	port, ok := c.Locals("port").(string)
	if ok && port != "80" {
		shortURLHost += ":" + port
	}

	shortURL := shortURLHost + "/" + shortCode

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"short_code":   shortCode,
		"short_url":    shortURL,
		"original_url": req.URL,
	})
}

func (h *URLHandler) Redirect(c fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Short code required")
	}

	originalURL, err := h.service.GetOriginalURL(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("URL not found")
	}

	return c.Redirect().To(originalURL)
}
