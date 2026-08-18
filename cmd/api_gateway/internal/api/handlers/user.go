package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
