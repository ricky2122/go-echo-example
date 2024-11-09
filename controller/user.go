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

type GetUsersResponse struct {
	Users []GetUserResponse `json:"users"`
}

type IUserUseCase interface {
	SignUp(usecase.SignUpUseCaseInput) (*usecase.SignUpUseCaseOutput, error)
	GetUser(usecase.GetUserUseCaseInput) (*usecase.GetUserUseCaseOutput, error)
	GetUsers() (*usecase.GetUsersUseCaseOutput, error)
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

	// get user usecase
	input := usecase.GetUserUseCaseInput{ID: req.ID}
	output, err := uc.uuc.GetUser(input)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// send response
	res := GetUserResponse{
		ID:       output.ID,
		Name:     output.Name,
		Email:    output.Email,
		BirthDay: output.BirthDay,
	}
	return c.JSONPretty(http.StatusOK, res, "  ")
}

func (ur *UserController) GetUsers(c echo.Context) error {
	// get users usecase
	output, err := ur.uuc.GetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// user is empty
	if output == nil {
		return c.JSONPretty(http.StatusOK, GetUsersResponse{Users: []GetUserResponse{}}, "  ")
	}

	// send response
	users := make([]GetUserResponse, 0, len(output.Users))
	for _, outputUser := range output.Users {
		userResp := GetUserResponse{
			ID:       outputUser.ID,
			Name:     outputUser.Name,
			Email:    outputUser.Email,
			BirthDay: outputUser.BirthDay,
		}
		users = append(users, userResp)
	}
	res := GetUsersResponse{Users: users}

	return c.JSONPretty(http.StatusOK, res, "  ")
}
