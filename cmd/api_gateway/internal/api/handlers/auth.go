package handlers

import (
	"net/http"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) Login(c *gin.Context) {
	// get email & password
	var req models.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.logger.Error("auth.Login: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// login
	tokens, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("auth.Login: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// answer
	c.JSON(http.StatusOK, tokens)
}
