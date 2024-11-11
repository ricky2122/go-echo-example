package controller_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

type TestStubSessionStore struct {
	sessionsStore map[string]*sessions.Session
}

func (s *TestStubSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	sess, ok := s.sessionsStore[name]
	if !ok {
		newSess := sessions.NewSession(s, name)
		s.sessionsStore[name] = newSess
		return newSess, nil
	}
	return sess, nil
}

func (s *TestStubSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	newSess := sessions.NewSession(s, name)
	s.sessionsStore[name] = newSess
	return newSess, nil
}

func (s *TestStubSessionStore) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	// Check if session_id is actually present
	value, ok := sess.Values["session_id"]
	if !ok {
		return fmt.Errorf("session_id not found in session value")
	}

	// Ensure the type assertion will not panic
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("session_id value cannot be asserted as string")
	}

	cookie := &http.Cookie{
		Name:     sess.Name(),
		Value:    strValue,
		Path:     sess.Options.Path,
		MaxAge:   sess.Options.MaxAge,
		HttpOnly: sess.Options.HttpOnly,
		Secure:   sess.Options.Secure,
	}
	http.SetCookie(w, cookie)

	// Log success for debugging
	log.Printf("Session saved: %v", sess.Values)
	return nil
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

		// Initialize the session store
		store := &TestStubSessionStore{
			sessionsStore: map[string]*sessions.Session{},
		}

		// Use the correct session middleware
		e.Use(session.Middleware(store))

		// Set the validator
		e.Validator = api.NewCustomValidator()

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginReq))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.Set("_session_store", store)

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
			// Check if the session cookie is correctly set
			cookie := rec.Header().Get("Set-Cookie")
			// Verify the session ID value within the cookie
			assert.Contains(t, cookie, "session_id=test_session_id")
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
