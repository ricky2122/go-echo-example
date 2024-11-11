package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ricky2122/go-echo-example/usecase"
)

type LoginRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type IAuthUseCase interface {
	Login(input usecase.LoginUseCaseInput) error
}

type AuthController struct {
	au IAuthUseCase
}

func NewAuthController(au IAuthUseCase) AuthController {
	return AuthController{au: au}
}

func (ac *AuthController) Login(c echo.Context) error {
	// parse request
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// validate
	if err := c.Validate(req); err != nil {
		return err
	}

	// login usecase
	input := usecase.LoginUseCaseInput{Name: req.Name, Password: req.Password}
	if err := ac.au.Login(input); err != nil {
		if errors.Is(err, usecase.ErrLoginFailed) {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed login")
		}
		return echo.NewHTTPError(http.StatusBadRequest, "internal server error")
	}

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
