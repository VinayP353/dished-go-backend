package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/dished-go-backend/internal/models"
	"github.com/yourusername/dished-go-backend/internal/service"
)

type CookHandler struct {
	service service.CookService
}

func NewCookHandler(service service.CookService) *CookHandler {
	return &CookHandler{service: service}
}

func (h *CookHandler) CreateCook(c *gin.Context) {
	var req models.CreateCookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cook, err := h.service.CreateCook(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cook)
}

func (h *CookHandler) GetCook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cook id"})
		return
	}

	cook, err := h.service.GetCook(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cook not found"})
		return
	}

	c.JSON(http.StatusOK, cook)
}

func (h *CookHandler) GetAllCooks(c *gin.Context) {
	cooks, err := h.service.GetAllCooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cooks)
}

func (h *CookHandler) UpdateCook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cook id"})
		return
	}

	var req models.UpdateCookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cook, err := h.service.UpdateCook(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cook)
}

func (h *CookHandler) DeleteCook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cook id"})
		return
	}

	if err := h.service.DeleteCook(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cook deleted successfully"})
}

func (h *CookHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cook, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cook)
}
