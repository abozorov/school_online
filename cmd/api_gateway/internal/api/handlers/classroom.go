package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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