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

type TestStubUserUseCase struct {
	signUpOutputStore   map[string]*usecase.SignUpUseCaseOutput
	getUserOutputStore  map[int]*usecase.GetUserUseCaseOutput
	getUsersOutputStore usecase.GetUsersUseCaseOutput
}

func (s *TestStubUserUseCase) SignUp(input usecase.SignUpUseCaseInput) (*usecase.SignUpUseCaseOutput, error) {
	_, ok := s.signUpOutputStore[input.Name]
	if ok {
		return nil, usecase.ErrUserAlreadyExists
	}
	output := &usecase.SignUpUseCaseOutput{ID: len(s.signUpOutputStore) + 1, Name: input.Name}
	s.signUpOutputStore[input.Name] = output

	return output, nil
}

func (s *TestStubUserUseCase) GetUser(input usecase.GetUserUseCaseInput) (*usecase.GetUserUseCaseOutput, error) {
	output, ok := s.getUserOutputStore[input.ID]
	if !ok {
		return nil, usecase.ErrUserNotFound
	}
	return output, nil
}

func (s *TestStubUserUseCase) GetUsers() (*usecase.GetUsersUseCaseOutput, error) {
	return &s.getUsersOutputStore, nil
}

func TestSignUp(t *testing.T) {
	signUpReq01 := `{
		"name":"test01",
		"password":"test01",
		"email":"test01@test.com",
		"birth_day":"2001-01-01"
	  }
	  `

	signUpRes01 := `{
		"id": 1,
		"name": "test01"
	  }
	  `

	signUpReq02 := `{
		"name":"test02",
		"password":"test02",
		"email":"test02@test.com",
		"birth_day":"2002-01-01"
	  }
	  `

	signUpRes02 := `{
		"id": 2,
		"name": "test02"
	  }
	  `
	// Setup
	e := echo.New()
	e.Validator = api.NewCustomValidator()

	t.Run("StatusCreated", func(t *testing.T) {
		store := map[string]*usecase.SignUpUseCaseOutput{}
		uc := controller.NewUserController(&TestStubUserUseCase{signUpOutputStore: store})

		cases := []struct {
			name    string
			reqJSON string
			want    string
		}{
			{
				name:    "id:1",
				reqJSON: signUpReq01,
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

	t.Run("StatusBadRequest", func(t *testing.T) {
		userAlreadyExistsMsg := "user already exists"

		store := map[string]*usecase.SignUpUseCaseOutput{}
		store["test01"] = &usecase.SignUpUseCaseOutput{ID: 1, Name: "test01"}
		uc := controller.NewUserController(&TestStubUserUseCase{signUpOutputStore: store})

		t.Run("user already exists", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(signUpReq01))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			// Assertions
			err := uc.SignUp(c)
			if assert.NotNil(t, err) {
				err, res := err.(*echo.HTTPError)
				if res {
					assert.Equal(t, http.StatusBadRequest, err.Code)
					assert.Equal(t, userAlreadyExistsMsg, err.Message)
				}
			}
		})
	})
}

func TestGetUser(t *testing.T) {
	getUserRes01 := `{
		"id": 1,
		"name": "test01",
		"email": "test01@test.com",
		"birth_day": "2001-01-01"
	  }
	  `

	getUserRes02 := `{
		"id": 2,
		"name": "test02",
		"email": "test02@test.com",
		"birth_day": "2002-01-01"
	  }
	  `

	userNotFoundMsg := "user not found"
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
		// Set up
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

func TestGetUsers(t *testing.T) {
	// Set up
	e := echo.New()
	e.Validator = api.NewCustomValidator()

	getUsersEmptyRes := `
	{
  	  "users": []
    }`
	getUsersTwoUsersRes := `
	{
	  "users": [
		{
			"id": 1,
			"name": "test01",
			"email": "test01@test.com",
			"birth_day": "2001-01-01"
		},
		{
			"id": 2,
			"name": "test02",
			"email": "test02@test.com",
			"birth_day": "2002-01-01"
		}
	  ]
	}
	`
	t.Run("StatusOK", func(t *testing.T) {
		cases := []struct {
			name string
			want string
		}{
			{
				name: "empty",
				want: getUsersEmptyRes,
			},
			{
				name: "two users",
				want: getUsersTwoUsersRes,
			},
		}

		for _, tt := range cases {
			// Create a test user controller with a stub user use case
			var store usecase.GetUsersUseCaseOutput
			switch tt.name {
			case "empty":
				store = usecase.GetUsersUseCaseOutput{
					Users: []usecase.GetUserUseCaseOutput{},
				}
			case "two users":
				store = usecase.GetUsersUseCaseOutput{
					Users: []usecase.GetUserUseCaseOutput{
						{
							ID:       1,
							Name:     "test01",
							Email:    "test01@test.com",
							BirthDay: "2001-01-01",
						},
						{
							ID:       2,
							Name:     "test02",
							Email:    "test02@test.com",
							BirthDay: "2002-01-01",
						},
					},
				}
			}
			uc := controller.NewUserController(&TestStubUserUseCase{
				getUsersOutputStore: store,
			})
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/users", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				c := e.NewContext(req, rec)

				// Assertions
				if assert.NoError(t, uc.GetUsers(c)) {
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.JSONEq(t, tt.want, rec.Body.String())
				}
			})
		}
	})
}
