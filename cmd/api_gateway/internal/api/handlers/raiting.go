package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func validateDateRange(value string) error {
	const layout = "02.01.2006"

	var firstDateStr, secondDateStr string

	if len(value) < 21 {
		return errs.ErrBadRequest
	}

	if len(value) != 21 {
		return errs.ErrBadRequest
	}

	if value[10] != '-' {
		return errs.ErrBadRequest
	}

	firstDateStr = value[:10]
	secondDateStr = value[11:]

	firstDate, err := time.Parse(layout, firstDateStr)
	if err != nil {
		return errs.ErrBadRequest
	}

	secondDate, err := time.Parse(layout, secondDateStr)
	if err != nil {
		return errs.ErrBadRequest
	}

	if !secondDate.After(firstDate) {
		return errs.ErrBadRequest
	}

	today := time.Now()

	today = time.Date(
		today.Year(),
		today.Month(),
		today.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	if secondDate.After(today) {
		return errs.ErrBadRequest
	}

	return nil
}

// GetJournalByStudentId godoc
// @Summary Get journal by student ID
// @Tags Raiting
// @Security BearerAuth
// @Produce json
// @Param id path int true "Student ID"
// @Param request body map[string]string true "Date range in format DD.MM.YYYY-DD.MM.YYYY"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/journal/student/{id} [post]
func (h *Handler) GetJournalByStudentId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("raiting.GetJournalByStudentId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	dataRange := make(map[string]string)
	err = c.ShouldBindJSON(&dataRange)
	if err != nil {
		h.logger.Error("raiting.GetJournalByStudentId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	dateRange, ok := dataRange["date_range"]
	if !ok {
		h.logger.Error("raiting.GetJournalByStudentId: ", zap.String("error", "date_range query param is required"))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	err = validateDateRange(dateRange)
	if err != nil {
		h.logger.Error("raiting.GetJournalByStudentId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// get journal by student id
	journal, err := h.raitingService.GetJournalByStudentId(c.Request.Context(), id, dateRange)
	if err != nil {
		h.logger.Error("raiting.GetJournalByStudentId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, journal)
}

// GetJournalByClassroomId godoc
// @Summary Get journal by classroom ID
// @Tags Raiting
// @Security BearerAuth
// @Produce json
// @Param id path int true "Classroom ID"
// @Param request body map[string]string true "Date range(date_range) in format DD.MM.YYYY-DD.MM.YYYY"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/journal/classroom/{id} [post]
func (h *Handler) GetJournalByClassroomId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("raiting.GetJournalByClassroomId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	dataRange := make(map[string]string)
	err = c.ShouldBindJSON(&dataRange)
	if err != nil {
		h.logger.Error("raiting.GetJournalByClassroomId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	dateRange, ok := dataRange["date_range"]
	if !ok {
		h.logger.Error("raiting.GetJournalByClassroomId: ", zap.String("error", "date_range query param is required"))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	err = validateDateRange(dateRange)
	if err != nil {
		h.logger.Error("raiting.GetJournalByClassroomId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	// get journal by classroom id
	journal, err := h.raitingService.GetJournalByClassroomId(c.Request.Context(), id, dateRange)
	if err != nil {
		h.logger.Error("raiting.GetJournalByClassroomId: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// write response
	c.JSON(http.StatusOK, journal)
}

// UpdateJournal godoc
// @Summary Update journal entry
// @Tags Raiting
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.Journal true "Journal update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/journal [patch]
func (h *Handler) UpdateJournal(c *gin.Context) {
	var req models.Journal

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("raiting.UpdateJournal: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequest)
		return
	}

	err := h.raitingService.UpdateJournal(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("raiting.UpdateJournal: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "journal updated successfully"})
}
