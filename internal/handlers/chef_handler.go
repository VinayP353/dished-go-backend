package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/dished-go-backend/internal/models"
	"github.com/yourusername/dished-go-backend/internal/service"
)

type ChefHandler struct {
	service service.ChefService
}

func NewChefHandler(service service.ChefService) *ChefHandler {
	return &ChefHandler{service: service}
}

// CreateChef godoc
// @Summary      Create a new chef
// @Tags         chefs
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateChefRequest  true  "Create chef payload"
// @Success      201   {object}  models.Chef
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /chefs [post]
func (h *ChefHandler) CreateChef(c *gin.Context) {
	// Step 1: Parse and validate the JSON body
	var req models.CreateChefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Step 2: Call the service
	chef, err := h.service.CreateChef(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Step 3: Return the response
	c.JSON(http.StatusCreated, gin.H{"message": "chef created successfully", "chef": chef})
}

// Register godoc
// @Summary      Register a new chef
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "Register payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /auth/register [post]
func (h *ChefHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chef, err := h.service.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registration successful", "chef": chef})
}

// Login godoc
// @Summary      Login as a chef
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Login payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *ChefHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chef, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "chef": chef})
}

// GetUsernames godoc
// @Summary      Get all chef usernames
// @Tags         chefs
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /chefs/usernames [get]
func (h *ChefHandler) GetUsernames(c *gin.Context) {
	usernames, err := h.service.GetAllUsernames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := make([]models.UsernameResponse, len(usernames))
	for i, u := range usernames {
		response[i] = models.UsernameResponse{Username: u}
	}
	c.JSON(http.StatusOK, gin.H{"usernames": response, "count": len(usernames)})
}

// GetChef godoc
// @Summary      Get a chef by ID
// @Tags         chefs
// @Produce      json
// @Param        id   path      int  true  "Chef ID"
// @Success      200  {object}  models.Chef
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /chefs/{id} [get]
func (h *ChefHandler) GetChef(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chef id"})
		return
	}
	chef, err := h.service.GetChef(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chef not found"})
		return
	}
	c.JSON(http.StatusOK, chef)
}

// GetAllChefs godoc
// @Summary      Get all chefs
// @Tags         chefs
// @Produce      json
// @Success      200  {array}   models.Chef
// @Failure      500  {object}  map[string]interface{}
// @Router       /chefs [get]
func (h *ChefHandler) GetAllChefs(c *gin.Context) {
	chefs, err := h.service.GetAllChefs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chefs)
}

// UpdateChef godoc
// @Summary      Update a chef
// @Tags         chefs
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "Chef ID"
// @Param        body  body      models.UpdateChefRequest  true  "Update payload"
// @Success      200   {object}  models.Chef
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /chefs/{id} [put]
func (h *ChefHandler) UpdateChef(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chef id"})
		return
	}
	var req models.UpdateChefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chef, err := h.service.UpdateChef(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chef)
}

// UpdateProfile godoc
// @Summary      Update a chef's profile
// @Tags         chefs
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Chef ID"
// @Param        body  body      models.UpdateChefProfileRequest  true  "Profile update payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /chefs/{id}/profile [put]
func (h *ChefHandler) UpdateProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chef id"})
		return
	}
	var req models.UpdateChefProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chef, err := h.service.UpdateProfile(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully", "chef": chef})
}

// DeleteChef godoc
// @Summary      Delete a chef
// @Tags         chefs
// @Produce      json
// @Param        id   path      int  true  "Chef ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /chefs/{id} [delete]
func (h *ChefHandler) DeleteChef(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chef id"})
		return
	}
	if err := h.service.DeleteChef(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "chef deleted successfully"})
}
