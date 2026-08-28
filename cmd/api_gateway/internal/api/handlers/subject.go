package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetSubjectById godoc
// @Summary Get subject by ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param id path int true "Subject ID"
// @Success 200 {object} models.Subject
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/user/subject/{id} [get]
func (h *Handler) GetSubjectById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("user.GetSubjectById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	subject, err := h.userService.GetSubjectById(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("user.GetSubjectById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.JSON(http.StatusOK, subject)
}

// ListSubjects godoc
// @Summary List subjects
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Subject
// @Failure 401 {object} map[string]string
// @Router /api/user/subject/list [get]
func (h *Handler) ListSubjects(c *gin.Context) {
	subjects, err := h.userService.ListSubjects(c.Request.Context())
	if err != nil {
		h.logger.Error("user.ListSubjects: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// UpdateSubject godoc
// @Summary Update subject
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateSubjectRequest true "Updated subject data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/user/subject [patch]
func (h *Handler) UpdateSubject(c *gin.Context) {
	var req models.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("user.UpdateSubject: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	if err := h.userService.UpdateSubject(c.Request.Context(), req); err != nil {
		h.logger.Error("user.UpdateSubject: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subject updated"})
}

// CreateSubject godoc
// @Summary Create subject
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateSubjectRequest true "Subject data"
// @Success 201 {object} map[string]int32
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/user/subject [post]
func (h *Handler) CreateSubject(c *gin.Context) {
	var req models.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("user.CreateSubject: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// create subject
	id, err := h.userService.CreateSubject(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("user.CreateSubject: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
