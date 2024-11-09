package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ricky2122/go-echo-example/controller"
	"github.com/ricky2122/go-echo-example/infrastructure/api"
	"github.com/ricky2122/go-echo-example/usecase"
	"github.com/stretchr/testify/assert"
)

var singUpReq01 = `{
  "name":"test01",
  "password":"test01",
  "email":"test01@test.com",
  "birth_day":"2001-01-01"
}
`

var signUpRes01 = `{
  "id": 1,
  "name": "test01"
}
`

var signUpReq02 = `{
  "name":"test02",
  "password":"test02",
  "email":"test02@test.com",
  "birth_day":"2002-01-01"
}
`

var signUpRes02 = `{
  "id": 2,
  "name": "test02"
}
`

var getUserRes01 = `{
  "id": 1,
  "name": "test01",
  "email": "test01@test.com",
  "birth_day": "2001-01-01"
}
`

var getUserRes02 = `{
  "id": 2,
  "name": "test02",
  "email": "test02@test.com",
  "birth_day": "2002-01-01"
}
`

var userNotFoundRes = `{
  "message": "user not found"
}
`

var userNotFoundMsg = "user not found"

type TestStubUserUseCase struct {
	counter            int
	getUserOutputStore map[int]*usecase.GetUserUseCaseOutput
}

func (s *TestStubUserUseCase) SignUp(input usecase.SignUpUseCaseInput) (*usecase.SignUpUseCaseOutput, error) {
	s.counter++
	output := &usecase.SignUpUseCaseOutput{ID: s.counter, Name: input.Name}
	return output, nil
}

func (s *TestStubUserUseCase) GetUser(input usecase.GetUserUseCaseInput) (*usecase.GetUserUseCaseOutput, error) {
	output, ok := s.getUserOutputStore[input.ID]
	if !ok {
		return nil, usecase.ErrUserNotFound
	}
	return output, nil
}

func TestSignUp(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = api.NewCustomValidator()

	uc := controller.NewUserController(&TestStubUserUseCase{})

	t.Run("StatusCreated", func(t *testing.T) {
		cases := []struct {
			name    string
			reqJSON string
			want    string
		}{
			{
				name:    "id:1",
				reqJSON: singUpReq01,
				want:    signUpRes01,
			},
			{
				name:    "id:2",
				reqJSON: signUpReq02,
				want:    signUpRes02,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(tt.reqJSON))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				c := e.NewContext(req, rec)

				// Assertions
				if assert.NoError(t, uc.SignUp(c)) {
					assert.Equal(t, http.StatusCreated, rec.Code)
					assert.JSONEq(t, tt.want, rec.Body.String())
				}
			})
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("StatusOK", func(t *testing.T) {
		// Setup
		e := echo.New()
		e.Validator = api.NewCustomValidator()

		store := map[int]*usecase.GetUserUseCaseOutput{}
		store[1] = &usecase.GetUserUseCaseOutput{
			ID:       1,
			Name:     "test01",
			Email:    "test01@test.com",
			BirthDay: "2001-01-01",
		}
		store[2] = &usecase.GetUserUseCaseOutput{
			ID:       2,
			Name:     "test02",
			Email:    "test02@test.com",
			BirthDay: "2002-01-01",
		}
		uc := controller.NewUserController(&TestStubUserUseCase{
			getUserOutputStore: store,
		})

		cases := []struct {
			name string
			id   int
			want string
		}{
			{
				name: "id:1",
				id:   1,
				want: getUserRes01,
			},
			{
				name: "id:2",
				id:   2,
				want: getUserRes02,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				c := e.NewContext(req, rec)
				c.SetPath("/users/:id")
				c.SetParamNames("id")
				c.SetParamValues(strconv.Itoa(tt.id))

				// Assertions
				if assert.NoError(t, uc.GetUser(c)) {
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.JSONEq(t, tt.want, rec.Body.String())
				}
			})
		}
	})

	t.Run("StatusNotFound", func(t *testing.T) {
		// Set up Echo with a custom error handler
		e := echo.New()
		e.Validator = api.NewCustomValidator()

		// Create a test user controller with a stub user use case
		store := map[int]*usecase.GetUserUseCaseOutput{}
		uc := controller.NewUserController(&TestStubUserUseCase{
			getUserOutputStore: store,
		})

		// Create a request and recorder
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Create a context
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Assertions
		err := uc.GetUser(c)
		if assert.NotNil(t, err) {
			err, res := err.(*echo.HTTPError)
			if res {
				assert.Equal(t, http.StatusNotFound, err.Code)
				assert.Equal(t, userNotFoundMsg, err.Message)
			}
		}
	})
}
