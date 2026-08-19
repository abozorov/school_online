package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetClassroomById godoc
// @Summary Get classroom by ID
// @Tags Classroom
// @Security BearerAuth
// @Produce json
// @Param id path int true "Classroom ID"
// @Success 200 {object} models.Classroom
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/classroom/{id} [get]
func (h *Handler) GetClassroomById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("classroom.GetClassroomById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// get by id
	classroom, err := h.classroomService.GetClassroomByID(c.Request.Context(), int32(id))
	if err != nil {
		h.logger.Error("classroom.GetClassroomById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, classroom)
}

// ListClassrooms godoc
// @Summary List classrooms
// @Tags Classroom
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Classroom
// @Failure 401 {object} map[string]string
// @Router /api/classroom/list [get]
func (h *Handler) ListClassrooms(c *gin.Context) {
	// get list
	classrooms, err := h.classroomService.List(c.Request.Context())
	if err != nil {
		h.logger.Error("classroom.ListClassrooms: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, classrooms)
}

// CreateClassroom godoc
// @Summary Create classroom
// @Tags Classroom
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ClassroomRequest true "Classroom data"
// @Success 200 {object} map[string]int32
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/classroom [post]
func (h *Handler) CreateClassroom(c *gin.Context) {
	var req models.ClassroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("classroom.CreateClassroom: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// create classroom
	id, err := h.classroomService.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("classroom.CreateClassroom: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// UpdateClassroomById godoc
// @Summary Update classroom
// @Tags Classroom
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.Classroom true "Updated classroom data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/classroom [patch]
func (h *Handler) UpdateClassroomById(c *gin.Context) {
	var req models.Classroom
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("classroom.UpdateClassroomById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// update classroom
	err := h.classroomService.UpdateByID(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("classroom.UpdateClassroomById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"message": "classroom updated"})
}

// DeleteClassroomById godoc
// @Summary Delete classroom by ID
// @Tags Classroom
// @Security BearerAuth
// @Produce json
// @Param id path int true "Classroom ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/classroom/{id} [delete]
func (h *Handler) DeleteClassroomById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("classroom.DeleteClassroomById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// delete classroom
	err = h.classroomService.DeleteByID(c.Request.Context(), int32(id))
	if err != nil {
		h.logger.Error("classroom.DeleteClassroomById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"message": "classroom deleted"})
}