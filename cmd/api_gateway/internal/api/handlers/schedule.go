package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetScheduleByClassroomId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("schedule.GetScheduleByClassroomId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// get by id
	schedule, err := h.classroomService.GetScheduleByClassroomId(c.Request.Context(), int32(id))
	if err != nil {
		h.logger.Error("schedule.GetScheduleByClassroomId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, schedule)
}

func (h *Handler) GetScheduleByTeacherId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("schedule.GetScheduleByTeacherId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// get by id
	schedule, err := h.classroomService.GetScheduleByTeacherId(c.Request.Context(), int32(id))
	if err != nil {
		h.logger.Error("schedule.GetScheduleByTeacherId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, schedule)
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	var req models.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("schedule.CreateSchedule: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// create schedule
	id, err := h.classroomService.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("schedule.CreateSchedule: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) UpdateScheduleById(c *gin.Context) {
	var req models.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("schedule.UpdateScheduleById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// update schedule
	err := h.classroomService.UpdateScheduleById(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("schedule.UpdateScheduleById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}

func (h *Handler) DeleteScheduleById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("schedule.DeleteScheduleById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// delete schedule
	err = h.classroomService.DeleteScheduleById(c.Request.Context(), int32(id))
	if err != nil {
		h.logger.Error("schedule.DeleteScheduleById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}