package handlers

import (
	"net/http"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Login godoc
// @Summary Login user
// @Tags Auth
// @Description Authenticate user and return JWT and refresh tokens.
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
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
