package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	BirthDay string `json:"birth_day"`
}

type SignUpResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetUserRequest struct {
	ID int `param:"id"`
}

type GetUserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	BirthDay string `json:"birth_day"`
}

func NewRouter() *echo.Echo {
	e := echo.New()

	e.POST("/signup", func(c echo.Context) error {
		// parse request
		req := new(SignUpRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}

		// Todo: validate

		// Todo: sign up usecase

		// send response
		res := SignUpResponse{ID: 1, Name: req.Name}
		return c.JSONPretty(http.StatusCreated, res, "  ")
	})

	e.GET("/users/:id", func(c echo.Context) error {
		// parse request
		req := new(GetUserRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}

		// Todo: validate

		// Todo: get user usecase

		// send response
		res := GetUserResponse{
			ID:       req.ID,
			Name:     "example01",
			Email:    "example01",
			BirthDay: "2001-01-01",
		}
		return c.JSONPretty(http.StatusOK, res, "  ")
	})

	e.GET("/users", func(c echo.Context) error {
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
	})

	e.POST("/login", func(c echo.Context) error {
		// parse request
		req := new(LoginRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}

		// Todo: validate

		// Todo: login usecase

		// Todo: set session_id to cookie

		// send response
		return c.NoContent(http.StatusOK)
	})

	e.DELETE("/logout", func(c echo.Context) error {
		// Todo: get session_id from cookie

		// Todo: logout usecase

		// send response
		return c.NoContent(http.StatusNoContent)
	})

	return e
}

func main() {
	router := NewRouter()

	router.Logger.Fatal(router.Start(":1323"))
}
