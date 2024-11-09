package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ricky2122/go-echo-example/controller"
	"github.com/ricky2122/go-echo-example/infrastructure/api"
	"github.com/ricky2122/go-echo-example/usecase"
	"github.com/stretchr/testify/assert"
)

var signUpRes01 = `{
  "id": 1,
  "name": "user01"
}
`

var signUpRes02 = `{
  "id": 2,
  "name": "user02"
}
`

type TestStubUserUseCase struct {
	counter int
}

func (s *TestStubUserUseCase) SignUp(input usecase.SignUpUseCaseInput) (*usecase.SignUpUseCaseOutput, error) {
	s.counter++
	output := &usecase.SignUpUseCaseOutput{ID: s.counter, Name: input.Name}
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
				reqJSON: `{"name":"user01","password":"example01","email":"example01@example.com","birth_day":"2001-01-01"}`,
				want:    signUpRes01,
			},
			{
				name:    "id:2",
				reqJSON: `{"name":"user02","password":"example02","email":"example02@example.com","birth_day":"2002-01-01"}`,
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
					assert.Equal(t, tt.want, rec.Body.String())
				}
			})
		}
	})
}
