package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ricky2122/go-echo-example/domain"
	"github.com/ricky2122/go-echo-example/usecase"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	BirthDay string `json:"birth_day" validate:"required"`
}

type SignUpResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetUserRequest struct {
	ID int `param:"id" validate:"gte=1"`
}

type GetUserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	BirthDay string `json:"birth_day"`
}

type IUserUseCase interface {
	SignUp(usecase.SignUpUseCaseInput) (*usecase.SignUpUseCaseOutput, error)
}

type UserController struct {
	uuc IUserUseCase
}

func NewUserController(uuc IUserUseCase) UserController {
	return UserController{uuc: uuc}
}

func (uc *UserController) SignUp(c echo.Context) error {
	// parse request
	req := new(SignUpRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// validate
	if err := c.Validate(req); err != nil {
		return err
	}

	// parse birthday to time.Time from string(YYYY-mm-dd)
	parseBirthDay, err := time.Parse(domain.BirthDayLayout, req.BirthDay)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid date format")
	}

	// sign up usecase
	input := usecase.SignUpUseCaseInput{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		BirthDay: parseBirthDay,
	}
	output, err := uc.uuc.SignUp(input)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// send response
	res := SignUpResponse{
		ID:   output.ID,
		Name: output.Name,
	}

	return c.JSONPretty(http.StatusCreated, res, "  ")
}

func (uc *UserController) GetUser(c echo.Context) error {
	// parse request
	req := new(GetUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// validate
	if err := c.Validate(req); err != nil {
		return err
	}

	// Todo: get user usecase

	// send response
	res := GetUserResponse{
		ID:       req.ID,
		Name:     "example01",
		Email:    "example01",
		BirthDay: "2001-01-01",
	}
	return c.JSONPretty(http.StatusOK, res, "  ")
}

func (ur *UserController) GetUsers(c echo.Context) error {
	// Todo: get users usecase

	// send response
	res := []GetUserResponse{
		{
			ID:       1,
			Name:     "example01",
			Email:    "example01",
			BirthDay: "2001-01-01",
		},
		{
			ID:       2,
			Name:     "example02",
			Email:    "example02",
			BirthDay: "2002-01-01",
		},
	}
	return c.JSONPretty(http.StatusOK, res, "  ")
}
