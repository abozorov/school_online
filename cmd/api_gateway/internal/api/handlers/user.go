package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetById godoc
// @Summary Get user by ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/user/{id} [get]
func (h *Handler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("user.GetById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// get by id
	usr, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("user.GetById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, usr)
}

// List godoc
// @Summary List users
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} map[string]string
// @Router /api/user/list [get]
func (h *Handler) List(c *gin.Context) {
	// get list
	users, err := h.userService.List(c.Request.Context())
	if err != nil {
		h.logger.Error("user.List: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, users)
}

// Create godoc
// @Summary Create user
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.RegisterUserRequest true "User data"
// @Success 201 {object} map[string]int32
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/user [post]
func (h *Handler) Create(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("user.Create: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// create user
	id, err := h.userService.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("user.Create: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateById godoc
// @Summary Update user
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateUserRequest true "Updated user data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/user [patch]
func (h *Handler) UpdateById(c *gin.Context) {

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("user.UpdateById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// update user
	err := h.userService.UpdateByID(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("user.UpdateById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

// DeleteById godoc
// @Summary Delete user by ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/user/{id} [delete]
func (h *Handler) DeleteById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("user.DeleteById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// delete user
	err = h.userService.DeleteByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("user.DeleteById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
