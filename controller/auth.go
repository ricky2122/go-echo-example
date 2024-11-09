package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthController struct{}

func NewAuthController() AuthController {
	return AuthController{}
}

func (ac *AuthController) Login(c echo.Context) error {
	// parse request
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// validate
	if err := c.Validate(req); err != nil {
		c.Logger().Info(err)
		return err
	}

	// Todo: login usecase

	// Todo: set session_id to cookie

	// send response
	return c.NoContent(http.StatusOK)
}

func (ac *AuthController) Logout(c echo.Context) error {
	// Todo: get session_id from cookie

	// Todo: logout usecase

	// send response
	return c.NoContent(http.StatusNoContent)
}
