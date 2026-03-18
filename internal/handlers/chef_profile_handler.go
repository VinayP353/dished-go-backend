package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/dished-go-backend/internal/models"
	"github.com/yourusername/dished-go-backend/internal/service"
)

type ChefProfileHandler struct {
	service service.ChefProfileService
}

func NewChefProfileHandler(service service.ChefProfileService) *ChefProfileHandler {
	return &ChefProfileHandler{service: service}
}

// CreateProfile godoc
// @Summary      Create a chef profile
// @Tags         chef-profiles
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateChefProfileRequest  true  "Profile payload"
// @Success      201   {object}  models.ChefProfile
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /chef-profiles [post]
func (h *ChefProfileHandler) CreateProfile(c *gin.Context) {
	var req models.CreateChefProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile, err := h.service.CreateProfile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, profile)
}

// GetProfile godoc
// @Summary      Get a chef profile by ID
// @Tags         chef-profiles
// @Produce      json
// @Param        id   path      int  true  "Profile ID"
// @Success      200  {object}  models.ChefProfile
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /chef-profiles/{id} [get]
func (h *ChefProfileHandler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile id"})
		return
	}
	profile, err := h.service.GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// GetAllProfiles godoc
// @Summary      Get all chef profiles
// @Tags         chef-profiles
// @Produce      json
// @Success      200  {array}   models.ChefProfile
// @Failure      500  {object}  map[string]interface{}
// @Router       /chef-profiles [get]
func (h *ChefProfileHandler) GetAllProfiles(c *gin.Context) {
	profiles, err := h.service.GetAllProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

// UpdateProfile godoc
// @Summary      Update a chef profile
// @Tags         chef-profiles
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Profile ID"
// @Param        body  body      models.UpdateChefProfileRequest  true  "Update payload"
// @Success      200   {object}  models.ChefProfile
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /chef-profiles/{id} [put]
func (h *ChefProfileHandler) UpdateProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile id"})
		return
	}
	var req models.UpdateChefProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile, err := h.service.UpdateProfile(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// DeleteProfile godoc
// @Summary      Delete a chef profile
// @Tags         chef-profiles
// @Produce      json
// @Param        id   path      int  true  "Profile ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /chef-profiles/{id} [delete]
func (h *ChefProfileHandler) DeleteProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile id"})
		return
	}
	if err := h.service.DeleteProfile(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile deleted successfully"})
}
