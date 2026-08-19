package handlers

import (
	"net/http"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
