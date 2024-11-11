package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ricky2122/go-echo-example/controller"
	"github.com/ricky2122/go-echo-example/domain"
	"github.com/ricky2122/go-echo-example/infrastructure/api"
	"github.com/ricky2122/go-echo-example/usecase"
	"github.com/stretchr/testify/assert"
)

type TestStubAuthUseCase struct {
	userStore []domain.User
}

func (s *TestStubAuthUseCase) Login(input usecase.LoginUseCaseInput) error {
	for _, user := range s.userStore {
		if input.Name == user.GetName() && input.Password == user.GetPassword() {
			return nil
		}
	}
	return usecase.ErrLoginFailed
}

func TestLoginTest(t *testing.T) {
	loginReq := `{
	  "name": "test01",
	  "password": "test01"
	}
	`
	t.Run("StatusOK", func(t *testing.T) {
		// Setup
		e := echo.New()
		e.Validator = api.NewCustomValidator()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginReq))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		user := domain.NewUser(
			"test01",
			"test01",
			"test01@test.com",
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		user.SetID(1)
		userStore := []domain.User{user}
		ac := controller.NewAuthController(&TestStubAuthUseCase{userStore: userStore})

		// Assertions
		if assert.NoError(t, ac.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("StatusUnAuthorized", func(t *testing.T) {
		loginFailedMsg := "failed login"
		// Setup
		e := echo.New()
		e.Validator = api.NewCustomValidator()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginReq))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ac := controller.NewAuthController(&TestStubAuthUseCase{userStore: []domain.User{}})

		// Assertions
		err := ac.Login(c)
		if assert.NotNil(t, err) {
			err, res := err.(*echo.HTTPError)
			if res {
				assert.Equal(t, http.StatusUnauthorized, err.Code)
				assert.Equal(t, loginFailedMsg, err.Message)
			}
		}
	})
}
