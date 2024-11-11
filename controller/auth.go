package controller

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

	// set session_id to cookie
	sess, err := session.Get("session_id", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
	}
	sess.Values[input.Name] = input.Name
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// send response
	return c.NoContent(http.StatusOK)
}

func (ac *AuthController) Logout(c echo.Context) error {
	// Todo: get session_id from cookie

	// Todo: logout usecase

	// send response
	return c.NoContent(http.StatusNoContent)
}
