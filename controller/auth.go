package controller

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ricky2122/go-echo-example/usecase"
)

const SessionKey = "session_id"

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
	// Parse request
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// Validate
	if err := c.Validate(req); err != nil {
		return err
	}

	// Login usecase
	input := usecase.LoginUseCaseInput{Name: req.Name, Password: req.Password}
	if err := ac.au.Login(input); err != nil {
		if errors.Is(err, usecase.ErrLoginFailed) {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed login")
		}
		return echo.NewHTTPError(http.StatusBadRequest, "internal server error")
	}

	// Set session_id to cookie
	sess, err := session.Get(SessionKey, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		// Secure:   true,
	}
	sess.Values[SessionKey] = input.Name
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// Send response
	return c.NoContent(http.StatusOK)
}

func (ac *AuthController) Logout(c echo.Context) error {
	// delete session
	sess, err := session.Get(SessionKey, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// send response
	return c.NoContent(http.StatusNoContent)
}
