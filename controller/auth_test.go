package controller_test

import (
	"fmt"
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

var testSessionUserID = "test01"

type TestStubSessionStore struct {
	sessionsStore map[string]*sessions.Session
	options       *sessions.Options
}

var StubDefaultOpts *sessions.Options = &sessions.Options{
	Path:     "/",
	MaxAge:   86400 * 7,
	HttpOnly: true,
}

func (s *TestStubSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *TestStubSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = s.options
	session.IsNew = true

	c, err := r.Cookie(name)
	if err != nil {
		return session, nil
	}
	session.ID = c.Value

	if _, ok := s.sessionsStore[name]; ok {
		session.IsNew = false
	} else {
		s.sessionsStore[c.Value] = session
	}

	return session, nil
}

func (s *TestStubSessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	cookie := &http.Cookie{
		Name:     session.Name(),
		Value:    "",
		Path:     session.Options.Path,
		MaxAge:   session.Options.MaxAge,
		HttpOnly: session.Options.HttpOnly,
		Secure:   session.Options.Secure,
	}

	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		delete(s.sessionsStore, session.ID)
		http.SetCookie(w, cookie)
		return nil
	}

	if session.ID == "" {
		session.ID = session.Name()
	}

	// Check if session_id is actually present
	value, ok := session.Values[controller.SessionKey]
	if !ok {
		return fmt.Errorf("session_id not found in session value")
	}

	// Ensure the type assertion will not panic
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("session_id value cannot be asserted as string")
	}

	s.sessionsStore[session.ID] = session

	cookie.Value = strValue

	http.SetCookie(w, cookie)

	return nil
}

func TestLogin(t *testing.T) {
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
			targetCookie := fmt.Sprintf("%s=%s", controller.SessionKey, testSessionUserID)
			assert.Contains(t, cookie, targetCookie)
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

func TestLogout(t *testing.T) {
	t.Run("StatusNoContent", func(t *testing.T) {
		e := echo.New()
		// Initialize the session store
		store := &TestStubSessionStore{
			sessionsStore: map[string]*sessions.Session{},
		}

		// Use the correct session middleware
		e.Use(session.Middleware(store))

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		cookie := fmt.Sprintf("%s=%s", controller.SessionKey, testSessionUserID)
		req.Header.Set(echo.HeaderCookie, cookie)

		rec := httptest.NewRecorder()

		sess, _ := store.New(req, controller.SessionKey)
		sess.Values[controller.SessionKey] = testSessionUserID
		sess.Options = StubDefaultOpts
		_ = store.Save(req, rec, sess)

		c := e.NewContext(req, rec)
		c.Set("_session_store", store)

		ac := controller.NewAuthController(&TestStubAuthUseCase{})

		// Assertions
		assert.NoError(t, ac.Logout(c))
		assert.Equal(t, map[string]*sessions.Session{}, store.sessionsStore)
	})
}
