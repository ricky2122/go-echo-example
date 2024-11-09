package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ricky2122/go-echo-example/controller"
	"github.com/ricky2122/go-echo-example/infrastructure/repository"
	"github.com/ricky2122/go-echo-example/usecase"
	"github.com/uptrace/bun"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewRouter(db *bun.DB) *echo.Echo {
	e := echo.New()

	// set validator
	e.Validator = &CustomValidator{validator: validator.New()}

	ur := repository.NewUserRepository(db)
	uuc := usecase.NewUserUseCase(ur)
	uc := controller.NewUserController(uuc)
	au := controller.NewAuthController()

	e.POST("/signup", uc.SignUp)
	e.GET("/users/:id", uc.GetUser)
	e.GET("/users", uc.GetUsers)

	e.POST("/login", au.Login)
	e.DELETE("/logout", au.Logout)

	return e
}
